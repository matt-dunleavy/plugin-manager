// Copyright (C) 2024 Matt Dunleavy. All rights reserved.
// Use of this source code is subject to the MIT license 
// that can be found in the LICENSE file.

package pluginmanager

import (
    "errors"
    "fmt"
)

var (
    ErrPluginAlreadyLoaded    = errors.New("plugin already loaded")
    ErrInvalidPluginInterface = errors.New("invalid plugin interface")
    ErrPluginNotFound         = errors.New("plugin not found")
    ErrIncompatibleVersion    = errors.New("incompatible plugin version")
    ErrMissingDependency      = errors.New("missing plugin dependency")
    ErrCircularDependency     = errors.New("circular plugin dependency detected")
    ErrPluginSandboxViolation = errors.New("plugin attempted to violate sandbox")
)

type PluginError struct {
    Op     string
    Err    error
    Plugin string
}

func (e *PluginError) Error() string {
    if e.Plugin != "" {
        return fmt.Sprintf("plugin error: %s: %s: %v", e.Plugin, e.Op, e.Err)
    }
    return fmt.Sprintf("plugin error: %s: %v", e.Op, e.Err)
}

func (e *PluginError) Unwrap() error {
    return e.Err
}