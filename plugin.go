package pluginmanager

import (
    "plugin"
    "time"
)

type PluginMetadata struct {
    Name         string
    Version      string
    Dependencies []string
}

type Plugin interface {
    Metadata() PluginMetadata
    Init() error
    Execute() error
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
        return nil, err
    }

    symPlugin, err := p.Lookup(PluginSymbol)
    if err != nil {
        return nil, err
    }

    plugin, ok := symPlugin.(Plugin)
    if !ok {
        return nil, ErrInvalidPluginInterface
    }

    return plugin, nil
}