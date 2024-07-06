// FILE: manager.go

package pluginmanager

import (
    "fmt"
    "path/filepath"
    "plugin"
    "strconv"
    "strings"
    "sync"
    "time"

    "go.uber.org/zap"
)

type Manager struct {
    plugins       map[string]*lazyPlugin
    config        *Config
    dependencies  map[string][]string
    stats         map[string]*PluginStats
    eventBus      *EventBus
    sandbox       Sandbox
    logger        *zap.Logger
    publicKeyPath string
    mu            sync.RWMutex
}

type lazyPlugin struct {
    path   string
    loaded Plugin
}

func (lp *lazyPlugin) load() error {
    if lp.loaded == nil {
        p, err := plugin.Open(lp.path)
        if err != nil {
            return fmt.Errorf("failed to open plugin: %w", err)
        }

        symPlugin, err := p.Lookup(PluginSymbol)
        if err != nil {
            return fmt.Errorf("failed to lookup plugin symbol: %w", err)
        }

        plugin, ok := symPlugin.(Plugin)
        if !ok {
            return fmt.Errorf("invalid plugin interface")
        }

        lp.loaded = plugin
    }
    return nil
}

func NewManager(configPath, pluginDir, publicKeyPath string) (*Manager, error) {
    config, err := LoadConfig(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load config: %w", err)
    }

    logger, _ := zap.NewProduction()

    sandboxDir := filepath.Join(pluginDir, "sandbox")

    return &Manager{
        plugins:       make(map[string]*lazyPlugin),
        config:        config,
        dependencies:  make(map[string][]string),
        stats:         make(map[string]*PluginStats),
        eventBus:      NewEventBus(),
        sandbox:       NewLinuxSandbox(sandboxDir),
        logger:        logger,
        publicKeyPath: publicKeyPath,
    }, nil
}

func (m *Manager) LoadPlugin(path string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    pluginName := filepath.Base(path)
    if _, exists := m.plugins[pluginName]; exists {
        return fmt.Errorf("plugin %s already loaded", pluginName)
    }

    if err := m.VerifyPluginSignature(path, m.publicKeyPath); err != nil {
        return fmt.Errorf("failed to verify plugin signature: %w", err)
    }

    lazyPlug := &lazyPlugin{path: path}
    if err := lazyPlug.load(); err != nil {
        return fmt.Errorf("failed to load plugin %s: %w", pluginName, err)
    }

    plugin := lazyPlug.loaded

    if err := plugin.PreLoad(); err != nil {
        return fmt.Errorf("pre-load hook failed for %s: %w", pluginName, err)
    }

    if err := plugin.Init(); err != nil {
        return fmt.Errorf("initialization failed for %s: %w", pluginName, err)
    }

    if err := plugin.PostLoad(); err != nil {
        return fmt.Errorf("post-load hook failed for %s: %w", pluginName, err)
    }

    m.plugins[pluginName] = lazyPlug
    m.stats[pluginName] = &PluginStats{}

    metadata := plugin.Metadata()
    m.dependencies[pluginName] = make([]string, 0, len(metadata.Dependencies))
    for dep, constraint := range metadata.Dependencies {
        m.dependencies[pluginName] = append(m.dependencies[pluginName], dep)
        if err := m.checkDependency(dep, constraint); err != nil {
            delete(m.plugins, pluginName)
            delete(m.stats, pluginName)
            delete(m.dependencies, pluginName)
            return fmt.Errorf("dependency check failed for %s: %w", pluginName, err)
        }
    }

    m.eventBus.Publish(PluginLoadedEvent{PluginName: pluginName})
    m.logger.Info("Plugin loaded", zap.String("plugin", pluginName))

    return nil
}

func (m *Manager) UnloadPlugin(name string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    plugin, exists := m.plugins[name]
    if !exists {
        return ErrPluginNotFound
    }

    if err := plugin.loaded.PreUnload(); err != nil {
        return fmt.Errorf("pre-unload hook failed for %s: %w", name, err)
    }

    if err := plugin.loaded.Shutdown(); err != nil {
        return fmt.Errorf("shutdown failed for %s: %w", name, err)
    }

    delete(m.plugins, name)
    delete(m.dependencies, name)
    delete(m.stats, name)

    m.eventBus.Publish(PluginUnloadedEvent{PluginName: name})
    m.logger.Info("Plugin unloaded", zap.String("plugin", name))

    return nil
}

func (m *Manager) ExecutePlugin(name string) error {
    m.mu.RLock()
    plugin, exists := m.plugins[name]
    stats := m.stats[name]
    m.mu.RUnlock()

    if !exists {
        return ErrPluginNotFound
    }

    if err := m.sandbox.Enable(); err != nil {
        return fmt.Errorf("failed to enable sandbox for %s: %w", name, err)
    }
    defer m.sandbox.Disable()

    if err := plugin.load(); err != nil {
        return fmt.Errorf("failed to load plugin %s: %w", name, err)
    }

    start := time.Now()
    err := plugin.loaded.Execute()
    executionTime := time.Since(start)

    m.mu.Lock()
    stats.ExecutionCount++
    stats.LastExecutionTime = executionTime
    stats.TotalExecutionTime += executionTime
    m.mu.Unlock()

    if err != nil {
        return fmt.Errorf("execution failed for %s: %w", name, err)
    }

    m.logger.Info("Plugin executed", zap.String("plugin", name), zap.Duration("duration", executionTime))
    return nil
}

