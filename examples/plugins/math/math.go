package main

import (
    "fmt"
    pm "github.com/matt-dunleavy/plugin-manager"
)

type MathPlugin struct{}

func (p *MathPlugin) Metadata() pm.PluginMetadata {
    return pm.PluginMetadata{
        Name:         "MathPlugin",
        Version:      "1.0.0",
        Dependencies: []string{},
    }
}

func (p *MathPlugin) Init() error {
    fmt.Println("MathPlugin initialized")
    return nil
}

func (p *MathPlugin) Execute() error {
    a, b := 10, 5
    fmt.Printf("Addition: %d + %d = %d\n", a, b, a+b)
    fmt.Printf("Subtraction: %d - %d = %d\n", a, b, a-b)
    fmt.Printf("Multiplication: %d * %d = %d\n", a, b, a*b)
    fmt.Printf("Division: %d / %d = %d\n", a, b, a/b)
    return nil
}

func (p *MathPlugin) Shutdown() error {
    fmt.Println("MathPlugin shut down")
    return nil
}

var Plugin MathPlugin