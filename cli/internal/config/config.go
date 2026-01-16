package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	APIBaseURL  string `json:"api_base_url"`
	Token       string `json:"token"`
	TokenPrefix string `json:"token_prefix"`
	Email       string `json:"email"`
	Plan        string `json:"plan"`
}

func GetConfigDir() string {
	var configDir string

	switch runtime.GOOS {
	case "windows":
		configDir = os.Getenv("APPDATA")
		if configDir == "" {
			configDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
		configDir = filepath.Join(configDir, "InstantTLS")
	default:
		home, _ := os.UserHomeDir()
		configDir = filepath.Join(home, ".config", "instanttls")
	}

	return configDir
}

func GetConfigPath() string {
	return filepath.Join(GetConfigDir(), "config.json")
}

func GetCertDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".instanttls")
}

func GetCADir() string {
	return filepath.Join(GetCertDir(), "ca")
}

func GetCertsDir() string {
	return filepath.Join(GetCertDir(), "certs")
}

func Load() (*Config, error) {
	configPath := GetConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	configDir := GetConfigDir()
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(GetConfigPath(), data, 0600)
}

func IsLoggedIn() bool {
	cfg, err := Load()
	if err != nil || cfg == nil {
		return false
	}
	return cfg.Token != ""
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil || cfg == nil {
		return &Config{}
	}
	return cfg
}
