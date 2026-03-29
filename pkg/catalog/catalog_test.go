package catalog

import (
	"testing"
)

func TestDefaultCatalog(t *testing.T) {
	cat := DefaultCatalog()
	if len(cat.Skills) == 0 {
		t.Fatal("expected skills in default catalog")
	}
	if len(cat.Bundles) == 0 {
		t.Fatal("expected bundles in default catalog")
	}

	// All bundle skill IDs must reference valid catalog entries.
	for _, b := range cat.Bundles {
		for _, id := range b.SkillIDs {
			if cat.ByID(id) == nil {
				t.Errorf("bundle %q references unknown skill %q", b.ID, id)
			}
		}
	}
}

func TestSearch(t *testing.T) {
	cat := DefaultCatalog()

	tests := []struct {
		query   string
		wantMin int
		wantID  string
	}{
		{"kubernetes", 1, "kubernetes-specialist"},
		{"pulumi", 5, "pulumi-best-practices"},
		{"rust", 3, "rust-best-practices"},
		{"commit", 1, "commit"},
		{"NONEXISTENT_QUERY_XYZ", 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			results := cat.Search(tt.query)
			if len(results) < tt.wantMin {
				t.Errorf("Search(%q) returned %d results, want >= %d", tt.query, len(results), tt.wantMin)
			}
			if tt.wantID != "" {
				found := false
				for _, r := range results {
					if r.ID == tt.wantID {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Search(%q) did not find %q", tt.query, tt.wantID)
				}
			}
		})
	}
}

func TestSearchEmpty(t *testing.T) {
	cat := DefaultCatalog()
	if results := cat.Search(""); results != nil {
		t.Errorf("Search(\"\") should return nil, got %d results", len(results))
	}

	var nilCat *Catalog
	if results := nilCat.Search("test"); results != nil {
		t.Errorf("nil catalog Search should return nil")
	}
}

func TestByCategory(t *testing.T) {
	cat := DefaultCatalog()

	langs := cat.ByCategory(CategoryLanguage)
	if len(langs) == 0 {
		t.Fatal("expected language skills")
	}
	for _, s := range langs {
		if s.Category != CategoryLanguage {
			t.Errorf("ByCategory(language) returned skill %q with category %q", s.ID, s.Category)
		}
	}
}

func TestByTag(t *testing.T) {
	cat := DefaultCatalog()

	results := cat.ByTag("pulumi")
	if len(results) == 0 {
		t.Fatal("expected skills tagged 'pulumi'")
	}
	for _, s := range results {
		found := false
		for _, tag := range s.Tags {
			if tag == "pulumi" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ByTag(pulumi) returned skill %q without 'pulumi' tag", s.ID)
		}
	}
}

func TestByID(t *testing.T) {
	cat := DefaultCatalog()

	e := cat.ByID("golang-pro")
	if e == nil {
		t.Fatal("expected to find golang-pro by ID")
	}
	if e.Name != "golang-pro" {
		t.Errorf("got name %q, want golang-pro", e.Name)
	}

	if cat.ByID("nonexistent") != nil {
		t.Error("expected nil for nonexistent ID")
	}
}

func TestByAlias(t *testing.T) {
	cat := DefaultCatalog()

	e := cat.ByAlias("go-pro")
	if e == nil {
		t.Fatal("expected to find golang-pro by alias go-pro")
	}
	if e.ID != "golang-pro" {
		t.Errorf("got ID %q, want golang-pro", e.ID)
	}

	e = cat.ByAlias("ts-pro")
	if e == nil {
		t.Fatal("expected to find typescript-pro by alias ts-pro")
	}
	if e.ID != "typescript-pro" {
		t.Errorf("got ID %q, want typescript-pro", e.ID)
	}
}

