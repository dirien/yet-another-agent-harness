package catalog

import (
	"slices"
	"strings"
)

// Search returns entries whose name, description, tags, or aliases match the query.
// Matching is case-insensitive substring.
func (c *Catalog) Search(query string) []CatalogEntry {
	if c == nil || query == "" {
		return nil
	}
	q := strings.ToLower(query)
	var results []CatalogEntry
	for _, e := range c.Skills {
		if matchesEntry(e, q) {
			results = append(results, e)
		}
	}
	return results
}

// ByCategory returns entries matching the given category.
func (c *Catalog) ByCategory(cat Category) []CatalogEntry {
	if c == nil {
		return nil
	}
	var results []CatalogEntry
	for _, e := range c.Skills {
		if e.Category == cat {
			results = append(results, e)
		}
	}
	return results
}

// ByTag returns entries that have the given tag (case-insensitive).
func (c *Catalog) ByTag(tag string) []CatalogEntry {
	if c == nil {
		return nil
	}
	t := strings.ToLower(tag)
	var results []CatalogEntry
	for _, e := range c.Skills {
		for _, et := range e.Tags {
			if strings.ToLower(et) == t {
				results = append(results, e)
				break
			}
		}
	}
	return results
}

// ByID returns a single entry by its canonical ID, or nil.
func (c *Catalog) ByID(id string) *CatalogEntry {
	if c == nil {
		return nil
	}
	for i := range c.Skills {
		if c.Skills[i].ID == id {
			return &c.Skills[i]
		}
	}
	return nil
}

// ByAlias returns a single entry by any of its aliases, or nil.
func (c *Catalog) ByAlias(alias string) *CatalogEntry {
	if c == nil {
		return nil
	}
	a := strings.ToLower(alias)
	for i := range c.Skills {
		for _, al := range c.Skills[i].Aliases {
			if strings.ToLower(al) == a {
				return &c.Skills[i]
			}
		}
	}
	return nil
}

// Resolve looks up an entry by ID first, then by alias.
func (c *Catalog) Resolve(nameOrAlias string) *CatalogEntry {
	if e := c.ByID(nameOrAlias); e != nil {
		return e
	}
	return c.ByAlias(nameOrAlias)
}

// BundleByID returns a bundle by its ID, or nil.
func (c *Catalog) BundleByID(id string) *Bundle {
	if c == nil {
		return nil
	}
	for i := range c.Bundles {
		if c.Bundles[i].ID == id {
			return &c.Bundles[i]
		}
	}
	return nil
}

// ResolveBundle returns all catalog entries for a bundle's skill IDs.
func (c *Catalog) ResolveBundle(bundleID string) []CatalogEntry {
	b := c.BundleByID(bundleID)
	if b == nil {
		return nil
	}
	idSet := make(map[string]struct{}, len(b.SkillIDs))
	for _, id := range b.SkillIDs {
		idSet[id] = struct{}{}
	}
	var results []CatalogEntry
	for _, e := range c.Skills {
		if _, ok := idSet[e.ID]; ok {
			results = append(results, e)
		}
	}
	return results
}

// AllCategories returns the distinct set of categories in the catalog, sorted.
func (c *Catalog) AllCategories() []Category {
	if c == nil {
		return nil
	}
	seen := make(map[Category]struct{})
	for _, e := range c.Skills {
		seen[e.Category] = struct{}{}
	}
	cats := make([]Category, 0, len(seen))
	for cat := range seen {
		cats = append(cats, cat)
	}
	slices.Sort(cats)
	return cats
}

// AllTags returns the distinct set of tags in the catalog, sorted.
func (c *Catalog) AllTags() []string {
	if c == nil {
		return nil
	}
	seen := make(map[string]struct{})
	for _, e := range c.Skills {
		for _, t := range e.Tags {
			seen[t] = struct{}{}
		}
	}
	tags := make([]string, 0, len(seen))
	for t := range seen {
		tags = append(tags, t)
	}
	slices.Sort(tags)
	return tags
}

// BundlesForSkill returns all bundles that include the given skill ID.
func (c *Catalog) BundlesForSkill(skillID string) []Bundle {
	if c == nil {
		return nil
	}
	var results []Bundle
	for _, b := range c.Bundles {
		if slices.Contains(b.SkillIDs, skillID) {
			results = append(results, b)
		}
	}
	return results
}

func matchesEntry(e CatalogEntry, q string) bool {
	if strings.Contains(strings.ToLower(e.ID), q) {
		return true
	}
	if strings.Contains(strings.ToLower(e.Name), q) {
		return true
	}
	if strings.Contains(strings.ToLower(e.Description), q) {
		return true
	}
	for _, t := range e.Tags {
		if strings.Contains(strings.ToLower(t), q) {
			return true
		}
	}
	for _, a := range e.Aliases {
		if strings.Contains(strings.ToLower(a), q) {
			return true
		}
	}
	return false
}
