// Copyright (C) 2024 Matt Dunleavy. All rights reserved.
// Use of this source code is subject to the MIT license 
// that can be found in the LICENSE file.

package pluginmanager

import (
    "plugin"
    "time"
)

type PluginMetadata struct {
    Name         string
    Version      string
    Dependencies map[string]string
    GoVersion    string
    Signature    []byte
}

type Plugin interface {
    Metadata() PluginMetadata
    PreLoad() error
    Init() error
    PostLoad() error
    Execute() error
    PreUnload() error
    Shutdown() error
}

type PluginStats struct {
    ExecutionCount     int64
    LastExecutionTime  time.Duration
    TotalExecutionTime time.Duration
}

const PluginSymbol = "Plugin"

func LoadPlugin(path string) (Plugin, error) {
    p, err := plugin.Open(path)
    if err != nil {
        return nil, &PluginError{Op: "open", Err: err}
    }

    symPlugin, err := p.Lookup(PluginSymbol)
    if err != nil {
        return nil, &PluginError{Op: "lookup", Err: err}
    }

    plugin, ok := symPlugin.(Plugin)
    if !ok {
        return nil, &PluginError{Op: "assert", Err: ErrInvalidPluginInterface}
    }

    return plugin, nil
}