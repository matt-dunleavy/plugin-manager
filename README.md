# Plugin Manager for Go

A flexible and robust plugin management system for Go applications.

## Features

- Dynamic loading and unloading of plugins
- Plugin versioning and dependency management
- Hot-reloading of plugins
- Event system for plugin lifecycle events
- Basic sandboxing for improved security
- Metrics collection for plugin performance

## Installation

To use this plugin manager in your Go project, run:

```bash
go get github.com/matt-dunleavy/plugin-manager
```

## Usage

### Creating a plugin

Plugins must implement the `Plugin` interface:

```go
package main

import (
    "fmt"
    pm "github.com/matt-dunleavy/plugin-manager"
)

type MyPlugin struct{}

func (p *MyPlugin) Metadata() pm.PluginMetadata {
    return pm.PluginMetadata{
        Name:         "MyPlugin",
        Version:      "1.0.0",
        Dependencies: []string{},
    }
}

func (p *MyPlugin) Init() error {
    fmt.Println("MyPlugin initialized")
    return nil
}

func (p *MyPlugin) Execute() error {
    fmt.Println("MyPlugin executed")
    return nil
}

func (p *MyPlugin) Shutdown() error {
    fmt.Println("MyPlugin shut down")
    return nil
}

var Plugin MyPlugin
```

### Compiling Plugins

Compile your plugin with:

```bash
go build -buildmode=plugin -o myplugin.so myplugin.go
```

### Using the Plugin Manager

Here's an example of how to use the plugin manager in your application:

```go
package main

import (
    "fmt"
    "log"
    pm "github.com/matt-dunleavy/plugin-manager"
)

func main() {
    
    // Create a new plugin manager
    manager, err := pm.NewManager("plugins.json", "./plugins")
    if err != nil {
        log.Fatalf("Failed to create plugin manager: %v", err)
    }

    // Load a plugin
    err = manager.LoadPlugin("./plugins/myplugin.so")
    if err != nil {
        log.Fatalf("Failed to load plugin: %v", err)
    }

    // Execute a plugin
    err = manager.ExecutePlugin("MyPlugin")
    if err != nil {
        log.Fatalf("Failed to execute plugin: %v", err)
    }

    // Hot-reload a plugin
    err = manager.HotReload("MyPlugin", "./plugins/myplugin_v2.so")
    if err != nil {
        log.Fatalf("Failed to hot-reload plugin: %v", err)
    }

    // Unload a plugin
    err = manager.UnloadPlugin("MyPlugin")
    if err != nil {
        log.Fatalf("Failed to unload plugin: %v", err)
    }

    // Subscribe to plugin events
    manager.GetEventBus().Subscribe("PluginLoaded", func(e pm.Event) {
        fmt.Printf("Plugin loaded: %s\n", e.(pm.PluginLoadedEvent).PluginName)
    })
}
```

## Configuration

The plugin manager uses a JSON configuration file to keep track of enabled plugins. Here's an example `plugins.json`:

```json
{
  "enabled": {
    "MyPlugin": true,
    "AnotherPlugin": false
  }
}
```

## API Reference

- ### Manager

  - `NewManager(configPath string, pluginDir string) (*Manager, error)`
  - `LoadPlugin(path string) error`
  - `UnloadPlugin(name string) error`
  - `ExecutePlugin(name string) error`
  - `HotReload(name string, path string) error`
  - `EnablePlugin(name string) error`
  - `DisablePlugin(name string) error`
  - `LoadEnabledPlugins(pluginDir string) error`
  - `ListPlugins() []string`
  - `GetPluginStats(name string) (*PluginStats, error)`
  - `SubscribeToEvent(eventName string, handler EventHandler)`

### EventBus

- `Subscribe(eventName string, handler EventHandler)`
- `Publish(event Event)`

### Sandbox

- `Enable() error`
- `Disable() error`
- `VerifyPluginPath(path string) error`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.