func TestResolve(t *testing.T) {
	cat := DefaultCatalog()

	// Resolve by ID.
	e := cat.Resolve("golang-pro")
	if e == nil || e.ID != "golang-pro" {
		t.Error("Resolve by ID failed")
	}

	// Resolve by alias.
	e = cat.Resolve("k8s")
	if e == nil || e.ID != "kubernetes-specialist" {
		t.Error("Resolve by alias 'k8s' failed")
	}

	// Resolve nonexistent.
	if cat.Resolve("nonexistent") != nil {
		t.Error("expected nil for nonexistent")
	}
}

func TestBundleByID(t *testing.T) {
	cat := DefaultCatalog()

	b := cat.BundleByID("go-dev")
	if b == nil {
		t.Fatal("expected to find go-dev bundle")
	}
	if b.Name != "Go Development" {
		t.Errorf("got name %q, want 'Go Development'", b.Name)
	}

	if cat.BundleByID("nonexistent") != nil {
		t.Error("expected nil for nonexistent bundle")
	}
}

func TestResolveBundle(t *testing.T) {
	cat := DefaultCatalog()

	entries := cat.ResolveBundle("rust")
	if len(entries) != 3 {
		t.Errorf("ResolveBundle(rust) returned %d entries, want 3", len(entries))
	}

	ids := make(map[string]bool)
	for _, e := range entries {
		ids[e.ID] = true
	}
	for _, want := range []string{"rust-best-practices", "rust-async-patterns", "rust-engineer"} {
		if !ids[want] {
			t.Errorf("ResolveBundle(rust) missing %q", want)
		}
	}

	if entries := cat.ResolveBundle("nonexistent"); entries != nil {
		t.Error("expected nil for nonexistent bundle")
	}
}

func TestAllCategories(t *testing.T) {
	cat := DefaultCatalog()
	cats := cat.AllCategories()
	if len(cats) == 0 {
		t.Fatal("expected categories")
	}
	// Check sorted.
	for i := 1; i < len(cats); i++ {
		if cats[i] < cats[i-1] {
			t.Errorf("categories not sorted: %q < %q", cats[i], cats[i-1])
		}
	}
}

func TestAllTags(t *testing.T) {
	cat := DefaultCatalog()
	tags := cat.AllTags()
	if len(tags) == 0 {
		t.Fatal("expected tags")
	}
	// Check sorted.
	for i := 1; i < len(tags); i++ {
		if tags[i] < tags[i-1] {
			t.Errorf("tags not sorted: %q < %q", tags[i], tags[i-1])
		}
	}
}

func TestBundlesForSkill(t *testing.T) {
	cat := DefaultCatalog()

	bundles := cat.BundlesForSkill("code-reviewer")
	if len(bundles) < 2 {
		t.Errorf("expected code-reviewer in at least 2 bundles, got %d", len(bundles))
	}

	bundleIDs := make(map[string]bool)
	for _, b := range bundles {
		bundleIDs[b.ID] = true
	}
	if !bundleIDs["security"] {
		t.Error("expected code-reviewer in security bundle")
	}
	if !bundleIDs["architecture"] {
		t.Error("expected code-reviewer in architecture bundle")
	}
}

func TestNilCatalog(t *testing.T) {
	var c *Catalog

	if c.ByCategory(CategoryLanguage) != nil {
		t.Error("nil catalog ByCategory should return nil")
	}
	if c.ByTag("go") != nil {
		t.Error("nil catalog ByTag should return nil")
	}
	if c.ByID("test") != nil {
		t.Error("nil catalog ByID should return nil")
	}
	if c.ByAlias("test") != nil {
		t.Error("nil catalog ByAlias should return nil")
	}
	if c.BundleByID("test") != nil {
		t.Error("nil catalog BundleByID should return nil")
	}
	if c.AllCategories() != nil {
		t.Error("nil catalog AllCategories should return nil")
	}
	if c.AllTags() != nil {
		t.Error("nil catalog AllTags should return nil")
	}
	if c.BundlesForSkill("test") != nil {
		t.Error("nil catalog BundlesForSkill should return nil")
	}
}
