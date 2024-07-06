package pluginmanager

import (
    "os"
    "path/filepath"
    "log"
)

type Sandbox interface {
    Enable() error
    Disable() error
    VerifyPluginPath(path string) error
}

type DefaultSandbox struct {
    pluginDir string
}

func NewDefaultSandbox(pluginDir string) *DefaultSandbox {
    absPath, err := filepath.Abs(pluginDir)
    if err != nil {
        log.Printf("Error getting absolute path for plugin directory: %v", err)
        return &DefaultSandbox{pluginDir: pluginDir}
    }
    return &DefaultSandbox{
        pluginDir: absPath,
    }
}

func (s *DefaultSandbox) Enable() error {
    return os.Chdir(s.pluginDir)
}

func (s *DefaultSandbox) Disable() error {
    return nil
}

func (s *DefaultSandbox) VerifyPluginPath(path string) error {
    absPath, err := filepath.Abs(path)
    if err != nil {
        return err
    }
    
    if !filepath.HasPrefix(absPath, s.pluginDir) {
        return ErrPluginSandboxViolation
    }
    
    return nil
}