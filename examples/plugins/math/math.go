package main

import (
    "fmt"

    pm "github.com/matt-dunleavy/plugin-manager"
)

type MathPlugin struct{}

func (p *MathPlugin) Metadata() pm.PluginMetadata {
    return pm.PluginMetadata{
        Name:    "MathPlugin",
        Version: "1.0.0",
        Dependencies: map[string]string{},
    }
}

func (p *MathPlugin) Init() error {
    fmt.Println("MathPlugin initialized")
    return nil
}

func (p *MathPlugin) Execute() error {
    result := p.Add(5, 3)
    fmt.Printf("MathPlugin: 5 + 3 = %d\n", result)
    return nil
}

func (p *MathPlugin) Shutdown() error {
    fmt.Println("MathPlugin shut down")
    return nil
}

func (p *MathPlugin) PreLoad() error {
    fmt.Println("MathPlugin pre-load")
    return nil
}

func (p *MathPlugin) PostLoad() error {
    fmt.Println("MathPlugin post-load")
    return nil
}

func (p *MathPlugin) PreUnload() error {
    fmt.Println("MathPlugin pre-unload")
    return nil
}

func (p *MathPlugin) Add(a, b int) int {
    return a + b
}

var Plugin MathPlugin