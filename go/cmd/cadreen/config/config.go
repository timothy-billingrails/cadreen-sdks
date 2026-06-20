package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	APIKey      string `json:"api_key,omitempty"`
	BaseURL     string `json:"base_url,omitempty"`
	WorkspaceID string `json:"workspace_id,omitempty"`
	Profile     string `json:"profile,omitempty"`
	Output      string `json:"output,omitempty"`
}

const (
	DefaultBaseURL = "https://accomplishanything.today"
	DefaultProfile = "full"
	DefaultOutput  = "text"
	EnvAPIKey      = "CADREEN_API_KEY"
	EnvBaseURL     = "CADREEN_BASE_URL"
)

func ConfigDir() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "cadreen")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(os.TempDir(), ".cadreen")
	}
	return filepath.Join(home, ".cadreen")
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.json")
}

func Load() (*Config, error) {
	return LoadFrom(ConfigPath())
}

func LoadFrom(path string) (*Config, error) {
	cfg := &Config{
		BaseURL: DefaultBaseURL,
		Profile: DefaultProfile,
		Output:  DefaultOutput,
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}

func (c *Config) Save() error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	path := ConfigPath()
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	return nil
}

func (c *Config) Delete() error {
	path := ConfigPath()
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("deleting config: %w", err)
	}
	return nil
}

func (c *Config) APIKeyResolved() string {
	if c.APIKey != "" {
		return c.APIKey
	}
	return os.Getenv(EnvAPIKey)
}

func (c *Config) BaseURLResolved() string {
	if env := os.Getenv(EnvBaseURL); env != "" {
		return env
	}
	if c.BaseURL != "" {
		return c.BaseURL
	}
	return DefaultBaseURL
}

func (c *Config) IsAuthenticated() bool {
	return c.APIKeyResolved() != ""
}

func (c *Config) Set(key, value string) error {
	switch key {
	case "api_key":
		c.APIKey = value
	case "base_url":
		c.BaseURL = value
	case "workspace_id":
		c.WorkspaceID = value
	case "profile":
		c.Profile = value
	case "output":
		c.Output = value
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}
	return c.Save()
}

func (c *Config) Get(key string) string {
	switch key {
	case "api_key":
		return c.APIKey
	case "base_url":
		return c.BaseURL
	case "workspace_id":
		return c.WorkspaceID
	case "profile":
		return c.Profile
	case "output":
		return c.Output
	default:
		return ""
	}
}
