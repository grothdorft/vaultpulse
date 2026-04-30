package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the full vaultpulse configuration.
type Config struct {
	Vault    VaultConfig    `yaml:"vault"`
	Alerts   AlertsConfig   `yaml:"alerts"`
	Webhooks []WebhookConfig `yaml:"webhooks"`
}

// VaultConfig contains Vault connection settings.
type VaultConfig struct {
	Address   string        `yaml:"address"`
	Token     string        `yaml:"token"`
	Namespace string        `yaml:"namespace"`
	PollInterval time.Duration `yaml:"poll_interval"`
}

// AlertsConfig defines thresholds for lease expiry warnings.
type AlertsConfig struct {
	WarnBefore  time.Duration `yaml:"warn_before"`
	CritBefore  time.Duration `yaml:"crit_before"`
}

// WebhookConfig describes a single webhook destination.
type WebhookConfig struct {
	Name    string            `yaml:"name"`
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers"`
	Timeout time.Duration     `yaml:"timeout"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	cfg.applyDefaults()
	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Vault.Address == "" {
		return fmt.Errorf("vault.address is required")
	}
	if c.Vault.Token == "" {
		return fmt.Errorf("vault.token is required")
	}
	for i, wh := range c.Webhooks {
		if wh.URL == "" {
			return fmt.Errorf("webhooks[%d].url is required", i)
		}
	}
	return nil
}

func (c *Config) applyDefaults() {
	if c.Vault.PollInterval == 0 {
		c.Vault.PollInterval = 60 * time.Second
	}
	if c.Alerts.WarnBefore == 0 {
		c.Alerts.WarnBefore = 24 * time.Hour
	}
	if c.Alerts.CritBefore == 0 {
		c.Alerts.CritBefore = 1 * time.Hour
	}
	for i := range c.Webhooks {
		if c.Webhooks[i].Timeout == 0 {
			c.Webhooks[i].Timeout = 10 * time.Second
		}
	}
}
