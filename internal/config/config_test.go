package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "vaultpulse-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	content := `
vault:
  address: "http://127.0.0.1:8200"
  token: "root"
  poll_interval: 30s
alerts:
  warn_before: 48h
  crit_before: 2h
webhooks:
  - name: slack
    url: "https://hooks.slack.com/services/test"
    timeout: 5s
`
	path := writeTempConfig(t, content)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Address != "http://127.0.0.1:8200" {
		t.Errorf("expected vault address, got %q", cfg.Vault.Address)
	}
	if cfg.Vault.PollInterval != 30*time.Second {
		t.Errorf("expected 30s poll interval, got %v", cfg.Vault.PollInterval)
	}
	if cfg.Alerts.WarnBefore != 48*time.Hour {
		t.Errorf("expected 48h warn_before, got %v", cfg.Alerts.WarnBefore)
	}
	if cfg.Alerts.CritBefore != 2*time.Hour {
		t.Errorf("expected 2h crit_before, got %v", cfg.Alerts.CritBefore)
	}
	if len(cfg.Webhooks) != 1 || cfg.Webhooks[0].Name != "slack" {
		t.Errorf("unexpected webhooks: %+v", cfg.Webhooks)
	}
	if cfg.Webhooks[0].Timeout != 5*time.Second {
		t.Errorf("expected 5s webhook timeout, got %v", cfg.Webhooks[0].Timeout)
	}
}

func TestLoad_Defaults(t *testing.T) {
	content := `
vault:
  address: "http://127.0.0.1:8200"
  token: "root"
webhooks:
  - url: "https://example.com/hook"
`
	path := writeTempConfig(t, content)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.PollInterval != 60*time.Second {
		t.Errorf("expected default 60s, got %v", cfg.Vault.PollInterval)
	}
	if cfg.Alerts.WarnBefore != 24*time.Hour {
		t.Errorf("expected default 24h, got %v", cfg.Alerts.WarnBefore)
	}
	if cfg.Webhooks[0].Timeout != 10*time.Second {
		t.Errorf("expected default 10s webhook timeout, got %v", cfg.Webhooks[0].Timeout)
	}
}

func TestLoad_MissingVaultAddress(t *testing.T) {
	content := `
vault:
  token: "root"
`
	path := writeTempConfig(t, content)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestLoad_MissingWebhookURL(t *testing.T) {
	content := `
vault:
  address: "http://127.0.0.1:8200"
  token: "root"
webhooks:
  - name: broken
`
	path := writeTempConfig(t, content)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing webhook URL")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	content := `vault: [invalid: yaml: content`
	path := writeTempConfig(t, content)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}
