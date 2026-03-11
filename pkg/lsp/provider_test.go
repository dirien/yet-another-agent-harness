package lsp_test

import (
	"testing"

	"github.com/dirien/yet-another-agent-harness/pkg/lsp"
	"github.com/dirien/yet-another-agent-harness/pkg/lsp/providers"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

func TestRegistry_RegisterAndProviders(t *testing.T) {
	r := lsp.NewRegistry()
	if len(r.Providers()) != 0 {
		t.Fatal("expected 0 providers on new registry")
	}

	r.Register(providers.NewGopls())
	if len(r.Providers()) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(r.Providers()))
	}
}

func TestRegistry_ToConfig(t *testing.T) {
	r := lsp.NewRegistry()
	r.Register(providers.NewGopls())
	r.Register(providers.NewPyright())

	cfg := r.ToConfig()
	if len(cfg.Servers) != 2 {
		t.Fatalf("expected 2 servers, got %d", len(cfg.Servers))
	}
	if cfg.Servers[0].ID != "gopls" {
		t.Errorf("first server ID: got %q, want %q", cfg.Servers[0].ID, "gopls")
	}
}

func TestRegistry_CheckAll(t *testing.T) {
	r := lsp.NewRegistry()
	r.Register(providers.NewCustom(
		schema.LSPServer{ID: "nonexistent", Command: []string{"this-binary-does-not-exist-12345"}},
		"install it",
	))

	results := r.CheckAll()
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Installed {
		t.Error("expected nonexistent binary to be not installed")
	}
	if results[0].InstallHint != "install it" {
		t.Errorf("InstallHint: got %q, want %q", results[0].InstallHint, "install it")
	}
}

func TestFormatCheckResults(t *testing.T) {
	results := []lsp.CheckResult{
		{Name: "gopls", Installed: true, BinaryPath: "/usr/bin/gopls"},
		{Name: "missing", Installed: false, InstallHint: "install me"},
	}
	out := lsp.FormatCheckResults(results)
	if out == "" {
		t.Fatal("FormatCheckResults returned empty string")
	}
}
