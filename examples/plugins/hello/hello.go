package main

import (
    "fmt"
    pm "github.com/matt-dunleavy/plugin-manager"
)

type HelloPlugin struct {
    greeting string
}

func (p *HelloPlugin) Metadata() pm.PluginMetadata {
    return pm.PluginMetadata{
        Name:         "HelloPlugin",
        Version:      "1.0.0",
        Dependencies: []string{},
    }
}

func (p *HelloPlugin) Init() error {
    p.greeting = "Hello, World!"
    fmt.Println("HelloPlugin initialized")
    return nil
}

func (p *HelloPlugin) Execute() error {
    fmt.Println(p.greeting)
    return nil
}

func (p *HelloPlugin) Shutdown() error {
    fmt.Println("HelloPlugin shut down")
    return nil
}

var Plugin HelloPlugin