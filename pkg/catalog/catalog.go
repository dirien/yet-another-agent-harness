// Package catalog provides a discoverable index of all available skills and bundles.
// It separates the "what is available" concern from the "what is registered" concern
// handled by the skills.Registry.
package catalog

// Category classifies a skill by its primary domain.
type Category string

const (
	CategoryLanguage       Category = "language"
	CategoryInfrastructure Category = "infrastructure"
	CategorySecurity       Category = "security"
	CategoryDevOps         Category = "devops"
	CategoryWorkflow       Category = "workflow"
	CategoryMigration      Category = "migration"
	CategoryArchitecture   Category = "architecture"
	CategoryTesting        Category = "testing"
)

// RiskLevel classifies a skill's potential impact on the workspace.
type RiskLevel string

const (
	RiskNone     RiskLevel = "none"     // Text/reasoning only.
	RiskSafe     RiskLevel = "safe"     // File reads, no mutations.
	RiskCritical RiskLevel = "critical" // Filesystem/database changes.
)

// QualityTier indicates the provenance and support level of a skill.
type QualityTier string

const (
	TierOfficial  QualityTier = "official"  // Maintained by yaah core.
	TierVerified  QualityTier = "verified"  // Third-party, manually reviewed.
	TierCommunity QualityTier = "community" // Third-party, unreviewed.
)

// CatalogEntry describes a single skill available for installation.
type CatalogEntry struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Category    Category    `json:"category"`
	Tags        []string    `json:"tags"`
	Risk        RiskLevel   `json:"risk"`
	Tier        QualityTier `json:"tier"`
	Uses        string      `json:"uses"`              // Git ref: "github.com/owner/repo@ref", or "builtin".
	Subpath     string      `json:"subpath,omitempty"` // Path within repo to SKILL.md.
	Aliases     []string    `json:"aliases,omitempty"` // Alternative names.
	Repo        string      `json:"repo"`              // Display-friendly repo name.
}

// Bundle is a named group of skill IDs for a role or use case.
type Bundle struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	SkillIDs    []string `json:"skillIds"`
}

// Catalog holds the full index of available skills and bundles.
type Catalog struct {
	Skills  []CatalogEntry `json:"skills"`
	Bundles []Bundle       `json:"bundles"`
}
