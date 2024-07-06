package main

import (
    "fmt"

    pm "github.com/matt-dunleavy/plugin-manager"
)

type HelloPlugin struct{}

func (p *HelloPlugin) Metadata() pm.PluginMetadata {
    return pm.PluginMetadata{
        Name:    "HelloPlugin",
        Version: "1.0.0",
        Dependencies: map[string]string{},
    }
}

func (p *HelloPlugin) Init() error {
    fmt.Println("HelloPlugin initialized")
    return nil
}

func (p *HelloPlugin) Execute() error {
    fmt.Println("Hello from HelloPlugin!")
    return nil
}

func (p *HelloPlugin) Shutdown() error {
    fmt.Println("HelloPlugin shut down")
    return nil
}

func (p *HelloPlugin) PreLoad() error {
    fmt.Println("HelloPlugin pre-load")
    return nil
}

func (p *HelloPlugin) PostLoad() error {
    fmt.Println("HelloPlugin post-load")
    return nil
}

func (p *HelloPlugin) PreUnload() error {
    fmt.Println("HelloPlugin pre-unload")
    return nil
}

var Plugin HelloPlugin