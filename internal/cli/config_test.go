package cli_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dirien/yet-another-agent-harness/internal/cli"
)

func TestFindConfig_Found(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "yaah.json")
	if err := os.WriteFile(cfgPath, []byte(`{"version":"1"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	found, err := cli.FindConfig(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found != cfgPath {
		t.Errorf("got %q, want %q", found, cfgPath)
	}
}

func TestFindConfig_WalksUp(t *testing.T) {
	parent := t.TempDir()
	child := filepath.Join(parent, "sub", "dir")
	if err := os.MkdirAll(child, 0o755); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(parent, "yaah.json")
	if err := os.WriteFile(cfgPath, []byte(`{"version":"1"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	found, err := cli.FindConfig(child)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found != cfgPath {
		t.Errorf("got %q, want %q", found, cfgPath)
	}
}

func TestFindConfig_NotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := cli.FindConfig(dir)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLoadConfig_Valid(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "yaah.json")
	if err := os.WriteFile(cfgPath, []byte(`{"version":"1","settings":{"model":"opus"}}`), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := cli.LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Version != "1" {
		t.Errorf("Version: got %q, want %q", cfg.Version, "1")
	}
	if cfg.Settings == nil || cfg.Settings.Model != "opus" {
		t.Errorf("expected model opus, got %v", cfg.Settings)
	}
}

func TestLoadConfig_DefaultVersion(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "yaah.json")
	if err := os.WriteFile(cfgPath, []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := cli.LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Version != "1" {
		t.Errorf("Version: got %q, want %q (expected default)", cfg.Version, "1")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "yaah.json")
	if err := os.WriteFile(cfgPath, []byte(`{broken`), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := cli.LoadConfig(cfgPath)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := cli.LoadConfig("/nonexistent/yaah.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
