package catalog

import (
	"strings"
	"sync"

	"github.com/dirien/yet-another-agent-harness/pkg/gitcache"
)

// ValidationResult holds the outcome of validating a catalog entry.
type ValidationResult struct {
	ID             string
	OK             bool
	Size           int
	Error          string
	HasFrontmatter bool
}

// ValidateEntry attempts to fetch and parse a skill, returning a validation result.
func ValidateEntry(entry CatalogEntry) ValidationResult {
	r := ValidationResult{ID: entry.ID}

	if entry.Uses == "builtin" {
		r.OK = true
		return r
	}

	content, _, err := gitcache.FetchFileWithExtras(entry.Uses, entry.Subpath, "skills", true)
	if err != nil {
		r.Error = err.Error()
		return r
	}

	if len(strings.TrimSpace(content)) == 0 {
		r.Error = "empty SKILL.md"
		return r
	}

	r.Size = len(content)
	r.HasFrontmatter = strings.HasPrefix(strings.TrimSpace(content), "---")
	r.OK = true
	return r
}

// ValidateAll validates all entries in a catalog concurrently.
// The concurrency parameter controls the number of parallel fetches.
func ValidateAll(cat *Catalog, concurrency int) []ValidationResult {
	if cat == nil || len(cat.Skills) == 0 {
		return nil
	}
	if concurrency < 1 {
		concurrency = 4
	}

	results := make([]ValidationResult, len(cat.Skills))
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for i, entry := range cat.Skills {
		wg.Add(1)
		go func(idx int, e CatalogEntry) {
			defer wg.Done()
			sem <- struct{}{}
			results[idx] = ValidateEntry(e)
			<-sem
		}(i, entry)
	}

	wg.Wait()
	return results
}
