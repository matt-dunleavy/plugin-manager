package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "time"

    pm "github.com/matt-dunleavy/plugin-manager"
)

func main() {
    // Get the current working directory
    cwd, err := os.Getwd()
    if err != nil {
        log.Fatalf("Failed to get current working directory: %v", err)
    }

    // Set up paths
    configPath := filepath.Join(cwd, "plugins.json")
    pluginDir := filepath.Join(cwd, "plugins")

    // Ensure plugin directory exists
    if err := os.MkdirAll(pluginDir, 0755); err != nil {
        log.Fatalf("Failed to create plugin directory: %v", err)
    }

    // Create a new plugin manager
    manager, err := pm.NewManager(configPath, pluginDir)
    if err != nil {
        log.Fatalf("Failed to create plugin manager: %v", err)
    }

    // Subscribe to plugin events
    manager.SubscribeToEvent("PluginLoaded", func(e pm.Event) {
        fmt.Printf("Event: Plugin loaded - %s\n", e.(pm.PluginLoadedEvent).PluginName)
    })
    manager.SubscribeToEvent("PluginUnloaded", func(e pm.Event) {
        fmt.Printf("Event: Plugin unloaded - %s\n", e.(pm.PluginUnloadedEvent).PluginName)
    })

    // Load plugins
    pluginsToLoad := []string{"hello.so", "math.so"}
    for _, plugin := range pluginsToLoad {
        pluginPath := filepath.Join(pluginDir, plugin)
        err := manager.LoadPlugin(pluginPath)
        if err != nil {
            log.Printf("Failed to load plugin %s: %v", plugin, err)
        }
    }

    // List loaded plugins
    loadedPlugins := manager.ListPlugins()
    fmt.Println("Loaded plugins:", loadedPlugins)

    // Execute plugins
    pluginsToExecute := []string{"HelloPlugin", "MathPlugin"}
    for _, plugin := range pluginsToExecute {
        err := manager.ExecutePlugin(plugin)
        if err != nil {
            log.Printf("Failed to execute %s: %v", plugin, err)
        }
    }

    // Get plugin stats
    for _, plugin := range pluginsToExecute {
        stats, err := manager.GetPluginStats(plugin)
        if err != nil {
            fmt.Printf("%s stats not available: %v\n", plugin, err)
        } else {
            fmt.Printf("%s stats: %+v\n", plugin, stats)
        }
    }

    // Hot-reload HelloPlugin
    time.Sleep(2 * time.Second) // Wait to simulate some time passing
    fmt.Println("\nHot-reloading HelloPlugin...")
    err = manager.HotReload("HelloPlugin", filepath.Join(pluginDir, "hello.so"))
    if err != nil {
        log.Printf("Failed to hot-reload HelloPlugin: %v", err)
    }

    // Execute hot-reloaded plugin
    err = manager.ExecutePlugin("HelloPlugin")
    if err != nil {
        log.Printf("Failed to execute hot-reloaded HelloPlugin: %v", err)
    }

    // Unload plugins
    for _, plugin := range loadedPlugins {
        err := manager.UnloadPlugin(plugin)
        if err != nil {
            log.Printf("Failed to unload %s: %v", plugin, err)
        }
    }

    fmt.Println("\nFinal list of loaded plugins:", manager.ListPlugins())
}