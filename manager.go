// File: internal/plugin/manager.go

package pluginmanager

import (
    "fmt"
    "sync"
    "time"
    "path/filepath"
)

type Manager struct {
    plugins      map[string]Plugin
    config       *Config
    dependencies map[string][]string
    stats        map[string]*PluginStats
    eventBus     *EventBus
    sandbox      Sandbox
    mu           sync.RWMutex
}

func NewManager(configPath string, pluginDir string) (*Manager, error) {
    config, err := LoadConfig(configPath)
    if err != nil {
        return nil, err
    }

    return &Manager{
        plugins:      make(map[string]Plugin),
        config:       config,
        dependencies: make(map[string][]string),
        stats:        make(map[string]*PluginStats),
        eventBus:     NewEventBus(),
        sandbox:      NewDefaultSandbox(pluginDir),
    }, nil
}

func (m *Manager) LoadPlugin(path string) error {
    if err := m.sandbox.VerifyPluginPath(path); err != nil {
        return err
    }

    if err := m.sandbox.Enable(); err != nil {
        return err
    }
    defer m.sandbox.Disable()

    plugin, err := LoadPlugin(path)
    if err != nil {
        return err
    }

    m.mu.Lock()
    defer m.mu.Unlock()

    metadata := plugin.Metadata()
    if _, exists := m.plugins[metadata.Name]; exists {
        return ErrPluginAlreadyLoaded
    }

    if err := m.checkDependencies(metadata); err != nil {
        return err
    }

    m.plugins[metadata.Name] = plugin
    m.dependencies[metadata.Name] = metadata.Dependencies
    m.stats[metadata.Name] = &PluginStats{}

    m.eventBus.Publish(PluginLoadedEvent{PluginName: metadata.Name})

    return plugin.Init()
}

func (m *Manager) UnloadPlugin(name string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    plugin, exists := m.plugins[name]
    if !exists {
        return ErrPluginNotFound
    }

    if err := plugin.Shutdown(); err != nil {
        return err
    }

    delete(m.plugins, name)
    delete(m.dependencies, name)
    delete(m.stats, name)

    m.eventBus.Publish(PluginUnloadedEvent{PluginName: name})

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
        return err
    }
    defer m.sandbox.Disable()

    start := time.Now()
    err := plugin.Execute()
    executionTime := time.Since(start)

    m.mu.Lock()
    stats.ExecutionCount++
    stats.LastExecutionTime = executionTime
    stats.TotalExecutionTime += executionTime
    m.mu.Unlock()

    return err
}

func (m *Manager) HotReload(name string, path string) error {
    if err := m.sandbox.VerifyPluginPath(path); err != nil {
        return err
    }

    if err := m.sandbox.Enable(); err != nil {
        return err
    }
    defer m.sandbox.Disable()

    newPlugin, err := LoadPlugin(path)
    if err != nil {
        return err
    }

    m.mu.Lock()
    defer m.mu.Unlock()

    oldPlugin, exists := m.plugins[name]
    if !exists {
        return ErrPluginNotFound
    }

    if err := oldPlugin.Shutdown(); err != nil {
        return err
    }

    m.plugins[name] = newPlugin
    m.eventBus.Publish(PluginHotReloadedEvent{PluginName: name})

    return newPlugin.Init()
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

func (m *Manager) checkDependencies(metadata PluginMetadata) error {
    visited := make(map[string]bool)
    return m.dfs(metadata.Name, visited)
}

func (m *Manager) dfs(name string, visited map[string]bool) error {
    if visited[name] {
        return ErrCircularDependency
    }

    visited[name] = true
    for _, dep := range m.dependencies[name] {
        if _, exists := m.plugins[dep]; !exists {
            return fmt.Errorf("%w: %s", ErrMissingDependency, dep)
        }
        if err := m.dfs(dep, visited); err != nil {
            return err
        }
    }
    visited[name] = false
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