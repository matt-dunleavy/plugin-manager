# Changelog

## [1.3.0] - 2024-07-06

### Added
- Enhanced plugin lifecycle management
  - Implemented `PreLoad`, `PostLoad`, and `PreUnload` hooks in the Plugin interface
  - Added graceful shutdown for plugins during hot-reload
- Improved version compatibility checking
  - Implemented custom version comparison logic without external dependencies
  - Added support for version constraints (>=, >, <=, <, ==)
- Lazy loading for plugins
  - Introduced `lazyPlugin` struct for deferred plugin loading
  - Optimized plugin loading to occur only when necessary
- Comprehensive error handling
  - Added detailed error messages throughout the codebase
  - Improved error propagation and logging
- Enhanced concurrency safety
  - Implemented fine-grained locking with `sync.RWMutex`
  - Ensured thread-safety for all shared data access operations
- Plugin statistics tracking
  - Added `PluginStats` struct to track execution count and times
  - Implemented `GetPluginStats` method for retrieving plugin performance data

### Changed
- Refactored `Manager` struct to support new features
- Updated `LoadPlugin` method to incorporate new lifecycle hooks and lazy loading
- Modified `HotReload` method to include graceful shutdown of old plugin versions
- Improved `ExecutePlugin` method to work with lazy-loaded plugins and update statistics
- Enhanced `checkDependency` method to use the new version compatibility checking

### Removed
- Dependency on external version comparison libraries

### Security
- Added placeholder for plugin signature verification in `LoadPlugin` and `HotReload` methods
- Implemented `VerifyPluginSignature` method (to be fully implemented in future)

## [1.2.0] - 2024-07-05

### Added
- Plugin discovery system
  - Implemented automatic plugin discovery in specified directories
- Remote plugin repository support
  - Added `PluginRepository` struct for managing remote repositories
  - Implemented `SetupRemoteRepository` function for configuring remote repositories
- Plugin update system
  - Added `CheckForUpdates` and `UpdatePlugin` functions (placeholders to be implemented)
- Redbean server integration for plugin repositories
  - Implemented `DeployRepository` function for easy deployment of plugin repositories
  - Added automatic download and setup of Redbean server
- Digital signature verification for plugins
  - Added `VerifyPluginSignature` method (placeholder to be fully implemented)
- Enhanced plugin lifecycle hooks
  - Added `PreLoad` and `PostLoad` hooks to the Plugin interface
  - Implemented `PreUnload` hook in addition to existing `Shutdown` method
- SSH key support for remote repositories
  - Integrated SSH public key authentication for secure remote repository access

### Changed
- Updated `PluginMetadata` struct to include `Signature` field
- Modified `Manager` struct to include `publicKey` field for signature verification
- Enhanced `LoadPlugin` method to utilize new lifecycle hooks and signature verification
- Updated `UnloadPlugin` method to use the new `PreUnload` hook
- Refactored `HotReload` method to incorporate new plugin lifecycle

### Fixed
- Resolved issues with unused imports in discovery.go
- Fixed type mismatch in metadata.Dependencies assignment in manager.go
- Corrected public key type handling in SetupRemoteRepository function

## [1.1.0] - 2024-07-04

### Added
- Improved dependency management system
  - Added version compatibility checks for plugin dependencies
  - Implemented `checkDependencies` method in Manager
- Toolchain compatibility checks
  - Added Go version compatibility check during plugin loading
- Enhanced sandboxing
  - Implemented `LinuxSandbox` with chroot functionality
  - Added `VerifyPluginPath` method to prevent loading plugins from outside the sandbox
- Lazy loading of plugins
  - Introduced `lazyPlugin` struct for deferred plugin loading
- Comprehensive logging system
  - Integrated zap logger for structured logging
- Enhanced error handling
  - Created custom `PluginError` type for more informative error messages
- Improved hot-reloading mechanism
  - Added support for graceful shutdown of old plugin versions
- Performance optimizations
  - Optimized plugin loading and execution paths

### Changed
- Refactored `Manager` struct to support new features
- Updated `PluginMetadata` to include dependency information and Go version
- Modified `LoadPlugin` method to support lazy loading
- Updated `ExecutePlugin` method to work with lazy-loaded plugins
- Replaced `GetEventBus` method with `SubscribeToEvent` for better encapsulation

### Fixed
- Resolved potential race conditions in plugin management operations
- Improved error handling and reporting across all operations

## [1.0.0] - 2024-07-01

### Added
- Initial release of the plugin manager
- Basic plugin loading and unloading functionality
- Simple event system for plugin lifecycle events
- Configuration management for enabled/disabled plugins
- Basic plugin execution and stats collection