package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigOptional_CodexHeaderDefaults(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	configYAML := []byte(`
codex-header-defaults:
  user-agent: "  my-codex-client/1.0  "
  beta-features: "  feature-a,feature-b  "
`)
	if err := os.WriteFile(configPath, configYAML, 0o600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := LoadConfigOptional(configPath, false)
	if err != nil {
		t.Fatalf("LoadConfigOptional() error = %v", err)
	}

	if got := cfg.CodexHeaderDefaults.UserAgent; got != "my-codex-client/1.0" {
		t.Fatalf("UserAgent = %q, want %q", got, "my-codex-client/1.0")
	}
	if got := cfg.CodexHeaderDefaults.BetaFeatures; got != "feature-a,feature-b" {
		t.Fatalf("BetaFeatures = %q, want %q", got, "feature-a,feature-b")
	}
}

func TestLoadConfigOptional_CodexIdentityConfuse(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	configYAML := []byte(`
codex:
  identity-confuse: true
`)
	if err := os.WriteFile(configPath, configYAML, 0o600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := LoadConfigOptional(configPath, false)
	if err != nil {
		t.Fatalf("LoadConfigOptional() error = %v", err)
	}

	if !cfg.Codex.IdentityConfuse {
		t.Fatalf("IdentityConfuse = false, want true")
	}
}

func TestLoadConfigOptional_CodexPlanProxy(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	configYAML := []byte(`
codex:
  plan-proxy:
    free: " socks5h://127.0.0.1:10810 "
    plus-team: " socks5h://127.0.0.1:10811 "
    pro: " socks5h://127.0.0.1:10812 "
`)
	if err := os.WriteFile(configPath, configYAML, 0o600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := LoadConfigOptional(configPath, false)
	if err != nil {
		t.Fatalf("LoadConfigOptional() error = %v", err)
	}

	if got := cfg.Codex.PlanProxy.ProxyURLForPlanType("free"); got != "socks5h://127.0.0.1:10810" {
		t.Fatalf("free proxy = %q", got)
	}
	for _, planType := range []string{"plus", "team", "business", "go"} {
		if got := cfg.Codex.PlanProxy.ProxyURLForPlanType(planType); got != "socks5h://127.0.0.1:10811" {
			t.Fatalf("%s proxy = %q", planType, got)
		}
	}
	if got := cfg.Codex.PlanProxy.ProxyURLForPlanType("pro"); got != "socks5h://127.0.0.1:10812" {
		t.Fatalf("pro proxy = %q", got)
	}
	if got := cfg.Codex.PlanProxy.ProxyURLForPlanType("unknown"); got != "" {
		t.Fatalf("unknown proxy = %q", got)
	}
}
