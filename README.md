# Plugin Manager for Go

A robust and flexible plugin management library for Go applications.

##### Features

- Dynamic Loading and Unloading
- Versioning and Dependency Management
- Hot-Reloading and Lazy Loading
- Event System for Lifecycle Events
- Enhanced Sandboxing for improved security

- Metrics Collection for Plugin Performance

##### Deployment and Updates Simplified

- Automatic Discovery (directory)
- Remote Plugin Repository
- Automated Plugin Update System
- Digital Signature Verification

Easy deployment of your plugin repository with Redbean!

## Getting Started

Obtain the latest version of the `plugin-manager` library with the Go package manager (recommended):

```bash
go get github.com/matt-dunleavy/plugin-manager
```

Or clone the repository to your local machine using the Git Command Line by running the following command in your terminal:

```sh
git clone https://github.com/matt-dunleavy/plugin-manager.git
cd plugin-manager
```

After you've obtained the library, import the package into your source code:

```go
import "github.com/matt-dunleavy/plugin-manager"
```

## Usage

Visit to the [examples]() directory for full-featured implementations.

#### Initialize a New Plugin Manager

Create a new plugin manager instance at the start of your application.

```go
manager, err := pm.NewManager("plugins.json", "./plugins", "public_key.pem")
```

**Parameters:**

- `configPath` (string): Path to the JSON configuration file for managing enabled/disabled plugins. ("plugins.json")
- `pluginDir` (string): Directory where plugins are stored. ("./plugins")
- `publicKeyPath` (string): Path to the public key file used for verifying plugin signatures. ("public_key.pem")

**Returns:**

- `*Manager`: Pointer to the newly created Manager instance.
- `error`: Any error encountered during initialization.

#### Load a Plugin

Load a plugin from the specified path into memory, making it available for execution.

```go
err = manager.LoadPlugin("./plugins/myplugin.so")
```

**Parameters:**

- `path` (string): Path to the plugin file (.so extension).

**Returns:**

- `error`: Any error encountered during the loading process.

#### Execute (Run) a Plugin

Run a loaded plugin's Execute() method.

```go
err = manager.ExecutePlugin("MyPlugin")
```

**Parameters:**

- `name` (string): Name of the plugin to execute.

**Returns:**

- `error`: Any error encountered during plugin execution.

#### Unload a Plugin

Safely remove a plugin from memory when it's no longer needed.

```go
err = manager.UnloadPlugin("MyPlugin")
```

**Parameters:**

- `name` (string): Name of the plugin to unload.

**Returns:**

- `error`: Any error encountered during the unloading process.

#### Hot-Reload a Plugin

Update a loaded plugin to a new version while the application is running (without stopping the application).

```go
err = manager.HotReload("MyPlugin", "./plugins/myplugin_v2.so")
```

**Parameters:**

- `name` (string): Name of the plugin to hot-reload.
- `path` (string): Path to the new version of the plugin.

**Returns:**

- `error`: Any error encountered during the hot-reload process.

#### **Enable Automatic Plugin Discovery**

Automatically discover and load all plugins from a specified directory.

```go
err = manager.DiscoverPlugins("./plugins")
```

**Parameters:**

- `dir` (string): Directory to search for plugins.

**Returns:**

- `error`: Any error encountered during the discovery process.

#### Subscribe to Plugin Events

Subscribes to a specific plugin event, executing the provided function when the event occurs. Use this function to set up event handlers for various plugin lifecycle events.

```go
manager.SubscribeToEvent("PluginLoaded", func(e pm.Event) {
	fmt.Printf("Plugin loaded: %s\n", e.(pm.PluginLoadedEvent).PluginName)
})
```

**Parameters:**

- `eventName` (string): Name of the event to subscribe to (e.g., "PluginLoaded").
- `handler` (func(Event)): Function to execute when the event occurs.

**Returns:** None

## Creating a plugin

Plugins must implement the `Plugin` interface.

#### Plugin (struct)

The `Plugin` struct defines the basic structure of your plugin. It can be empty if your plugin doesn't need to store any state, or you can add fields if your plugin needs to maintain state across method calls.

```go
type MyPlugin struct{}
```

#### Metadata()

