package cliproxy

import (
	"context"
	"testing"

	coreauth "github.com/router-for-me/CLIProxyAPI/v7/sdk/cliproxy/auth"
	"github.com/router-for-me/CLIProxyAPI/v7/sdk/config"
)

func TestApplyCodexPlanProxyGroups(t *testing.T) {
	t.Parallel()

	svc := &Service{cfg: &config.Config{}}
	svc.cfg.Codex.PlanProxy.Free = "socks5h://127.0.0.1:10810"
	svc.cfg.Codex.PlanProxy.PlusTeam = "socks5h://127.0.0.1:10811"
	svc.cfg.Codex.PlanProxy.Pro = "socks5h://127.0.0.1:10812"

	tests := []struct {
		name     string
		planType string
		want     string
	}{
		{name: "free", planType: "free", want: "socks5h://127.0.0.1:10810"},
		{name: "plus", planType: "plus", want: "socks5h://127.0.0.1:10811"},
		{name: "team", planType: "team", want: "socks5h://127.0.0.1:10811"},
		{name: "business", planType: "business", want: "socks5h://127.0.0.1:10811"},
		{name: "go", planType: "go", want: "socks5h://127.0.0.1:10811"},
		{name: "pro", planType: "pro", want: "socks5h://127.0.0.1:10812"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			auth := &coreauth.Auth{
				ID:         "codex-" + tc.name,
				Provider:   "codex",
				Attributes: map[string]string{"plan_type": tc.planType},
			}
			if !svc.applyCodexPlanProxy(auth) {
				t.Fatalf("applyCodexPlanProxy() changed = false")
			}
			if auth.ProxyURL != tc.want {
				t.Fatalf("ProxyURL = %q, want %q", auth.ProxyURL, tc.want)
			}
			if auth.Attributes[codexPlanProxyDerivedAttribute] != "true" {
				t.Fatalf("derived marker = %q, want true", auth.Attributes[codexPlanProxyDerivedAttribute])
			}
		})
	}
}

func TestApplyCodexPlanProxyPreservesExplicitProxy(t *testing.T) {
	t.Parallel()

	svc := &Service{cfg: &config.Config{}}
	svc.cfg.Codex.PlanProxy.Pro = "socks5h://127.0.0.1:10812"
	auth := &coreauth.Auth{
		ID:         "codex-pro",
		Provider:   "codex",
		ProxyURL:   "socks5h://explicit.local:1080",
		Attributes: map[string]string{"plan_type": "pro"},
		Metadata:   map[string]any{"proxy_url": "socks5h://explicit.local:1080"},
	}

	if svc.applyCodexPlanProxy(auth) {
		t.Fatalf("applyCodexPlanProxy() changed = true")
	}
	if auth.ProxyURL != "socks5h://explicit.local:1080" {
		t.Fatalf("ProxyURL = %q", auth.ProxyURL)
	}
}

func TestApplyCodexPlanProxyRefreshesDerivedProxy(t *testing.T) {
	t.Parallel()

	svc := &Service{cfg: &config.Config{}}
	svc.cfg.Codex.PlanProxy.Free = "socks5h://127.0.0.1:10810"
	auth := &coreauth.Auth{
		ID:         "codex-free",
		Provider:   "codex",
		Attributes: map[string]string{"plan_type": "free"},
	}

	if !svc.applyCodexPlanProxy(auth) {
		t.Fatalf("initial apply changed = false")
	}
	svc.cfg.Codex.PlanProxy.Free = "socks5h://127.0.0.1:11810"
	if !svc.applyCodexPlanProxy(auth) {
		t.Fatalf("second apply changed = false")
	}
	if auth.ProxyURL != "socks5h://127.0.0.1:11810" {
		t.Fatalf("ProxyURL = %q", auth.ProxyURL)
	}
}

func TestRefreshCodexPlanProxyAuthsUpdatesManager(t *testing.T) {
	t.Parallel()

	manager := coreauth.NewManager(nil, nil, nil)
	auth := &coreauth.Auth{
		ID:         "codex-free",
		Provider:   "codex",
		Attributes: map[string]string{"plan_type": "free"},
	}
	if _, err := manager.Register(context.Background(), auth); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	svc := &Service{cfg: &config.Config{}, coreManager: manager}
	svc.cfg.Codex.PlanProxy.Free = "socks5h://127.0.0.1:10810"
	svc.refreshCodexPlanProxyAuths(context.Background())

	got, ok := manager.GetByID("codex-free")
	if !ok {
		t.Fatalf("codex-free auth missing")
	}
	if got.ProxyURL != "socks5h://127.0.0.1:10810" {
		t.Fatalf("ProxyURL = %q", got.ProxyURL)
	}
}