func (m *Manager) HotReload(name string, path string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    oldPlugin, ok := m.plugins[name]
    if !ok {
        return ErrPluginNotFound
    }

    if err := m.VerifyPluginSignature(path, m.publicKeyPath); err != nil {
        return fmt.Errorf("failed to verify new plugin signature: %w", err)
    }

    newLazyPlugin := &lazyPlugin{path: path}
    if err := newLazyPlugin.load(); err != nil {
        return fmt.Errorf("failed to load new version of %s: %w", name, err)
    }

    newPlugin := newLazyPlugin.loaded

    metadata := newPlugin.Metadata()
    for dep, constraint := range metadata.Dependencies {
        if err := m.checkDependency(dep, constraint); err != nil {
            return fmt.Errorf("dependency check failed for new version of %s: %w", name, err)
        }
    }

    if err := newPlugin.Init(); err != nil {
        return fmt.Errorf("initialization failed for new version of %s: %w", name, err)
    }

    if err := oldPlugin.loaded.PreUnload(); err != nil {
        m.logger.Warn("Pre-unload hook failed for old version", zap.String("plugin", name), zap.Error(err))
    }
    if err := oldPlugin.loaded.Shutdown(); err != nil {
        m.logger.Warn("Shutdown failed for old version", zap.String("plugin", name), zap.Error(err))
    }

    m.plugins[name] = newLazyPlugin
    m.dependencies[name] = make([]string, 0, len(metadata.Dependencies))
    for dep := range metadata.Dependencies {
        m.dependencies[name] = append(m.dependencies[name], dep)
    }

    m.eventBus.Publish(PluginHotReloadedEvent{PluginName: name})
    m.logger.Info("Plugin hot-reloaded", zap.String("plugin", name))

    return nil
}

func (m *Manager) checkDependency(depName, constraint string) error {
    depPlugin, exists := m.plugins[depName]
    if !exists {
        return fmt.Errorf("missing dependency: %s", depName)
    }

    if err := depPlugin.load(); err != nil {
        return fmt.Errorf("failed to load dependency %s: %w", depName, err)
    }

    depVersion := depPlugin.loaded.Metadata().Version
    if !isVersionCompatible(depVersion, constraint) {
        return fmt.Errorf("incompatible version for dependency %s: required %s, got %s", depName, constraint, depVersion)
    }

    return nil
}

func isVersionCompatible(currentVersion, constraint string) bool {
    parts := strings.Split(constraint, " ")
    if len(parts) != 2 {
        return false
    }

    operator := parts[0]
    requiredVersion := parts[1]

    switch operator {
    case ">=":
        return compareVersions(currentVersion, requiredVersion) >= 0
    case ">":
        return compareVersions(currentVersion, requiredVersion) > 0
    case "<=":
        return compareVersions(currentVersion, requiredVersion) <= 0
    case "<":
        return compareVersions(currentVersion, requiredVersion) < 0
    case "==":
        return compareVersions(currentVersion, requiredVersion) == 0
    default:
        return false
    }
}

func compareVersions(v1, v2 string) int {
    parts1 := strings.Split(v1, ".")
    parts2 := strings.Split(v2, ".")

    for i := 0; i < len(parts1) && i < len(parts2); i++ {
        n1, _ := strconv.Atoi(parts1[i])
        n2, _ := strconv.Atoi(parts2[i])

        if n1 < n2 {
            return -1
        } else if n1 > n2 {
            return 1
        }
    }

    if len(parts1) < len(parts2) {
        return -1
    } else if len(parts1) > len(parts2) {
        return 1
    }

    return 0
}

func (m *Manager) EnablePlugin(name string) error {
    if err := m.config.EnablePlugin(name); err != nil {
        return err
    }
    return m.config.Save()
}

func (m *Manager) DisablePlugin(name string) error {
    if err := m.config.DisablePlugin(name); err != nil {
        return err
    }
    return m.config.Save()
}

func (m *Manager) LoadEnabledPlugins(pluginDir string) error {
    enabled := m.config.EnabledPlugins()
    for _, name := range enabled {
        path := filepath.Join(pluginDir, name+".so")
        if err := m.LoadPlugin(path); err != nil {
            return err
        }
    }
    return nil
}

func (m *Manager) ListPlugins() []string {
    m.mu.RLock()
    defer m.mu.RUnlock()

    plugins := make([]string, 0, len(m.plugins))
    for name := range m.plugins {
        plugins = append(plugins, name)
    }
    return plugins
}

func (m *Manager) GetPluginStats(name string) (*PluginStats, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    stats, ok := m.stats[name]
    if !ok {
        return nil, ErrPluginNotFound
    }
    return stats, nil
}

func (m *Manager) SubscribeToEvent(eventName string, handler EventHandler) {
    m.eventBus.Subscribe(eventName, handler)
}
