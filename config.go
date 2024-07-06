package pluginmanager

import (
    "encoding/json"
    "os"
    "sync"
)

type Config struct {
    Enabled map[string]bool `json:"enabled"`
    path    string
    mu      sync.RWMutex
}

func LoadConfig(path string) (*Config, error) {
    config := &Config{
        Enabled: make(map[string]bool),
        path:    path,
    }

    file, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return config, nil
        }
        return nil, err
    }

    if err := json.Unmarshal(file, &config); err != nil {
        return nil, err
    }

    return config, nil
}

func (c *Config) Save() error {
    c.mu.RLock()
    defer c.mu.RUnlock()

    data, err := json.MarshalIndent(c, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(c.path, data, 0644)
}

func (c *Config) EnablePlugin(name string) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.Enabled[name] = true
    return nil
}

func (c *Config) DisablePlugin(name string) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.Enabled[name] = false
    return nil
}

func (c *Config) EnabledPlugins() []string {
    c.mu.RLock()
    defer c.mu.RUnlock()

    var enabled []string
    for name, isEnabled := range c.Enabled {
        if isEnabled {
            enabled = append(enabled, name)
        }
    }
    return enabled
}