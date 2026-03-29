package catalog

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// LoadRegistry loads additional catalog entries from a JSON file.
// The path can be a local file path or an HTTP(S) URL.
func LoadRegistry(path string) (*Catalog, error) {
	var data []byte
	var err error

	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		resp, fetchErr := http.Get(path) //nolint:gosec // User-provided URL is intentional.
		if fetchErr != nil {
			return nil, fmt.Errorf("fetch registry %s: %w", path, fetchErr)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("fetch registry %s: HTTP %d", path, resp.StatusCode)
		}
		data, err = io.ReadAll(resp.Body)
	} else {
		data, err = os.ReadFile(path)
	}
	if err != nil {
		return nil, fmt.Errorf("read registry %s: %w", path, err)
	}

	var cat Catalog
	if err := json.Unmarshal(data, &cat); err != nil {
		return nil, fmt.Errorf("parse registry %s: %w", path, err)
	}
	return &cat, nil
}

// MergeCatalogs combines multiple catalogs. Later entries override earlier ones
// when they share the same skill ID. Bundle lists are concatenated (no dedup).
func MergeCatalogs(catalogs ...*Catalog) *Catalog {
	seen := make(map[string]int) // skill ID → index in result.
	var merged Catalog

	for _, cat := range catalogs {
		if cat == nil {
			continue
		}
		for _, entry := range cat.Skills {
			if idx, ok := seen[entry.ID]; ok {
				merged.Skills[idx] = entry // override.
			} else {
				seen[entry.ID] = len(merged.Skills)
				merged.Skills = append(merged.Skills, entry)
			}
		}
		merged.Bundles = append(merged.Bundles, cat.Bundles...)
	}
	return &merged
}
