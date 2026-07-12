package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Version         string `json:"version"`
	RuntimeVersion  string `json:"runtime_version"`
	StandardVersion string `json:"standard_version"`
	InstallPath     string `json:"install_path"`
}

func Default() *Config {
	home, _ := os.UserHomeDir()
	return &Config{
		Version:         "1.0.0",
		RuntimeVersion:  "1.0.0",
		StandardVersion: "1.0.0",
		InstallPath:     filepath.Join(home, ".kdse"),
	}
}

func (c *Config) Save() error {
	if err := os.MkdirAll(c.InstallPath, 0755); err != nil {
		return err
	}
	cfgPath := filepath.Join(c.InstallPath, "config.json")
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cfgPath, data, 0644)
}

func Load() (*Config, error) {
	cfg := Default()
	cfgPath := filepath.Join(cfg.InstallPath, "config.json")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
