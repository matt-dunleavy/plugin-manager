// Copyright (C) 2024 Matt Dunleavy. All rights reserved.
// Use of this source code is subject to the MIT license 
// that can be found in the LICENSE file.

package pluginmanager

import (
    "crypto"
	"crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

type PluginRepository struct {
    URL       string
    SSHKey    string
    PublicKey ssh.PublicKey
}

func (m *Manager) DiscoverPlugins(dir string) error {
    return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if filepath.Ext(path) == ".so" {
            pluginName := strings.TrimSuffix(filepath.Base(path), ".so")
            if err := m.LoadPlugin(path); err != nil {
                m.logger.Warn("Failed to load discovered plugin", zap.String("plugin", pluginName), zap.Error(err))
            } else {
                m.logger.Info("Discovered and loaded plugin", zap.String("plugin", pluginName))
            }
        }
        return nil
    })
}

func (m *Manager) SetupRemoteRepository(url, sshKeyPath string) (*PluginRepository, error) {
    key, err := os.ReadFile(sshKeyPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read SSH key: %w", err)
    }

    signer, err := ssh.ParsePrivateKey(key)
    if err != nil {
        return nil, fmt.Errorf("failed to parse SSH key: %w", err)
    }

    return &PluginRepository{
        URL:       url,
        SSHKey:    string(key),
        PublicKey: signer.PublicKey(),
    }, nil
}

func (m *Manager) DeployRepository(repo *PluginRepository, localPath string) error {
    if err := m.downloadRedbean(localPath); err != nil {
        return err
    }

    cmd := exec.Command(filepath.Join(localPath, "redbean.com"), "-v")
    if repo.URL != "" {
        // Deploy via SSH
        cmd = exec.Command("ssh", "-i", repo.SSHKey, repo.URL, filepath.Join(localPath, "redbean.com"), "-v")
    }

    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("failed to deploy repository: %w\nOutput: %s", err, string(output))
    }

    m.logger.Info("Repository deployed successfully", zap.String("output", string(output)))
    return nil
}

func (m *Manager) downloadRedbean(localPath string) error {
    resp, err := http.Get("https://redbean.dev/redbean-latest.com")
    if err != nil {
        return fmt.Errorf("failed to download redbean: %w", err)
    }
    defer resp.Body.Close()

    out, err := os.Create(filepath.Join(localPath, "redbean.com"))
    if err != nil {
        return fmt.Errorf("failed to create redbean file: %w", err)
    }
    defer out.Close()

    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return fmt.Errorf("failed to save redbean file: %w", err)
    }

    if runtime.GOOS != "windows" {
        if err := os.Chmod(filepath.Join(localPath, "redbean.com"), 0755); err != nil {
            return fmt.Errorf("failed to set execute permission on redbean: %w", err)
        }
    }

    return nil
}

func (m *Manager) CheckForUpdates(repo *PluginRepository) ([]string, error) {
    // Implement logic to check for updates from the repository
    // This would typically involve making an HTTP request to the repository
    // and comparing versions of installed plugins with available versions
    return []string{}, nil
}

func (m *Manager) UpdatePlugin(repo *PluginRepository, pluginName string) error {
    // Implement logic to download and update a specific plugin
    return nil
}

func (m *Manager) VerifyPluginSignature(pluginPath string, publicKeyPath string) error {
    // Read the plugin file
    pluginData, err := os.ReadFile(pluginPath)
    if err != nil {
        return fmt.Errorf("failed to read plugin file: %w", err)
    }

    // Read the signature file
    signaturePath := pluginPath + ".sig"
    signatureData, err := os.ReadFile(signaturePath)
    if err != nil {
        return fmt.Errorf("failed to read signature file: %w", err)
    }

    // Read the public key
    publicKeyData, err := os.ReadFile(publicKeyPath)
    if err != nil {
        return fmt.Errorf("failed to read public key file: %w", err)
    }

    // Parse the public key
    block, _ := pem.Decode(publicKeyData)
    if block == nil {
        return fmt.Errorf("failed to parse PEM block containing the public key")
    }

    publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return fmt.Errorf("failed to parse public key: %w", err)
    }

    rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
    if !ok {
        return fmt.Errorf("public key is not an RSA public key")
    }

    // Verify the signature
    hashed := sha256.Sum256(pluginData)
    err = rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, hashed[:], signatureData)
    if err != nil {
        return fmt.Errorf("failed to verify signature: %w", err)
    }

    return nil
}