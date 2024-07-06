package main

import (
    "fmt"
    "log"
    "time"

    pm "github.com/matt-dunleavy/plugin-manager"
)

func main() {
    // Create a new plugin manager
    manager, err := pm.NewManager("../plugins.json", "../plugins")
    if err != nil {
        log.Fatalf("Failed to create plugin manager: %v", err)
    }

    // Subscribe to plugin events
    manager.GetEventBus().Subscribe("PluginLoaded", func(e pm.Event) {
        fmt.Printf("Event: Plugin loaded - %s\n", e.(pm.PluginLoadedEvent).PluginName)
    })
    manager.GetEventBus().Subscribe("PluginUnloaded", func(e pm.Event) {
        fmt.Printf("Event: Plugin unloaded - %s\n", e.(pm.PluginUnloadedEvent).PluginName)
    })

    // Load plugins
    err = manager.LoadPlugin("../plugins/hello.so")
    if err != nil {
        log.Printf("Failed to load hello plugin: %v", err)
    }

    err = manager.LoadPlugin("../plugins/math.so")
    if err != nil {
        log.Printf("Failed to load math plugin: %v", err)
    }

    // List loaded plugins
    plugins := manager.ListPlugins()
    fmt.Println("Loaded plugins:", plugins)

    // Execute plugins
    err = manager.ExecutePlugin("HelloPlugin")
    if err != nil {
        log.Printf("Failed to execute HelloPlugin: %v", err)
    }

    err = manager.ExecutePlugin("MathPlugin")
    if err != nil {
        log.Printf("Failed to execute MathPlugin: %v", err)
    }

    // Get plugin stats
    helloStats := manager.GetPluginStats("HelloPlugin")
    mathStats := manager.GetPluginStats("MathPlugin")

    fmt.Printf("HelloPlugin stats: %+v\n", helloStats)
    fmt.Printf("MathPlugin stats: %+v\n", mathStats)

    // Hot-reload HelloPlugin
    time.Sleep(2 * time.Second) // Wait to simulate some time passing
    fmt.Println("\nHot-reloading HelloPlugin...")
    err = manager.HotReload("HelloPlugin", "../plugins/hello.so")
    if err != nil {
        log.Printf("Failed to hot-reload HelloPlugin: %v", err)
    }

    // Execute hot-reloaded plugin
    err = manager.ExecutePlugin("HelloPlugin")
    if err != nil {
        log.Printf("Failed to execute hot-reloaded HelloPlugin: %v", err)
    }

    // Unload plugins
    err = manager.UnloadPlugin("HelloPlugin")
    if err != nil {
        log.Printf("Failed to unload HelloPlugin: %v", err)
    }

    err = manager.UnloadPlugin("MathPlugin")
    if err != nil {
        log.Printf("Failed to unload MathPlugin: %v", err)
    }
}