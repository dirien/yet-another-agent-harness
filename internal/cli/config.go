package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dirien/yet-another-agent-harness/pkg/schema"
)

const configFileName = "yaah.json"

// FindConfig walks up from dir looking for yaah.json.
func FindConfig(dir string) (string, error) {
	for {
		candidate := filepath.Join(dir, configFileName)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("%s not found in any parent directory", configFileName)
}

// LoadConfig reads and validates a yaah config file.
func LoadConfig(path string) (*schema.HarnessConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg schema.HarnessConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if cfg.Version == "" {
		cfg.Version = "1"
	}
	return &cfg, nil
}
