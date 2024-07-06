package main

import (
    "fmt"
    "log"
    "path/filepath"

    pm "github.com/matt-dunleavy/plugin-manager"
)

func main() {
    // Initialize the plugin manager
    manager, err := pm.NewManager("plugins.json", "./plugins", "public_key.pem")
    if err != nil {
        log.Fatalf("Failed to create plugin manager: %v", err)
    }

    // Subscribe to plugin events
    manager.SubscribeToEvent("PluginLoaded", func(e pm.Event) {
        fmt.Printf("Plugin loaded: %s\n", e.(pm.PluginLoadedEvent).PluginName)
    })

    // Load plugins
    plugins := []string{"hello.so", "math.so"}
    for _, plugin := range plugins {
        err := manager.LoadPlugin(filepath.Join("./plugins", plugin))
        if err != nil {
            log.Printf("Failed to load plugin %s: %v", plugin, err)
        }
    }

    // List loaded plugins
    loadedPlugins := manager.ListPlugins()
    fmt.Println("Loaded plugins:", loadedPlugins)

    // Execute plugins
    for _, name := range loadedPlugins {
        err := manager.ExecutePlugin(name)
        if err != nil {
            log.Printf("Failed to execute plugin %s: %v", name, err)
        }
    }

    // Get and print plugin stats
    for _, name := range loadedPlugins {
        stats, err := manager.GetPluginStats(name)
        if err != nil {
            log.Printf("Failed to get stats for plugin %s: %v", name, err)
        } else {
            fmt.Printf("Stats for %s: Executions: %d, Last execution time: %v\n",
                name, stats.ExecutionCount, stats.LastExecutionTime)
        }
    }

    // Unload a plugin
    err = manager.UnloadPlugin("hello.so")
    if err != nil {
        log.Printf("Failed to unload plugin: %v", err)
    }

    fmt.Println("Remaining plugins:", manager.ListPlugins())
}