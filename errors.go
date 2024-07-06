package pluginmanager

import "errors"

var (
    ErrPluginAlreadyLoaded    = errors.New("plugin already loaded")
    ErrInvalidPluginInterface = errors.New("invalid plugin interface")
    ErrPluginNotFound         = errors.New("plugin not found")
    ErrIncompatibleVersion    = errors.New("incompatible plugin version")
    ErrMissingDependency      = errors.New("missing plugin dependency")
    ErrCircularDependency     = errors.New("circular plugin dependency detected")
    ErrPluginSandboxViolation = errors.New("plugin attempted to violate sandbox")
)