The `Metadata()` method returns metadata about the plugin, including the `Name`, `Version`, and `Dependencies` (a map of other plugins that this plugin depends on, with the key being the plugin's name, and the value being the version constraint).

```go
func (p *MyPlugin) Metadata() pm.PluginMetadata {
    return pm.PluginMetadata{
        Name:         "MyPlugin",
        Version:      "1.0.0",
        Dependencies: map[string]string{},
    }
}
```

#### Preload()

The `Preload()` method is called before the plugin is fully loaded. Use it for any setup that needs to happen before initialization.

```go
func (p *MyPlugin) PreLoad() error {
    fmt.Println("MyPlugin pre-load")
    return nil
}
```

#### Init()

The `Init()` method is called to initialize the plugin. Use it to set up any resources or state the plugin needs.

```go
func (p *MyPlugin) Init() error {
    fmt.Println("MyPlugin initialized")
    return nil
}
```

#### PostLoad()

The `PostLoad()` method is called after the plugin is fully loaded. Use it for any final setup steps.

```go
func (p *MyPlugin) PostLoad() error {
    fmt.Println("MyPlugin post-load")
    return nil
}
```

#### Execute()

`Execute()` is the main method of your plugin. It is called when the plugin manager executes your plugin.

```go
func (p *MyPlugin) Execute() error {
    fmt.Println("MyPlugin executed")
    return nil
}
```

#### PreUnload()

The `PreUnload()` method is called before the plugin is unloaded. Use it to prepare for shutdown.

```Go
func (p *MyPlugin) PreUnload() error {
    fmt.Println("MyPlugin pre-unload")
    return nil
}
```

#### Shutdown()

The `Shutdown()` method is called when the plugin is being unloaded. Use it to clean up any resources the plugin has allocated.

```Go
func (p *MyPlugin) Shutdown() error {
    fmt.Println("MyPlugin shut down")
    return nil
}
```

#### Plugin Variable

This variable is how the plugin manager discovers your plugin. It must be named `Plugin` and be of the type that implements the plugin interface.

```go
var Plugin MyPlugin
```

> [!IMPORTANT]
>
> Each of these methods should return an error if something goes wrong during their execution. Returning `nil` indicates successful completion.

> [!NOTE]
>
> When implementing your own plugin, you would replace the `fmt.Println` statements with your actual plugin logic. The `PreLoad`, `PostLoad`, `PreUnload`, and `Shutdown` methods allow you to manage the lifecycle of your plugin, while `Init` and `Execute` form the core functionality.

## Compiling Plugins

Compile a plugin using the standard Go compiler toolchain by setting the `-buildmode` flag to `plugin`:

```bash
go build -buildmode=plugin -o myplugin.so myplugin.go
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

## Simplified Deployment Plugin Repositories

<img src="assets/img/redbean.png" style="float:right"/>An efficient and straightforward way to deploy and manage remote plugin repositories.

### Benefits

1. **Simplified Management**: Manage your entire plugin repository with a single file.
2. **Reduced Dependencies**: No need for complex web server setups or databases.
3. **Easy Updates**: Update your plugin repository by simply replacing the Redbean executable.
4. **Scalability**: Redbean can handle repositories of various sizes efficiently.

### Advantages

- **Security**: Redbean includes built-in security features, reducing the attack surface of your plugin repository.
- **Performance**: As a compiled C application, Redbean offers excellent performance for serving plugin files.
- **Flexibility**: Easily customize your repository structure and access controls.
- **Low Overhead**: Minimal resource usage makes it suitable for various hosting environments.

### Features

- **Single-File Deployment**: Redbean combines the web server and your content into a single executable file, simplifying deployment and distribution.
- **Automatic Download**: The plugin manager can automatically download and set up the Redbean server.
- **Easy Repository Deployment**: Deploy your plugin repository with a single function call.
- **Cross-Platform Compatibility**: Redbean works on various platforms, including Linux, macOS, and Windows.
- **Lightweight**: Redbean has a small footprint, making it ideal for plugin repositories of all sizes.

### Step-by-Step Guide to Implementing and Deploying Redbean

#### **Setup a Remote Repository**

```go
repo, err := manager.SetupRemoteRepository("user@example.com:/path/to/repo", "/path/to/ssh/key")
```

#### Prepare Local Directory for Deployment

Create a local directory to store your plugins and repository structure:

```go
localRepoPath := "./repository"
```

#### **Add Plugins to the Local Repository.**

Copy or move your plugin files to the local repository directory.

#### **Deploy the Repository**

```go
err = manager.DeployRepository(repo, localRepoPath)
```
This step will:

- Download the latest version of Redbean (if not present)
- Package your plugins and repository structure into a Redbean executable
- Deploy the Redbean executable to your remote server (if a remote URL was provided)

#### **Verify Deployment**

If deployed remotely, SSH into your server and check that the Redbean executable is present and running:

  ```bash
  ssh user@example.com
  ls /path/to/repo/redbean.com
  ps aux | grep redbean
  ```

  #### **Access Your Plugin Repository**

 Your plugins are now accessible via HTTP/HTTPS. If Redbean is running on the default port, you can access your plugins at: `http://your-server-address:8080/plugins/`

#### **Update Repository**

To update your repository, simply repeat steps 4-5. The plugin manager will handle updating the Redbean executable and redeploying your changes.

### Advanced Configuration

When deploying your plugin repository via redbean, the plugin manager will include a `redbean.ini` file that you can use to customize your repository's server configuration.

```ini
[server]
port = 9000
addr = 127.0.0.1
```

Refer to the [redbean.ini reference](docs/redbean.md) for a comprehensive list of commands.

> [!IMPORTANT]
>
> - Always use HTTPS in production environments.
> - Implement proper access controls to restrict repository access.
> - Regularly update both your plugins and the Redbean executable to ensure you have the latest security patches.

## API Reference

##### Management

- `NewManager(configPath string, pluginDir string, publicKeyPath string) (*Manager, error)`
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

##### Automatic Discovery and Updates

- `DiscoverPlugins(dir string) error`
- `CheckForUpdates(repo *PluginRepository) ([]string, error)`
- `UpdatePlugin(repo *PluginRepository, pluginName string) error`

##### Remote Repository (via [redbean](https://redbean.dev/))

- `SetupRemoteRepository(url, sshKeyPath string) (*PluginRepository, error)`
- `DeployRepository(repo *PluginRepository, localPath string) error`

##### EventBus

- `Subscribe(eventName string, handler EventHandler)`
- `Publish(event Event)`

##### Sandbox

- `Enable() error`
- `Disable() error`
- `VerifyPluginPath(path string) error`

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
