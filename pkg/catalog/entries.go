package catalog

const (
	// Remote refs pinned to specific commits.
	pulumiAgentSkillsRef     = "github.com/pulumi/agent-skills@b6b942fc6e34517e2bbc52d6db04ca529baf3ad4"
	dirienClaudeSkillsRef    = "github.com/dirien/claude-skills@85b0ee2a07cb1e3420d445d3f2336eadca45cde5"
	jeffallanClaudeSkillsRef = "github.com/jeffallan/claude-skills@3bf9a24b76a7c122f1fc05e83929fbc84e1c207a"
	apolloSkillsRef          = "github.com/apollographql/skills@e1979d2f1e7c38cef58753b2bfd6fc9509101bdc"
	wshobsonAgentsRef        = "github.com/wshobson/agents@1ad2f007d5e9ec822a2d79e727ac1dcdf5f66f11"
	netresearchAgentRulesRef = "github.com/netresearch/agent-rules-skill@9b67bf594a52b1a7d38d8b0ec0a076a31f8d3d7e"
	rshadeAgentSkillsRef     = "github.com/rshade/agent-skills@4aff11fe89bb156337c2c7c303bb2db234cc9740"
)

// DefaultCatalog returns the complete catalog of all built-in and remote skills with bundles.
func DefaultCatalog() *Catalog {
	return &Catalog{
		Skills:  defaultSkills(),
		Bundles: defaultBundles(),
	}
}

func defaultSkills() []CatalogEntry {
	return []CatalogEntry{
		// Built-in skills.
		{
			ID: "commit", Name: "commit",
			Description: "Atomic, semantic-boundary git commit workflow",
			Category:    CategoryWorkflow, Tags: []string{"git", "commit", "conventional-commits"},
			Risk: RiskCritical, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "pr", Name: "pr",
			Description: "Create pull requests with structured description",
			Category:    CategoryWorkflow, Tags: []string{"git", "pull-request", "github"},
			Risk: RiskCritical, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "review", Name: "review",
			Description: "Review code changes for quality, security, and correctness",
			Category:    CategorySecurity, Tags: []string{"code-review", "security", "quality"},
			Risk: RiskSafe, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},

		// Built-in workflow commands.
		{
			ID: "yaah/init", Name: "yaah/init",
			Description: "Project onboarding: discover codebase, set vision, create roadmap",
			Category:    CategoryWorkflow, Tags: []string{"planning", "onboarding", "project-setup"},
			Risk: RiskCritical, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/discuss", Name: "yaah/discuss",
			Description: "Capture implementation decisions before planning",
			Category:    CategoryWorkflow, Tags: []string{"planning", "decisions", "context"},
			Risk: RiskSafe, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/plan", Name: "yaah/plan",
			Description: "Create a structured implementation plan for a project phase",
			Category:    CategoryWorkflow, Tags: []string{"planning", "tasks", "waves"},
			Risk: RiskCritical, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/execute", Name: "yaah/execute",
			Description: "Execute implementation plans wave by wave",
			Category:    CategoryWorkflow, Tags: []string{"planning", "execution", "waves"},
			Risk: RiskCritical, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/verify", Name: "yaah/verify",
			Description: "Verify implementation against plan requirements",
			Category:    CategoryWorkflow, Tags: []string{"planning", "verification", "quality"},
			Risk: RiskSafe, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/docs", Name: "yaah/docs",
			Description: "Generate or update project documentation from codebase analysis",
			Category:    CategoryWorkflow, Tags: []string{"documentation", "generation", "verification"},
			Risk: RiskCritical, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/next", Name: "yaah/next",
			Description: "Auto-detect and recommend the next workflow step",
			Category:    CategoryWorkflow, Tags: []string{"planning", "navigation", "status"},
			Risk: RiskNone, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/quick", Name: "yaah/quick",
			Description: "Execute a task without full planning overhead",
			Category:    CategoryWorkflow, Tags: []string{"planning", "quick", "task"},
			Risk: RiskCritical, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/ship", Name: "yaah/ship",
			Description: "Create a pull request from verified phase work",
			Category:    CategoryWorkflow, Tags: []string{"shipping", "pull-request", "release"},
			Risk: RiskCritical, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/pause", Name: "yaah/pause",
			Description: "Save current session state for later resumption",
			Category:    CategoryWorkflow, Tags: []string{"session", "handoff", "pause"},
			Risk: RiskSafe, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/resume", Name: "yaah/resume",
			Description: "Resume work from a previous session handoff",
			Category:    CategoryWorkflow, Tags: []string{"session", "handoff", "resume"},
			Risk: RiskSafe, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/complete-milestone", Name: "yaah/complete-milestone",
			Description: "Archive current milestone, tag release, generate changelog",
			Category:    CategoryWorkflow, Tags: []string{"milestone", "release", "archive"},
			Risk: RiskCritical, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/new-milestone", Name: "yaah/new-milestone",
			Description: "Start a new version cycle with fresh goals",
			Category:    CategoryWorkflow, Tags: []string{"milestone", "version", "planning"},
			Risk: RiskCritical, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/settings", Name: "yaah/settings",
			Description: "View or update workflow configuration",
			Category:    CategoryWorkflow, Tags: []string{"config", "settings", "workflow"},
			Risk: RiskSafe, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/add-phase", Name: "yaah/add-phase",
			Description: "Add a new phase to the end of the roadmap",
			Category:    CategoryWorkflow, Tags: []string{"planning", "phase", "roadmap"},
			Risk: RiskCritical, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/insert-phase", Name: "yaah/insert-phase",
			Description: "Insert an urgent phase between existing phases",
			Category:    CategoryWorkflow, Tags: []string{"planning", "phase", "roadmap"},
			Risk: RiskCritical, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/remove-phase", Name: "yaah/remove-phase",
			Description: "Remove a future phase from the roadmap",
			Category:    CategoryWorkflow, Tags: []string{"planning", "phase", "roadmap"},
			Risk: RiskCritical, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/health", Name: "yaah/health",
			Description: "Validate planning directory integrity and consistency",
			Category:    CategoryWorkflow, Tags: []string{"health", "validation", "integrity"},
			Risk: RiskNone, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/progress", Name: "yaah/progress",
			Description: "Show detailed project progress with metrics",
			Category:    CategoryWorkflow, Tags: []string{"status", "progress", "metrics"},
			Risk: RiskNone, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/review", Name: "yaah/review",
			Description: "Structured code review of phase implementation",
			Category:    CategoryWorkflow, Tags: []string{"code-review", "quality", "security"},
			Risk: RiskSafe, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/secure", Name: "yaah/secure",
			Description: "Security threat modeling and vulnerability analysis for a phase",
			Category:    CategorySecurity, Tags: []string{"security", "threat-model", "stride"},
			Risk: RiskSafe, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/todo", Name: "yaah/todo",
			Description: "Capture, list, or complete quick todo items",
			Category:    CategoryWorkflow, Tags: []string{"task", "todo", "capture"},
			Risk: RiskSafe, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/note", Name: "yaah/note",
			Description: "Capture an idea or observation for later",
			Category:    CategoryWorkflow, Tags: []string{"note", "idea", "capture"},
			Risk: RiskNone, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/cleanup", Name: "yaah/cleanup",
			Description: "Clean up temporary planning artifacts and state",
			Category:    CategoryWorkflow, Tags: []string{"cleanup", "maintenance"},
			Risk: RiskCritical, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/forensics", Name: "yaah/forensics",
			Description: "Investigate failed or stuck workflow runs",
			Category:    CategoryWorkflow, Tags: []string{"debug", "forensics", "recovery"},
			Risk: RiskSafe, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/explore", Name: "yaah/explore",
			Description: "Interactive codebase exploration and analysis",
			Category:    CategoryWorkflow, Tags: []string{"exploration", "analysis", "codebase"},
			Risk: RiskSafe, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/scan", Name: "yaah/scan",
			Description: "Scan codebase for patterns, issues, and improvement opportunities",
			Category:    CategorySecurity, Tags: []string{"scan", "security", "quality", "dependencies"},
			Risk: RiskSafe, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/import", Name: "yaah/import",
			Description: "Import an existing project into the planning workflow",
			Category:    CategoryWorkflow, Tags: []string{"import", "onboarding", "migration"},
			Risk: RiskCritical, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},
		{
			ID: "yaah/autonomous", Name: "yaah/autonomous",
			Description: "Run the full workflow autonomously for a phase",
			Category:    CategoryWorkflow, Tags: []string{"autonomous", "automation", "workflow"},
			Risk: RiskCritical, Tier: TierOfficial, Uses: "builtin", Repo: "yaah",
		},

		// pulumi/agent-skills — IaC authoring.
		{
			ID: "pulumi-best-practices", Name: "pulumi-best-practices",
			Description: "Pulumi best practices for reliable programs",
			Category:    CategoryInfrastructure, Tags: []string{"pulumi", "iac", "best-practices"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: pulumiAgentSkillsRef, Subpath: "authoring/skills/pulumi-best-practices/SKILL.md",
			Repo: "pulumi/agent-skills",
		},
		{
			ID: "pulumi-component", Name: "pulumi-component",
			Description: "Pulumi ComponentResource authoring",
			Category:    CategoryInfrastructure, Tags: []string{"pulumi", "iac", "component"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: pulumiAgentSkillsRef, Subpath: "authoring/skills/pulumi-component/SKILL.md",
			Repo: "pulumi/agent-skills",
		},
		{
			ID: "pulumi-automation-api", Name: "pulumi-automation-api",
			Description: "Pulumi Automation API best practices",
			Category:    CategoryInfrastructure, Tags: []string{"pulumi", "iac", "automation-api"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: pulumiAgentSkillsRef, Subpath: "authoring/skills/pulumi-automation-api/SKILL.md",
			Repo: "pulumi/agent-skills",
		},
		{
			ID: "pulumi-esc", Name: "pulumi-esc",
			Description: "Pulumi ESC environments, secrets, and configuration",
			Category:    CategoryInfrastructure, Tags: []string{"pulumi", "iac", "secrets", "esc"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: pulumiAgentSkillsRef, Subpath: "authoring/skills/pulumi-esc/SKILL.md",
			Repo: "pulumi/agent-skills",
		},

		// pulumi/agent-skills — migrations.
		{
			ID: "pulumi-terraform-to-pulumi", Name: "pulumi-terraform-to-pulumi",
			Description: "Convert Terraform to Pulumi",
			Category:    CategoryMigration, Tags: []string{"pulumi", "terraform", "migration"},
			Risk: RiskCritical, Tier: TierVerified,
			Uses: pulumiAgentSkillsRef, Subpath: "migration/skills/pulumi-terraform-to-pulumi/SKILL.md",
			Repo: "pulumi/agent-skills",
		},
		{
			ID: "pulumi-cdk-to-pulumi", Name: "pulumi-cdk-to-pulumi",
			Description: "Convert AWS CDK to Pulumi",
			Category:    CategoryMigration, Tags: []string{"pulumi", "cdk", "aws", "migration"},
			Risk: RiskCritical, Tier: TierVerified,
			Uses: pulumiAgentSkillsRef, Subpath: "migration/skills/pulumi-cdk-to-pulumi/SKILL.md",
			Repo: "pulumi/agent-skills",
		},
		{
			ID: "cloudformation-to-pulumi", Name: "cloudformation-to-pulumi",
			Description: "Convert CloudFormation to Pulumi",
			Category:    CategoryMigration, Tags: []string{"pulumi", "cloudformation", "aws", "migration"},
			Risk: RiskCritical, Tier: TierVerified,
			Uses: pulumiAgentSkillsRef, Subpath: "migration/skills/cloudformation-to-pulumi/SKILL.md",
			Repo: "pulumi/agent-skills",
		},
		{
			ID: "pulumi-arm-to-pulumi", Name: "pulumi-arm-to-pulumi",
			Description: "Convert Azure ARM/Bicep to Pulumi",
			Category:    CategoryMigration, Tags: []string{"pulumi", "arm", "bicep", "azure", "migration"},
			Risk: RiskCritical, Tier: TierVerified,
			Uses: pulumiAgentSkillsRef, Subpath: "migration/skills/pulumi-arm-to-pulumi/SKILL.md",
			Repo: "pulumi/agent-skills",
		},

		// dirien/claude-skills — Pulumi language-specific.
		{
			ID: "pulumi-typescript", Name: "pulumi-typescript",
			Description: "Pulumi TypeScript IaC with ESC and OIDC",
			Category:    CategoryInfrastructure, Tags: []string{"pulumi", "typescript", "iac", "oidc"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: dirienClaudeSkillsRef, Subpath: "pulumi-typescript/SKILL.md",
			Repo: "dirien/claude-skills",
		},
		{
			ID: "pulumi-go", Name: "pulumi-go",
			Description: "Pulumi Go IaC with ESC and OIDC",
			Category:    CategoryInfrastructure, Tags: []string{"pulumi", "go", "iac", "oidc"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: dirienClaudeSkillsRef, Subpath: "pulumi-go/SKILL.md",
			Repo: "dirien/claude-skills",
		},
		{
			ID: "pulumi-python", Name: "pulumi-python",
			Description: "Pulumi Python IaC with ESC and OIDC",
			Category:    CategoryInfrastructure, Tags: []string{"pulumi", "python", "iac", "oidc"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: dirienClaudeSkillsRef, Subpath: "pulumi-python/SKILL.md",
			Repo: "dirien/claude-skills",
		},
		{
			ID: "pulumi-neo", Name: "pulumi-neo",
			Description: "Pulumi Neo conversational infrastructure management",
			Category:    CategoryInfrastructure, Tags: []string{"pulumi", "neo", "ai", "conversational"},
			Risk: RiskCritical, Tier: TierVerified,
			Uses: dirienClaudeSkillsRef, Subpath: "pulumi-neo/SKILL.md",
			Repo: "dirien/claude-skills",
		},
		{
			ID: "pulumi-cli", Name: "pulumi-cli",
			Description: "Pulumi CLI command reference for deployments",
			Category:    CategoryInfrastructure, Tags: []string{"pulumi", "cli", "deployment"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: dirienClaudeSkillsRef, Subpath: "pulumi-cli/SKILL.md",
			Repo: "dirien/claude-skills",
		},

		// dirien/claude-skills — Flux CD GitOps.
		{
			ID: "flux-cli", Name: "flux-cli",
			Description: "GitOps for Kubernetes using Flux CD CLI",
			Category:    CategoryInfrastructure, Tags: []string{"flux", "gitops", "kubernetes", "cli"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: dirienClaudeSkillsRef, Subpath: "flux-cli/SKILL.md",
			Repo: "dirien/claude-skills",
		},
		{
			ID: "flux-operator-cli", Name: "flux-operator-cli",
			Description: "Flux Operator CLI for managing Flux CD deployments on Kubernetes",
			Category:    CategoryInfrastructure, Tags: []string{"flux", "gitops", "kubernetes", "operator", "cli"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: dirienClaudeSkillsRef, Subpath: "flux-operator-cli/SKILL.md",
			Repo: "dirien/claude-skills",
		},

		// jeffallan/claude-skills — Development and operations.
		{
			ID: "golang-pro", Name: "golang-pro",
			Description: "Go concurrent patterns, microservices, gRPC, and performance optimization",
			Category:    CategoryLanguage, Tags: []string{"go", "grpc", "microservices", "concurrency", "performance"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: jeffallanClaudeSkillsRef, Subpath: "skills/golang-pro/SKILL.md",
			Aliases: []string{"go-pro", "go-developer"},
			Repo:    "jeffallan/claude-skills",
		},
		{
			ID: "kubernetes-specialist", Name: "kubernetes-specialist",
			Description: "Kubernetes deployments, Helm, RBAC, NetworkPolicies, and multi-cluster",
			Category:    CategoryDevOps, Tags: []string{"kubernetes", "k8s", "helm", "rbac", "containers"},
			Risk: RiskCritical, Tier: TierVerified,
			Uses: jeffallanClaudeSkillsRef, Subpath: "skills/kubernetes-specialist/SKILL.md",
			Aliases: []string{"k8s-specialist", "k8s"},
			Repo:    "jeffallan/claude-skills",
		},
		{
			ID: "devops-engineer", Name: "devops-engineer",
			Description: "CI/CD pipelines, Docker, Kubernetes, Terraform, and GitOps",
			Category:    CategoryDevOps, Tags: []string{"cicd", "docker", "kubernetes", "terraform", "gitops"},
			Risk: RiskCritical, Tier: TierVerified,
			Uses: jeffallanClaudeSkillsRef, Subpath: "skills/devops-engineer/SKILL.md",
			Repo: "jeffallan/claude-skills",
		},
		{
			ID: "python-pro", Name: "python-pro",
			Description: "Python 3.11+ with type safety, async, pytest, and ruff",
			Category:    CategoryLanguage, Tags: []string{"python", "async", "pytest", "typing", "ruff"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: jeffallanClaudeSkillsRef, Subpath: "skills/python-pro/SKILL.md",
			Aliases: []string{"python-developer"},
			Repo:    "jeffallan/claude-skills",
		},
		{
			ID: "typescript-pro", Name: "typescript-pro",
			Description: "Advanced TypeScript types, generics, tRPC, and monorepo setup",
			Category:    CategoryLanguage, Tags: []string{"typescript", "generics", "trpc", "monorepo"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: jeffallanClaudeSkillsRef, Subpath: "skills/typescript-pro/SKILL.md",
			Aliases: []string{"ts-pro"},
			Repo:    "jeffallan/claude-skills",
		},
		{
			ID: "csharp-developer", Name: "csharp-developer",
			Description: "C# .NET 8+, ASP.NET Core, Blazor, EF Core, and MediatR",
			Category:    CategoryLanguage, Tags: []string{"csharp", "dotnet", "aspnet", "blazor", "efcore"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: jeffallanClaudeSkillsRef, Subpath: "skills/csharp-developer/SKILL.md",
			Aliases: []string{"dotnet-developer", "cs-developer"},
			Repo:    "jeffallan/claude-skills",
		},
		{
			ID: "javascript-pro", Name: "javascript-pro",
			Description: "Modern ES2023+ JavaScript, async/await, ESM, and Node.js",
			Category:    CategoryLanguage, Tags: []string{"javascript", "nodejs", "esm", "async"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: jeffallanClaudeSkillsRef, Subpath: "skills/javascript-pro/SKILL.md",
			Aliases: []string{"js-pro"},
			Repo:    "jeffallan/claude-skills",
		},
		{
			ID: "cli-developer", Name: "cli-developer",
			Description: "CLI tools with argument parsing, completions, and cross-platform support",
			Category:    CategoryLanguage, Tags: []string{"cli", "argument-parsing", "shell-completions"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: jeffallanClaudeSkillsRef, Subpath: "skills/cli-developer/SKILL.md",
			Repo: "jeffallan/claude-skills",
		},
		{
			ID: "sre-engineer", Name: "sre-engineer",
			Description: "SLOs, error budgets, incident response, and capacity planning",
			Category:    CategoryDevOps, Tags: []string{"sre", "slo", "observability", "incident-response"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: jeffallanClaudeSkillsRef, Subpath: "skills/sre-engineer/SKILL.md",
			Aliases: []string{"sre"},
			Repo:    "jeffallan/claude-skills",
		},
		{
			ID: "the-fool", Name: "the-fool",
			Description: "Devil's advocate, pre-mortems, red teaming, and assumption auditing",
			Category:    CategoryArchitecture, Tags: []string{"critical-thinking", "red-team", "pre-mortem"},
			Risk: RiskNone, Tier: TierVerified,
			Uses: jeffallanClaudeSkillsRef, Subpath: "skills/the-fool/SKILL.md",
			Repo: "jeffallan/claude-skills",
		},
		{
			ID: "architecture-designer", Name: "architecture-designer",
			Description: "System architecture, ADRs, trade-offs, and scalability planning",
			Category:    CategoryArchitecture, Tags: []string{"architecture", "adr", "system-design", "scalability"},
			Risk: RiskNone, Tier: TierVerified,
			Uses: jeffallanClaudeSkillsRef, Subpath: "skills/architecture-designer/SKILL.md",
			Aliases: []string{"architect"},
			Repo:    "jeffallan/claude-skills",
		},
		{
			ID: "spring-boot-engineer", Name: "spring-boot-engineer",
			Description: "Spring Boot 3.x, Spring Security 6, JPA, WebFlux, and Spring Cloud",
			Category:    CategoryLanguage, Tags: []string{"java", "spring-boot", "spring-security", "jpa", "webflux"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: jeffallanClaudeSkillsRef, Subpath: "skills/spring-boot-engineer/SKILL.md",
			Repo: "jeffallan/claude-skills",
		},
		{
			ID: "code-reviewer", Name: "code-reviewer",
			Description: "Code review for bugs, security, performance, and maintainability",
			Category:    CategorySecurity, Tags: []string{"code-review", "security", "performance", "bugs"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: jeffallanClaudeSkillsRef, Subpath: "skills/code-reviewer/SKILL.md",
			Repo: "jeffallan/claude-skills",
		},

		// Rust skills.
		{
			ID: "rust-best-practices", Name: "rust-best-practices",
			Description: "Idiomatic Rust code, borrowing, error handling, and performance optimization",
			Category:    CategoryLanguage, Tags: []string{"rust", "borrowing", "error-handling", "performance"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: apolloSkillsRef, Subpath: "skills/rust-best-practices/SKILL.md",
			Repo: "apollographql/skills",
		},
		{
			ID: "rust-async-patterns", Name: "rust-async-patterns",
			Description: "Rust async programming with Tokio, async traits, and concurrent patterns",
			Category:    CategoryLanguage, Tags: []string{"rust", "async", "tokio", "concurrency"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: wshobsonAgentsRef, Subpath: "plugins/systems-programming/skills/rust-async-patterns/SKILL.md",
			Repo: "wshobson/agents",
		},
		{
			ID: "rust-engineer", Name: "rust-engineer",
			Description: "Idiomatic Rust with ownership, lifetimes, traits, tokio, and error handling",
			Category:    CategoryLanguage, Tags: []string{"rust", "ownership", "lifetimes", "traits", "tokio"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: jeffallanClaudeSkillsRef, Subpath: "skills/rust-engineer/SKILL.md",
			Repo: "jeffallan/claude-skills",
		},

		// netresearch/agent-rules-skill.
		{
			ID: "agent-rules", Name: "agent-rules",
			Description: "Generate and maintain AGENTS.md files following the agents.md convention",
			Category:    CategoryWorkflow, Tags: []string{"agents-md", "onboarding", "documentation"},
			Risk: RiskCritical, Tier: TierVerified,
			Uses: netresearchAgentRulesRef, Subpath: "skills/agent-rules/SKILL.md",
			Repo: "netresearch/agent-rules-skill",
		},

		// rshade/agent-skills — Code quality, security, and workflow automation.
		{
			ID: "agent-ready-go", Name: "agent-ready-go",
			Description: "Prepare Go apps for AI agent interaction with structured logging and CLI design",
			Category:    CategoryLanguage, Tags: []string{"go", "agent-ready", "logging", "cli"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: rshadeAgentSkillsRef, Subpath: "skills/agent-ready-go/SKILL.md",
			Repo: "rshade/agent-skills",
		},
		{
			ID: "commitlint", Name: "commitlint",
			Description: "Validate commit messages against Conventional Commits specification",
			Category:    CategoryWorkflow, Tags: []string{"git", "commits", "linting", "conventional-commits"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: rshadeAgentSkillsRef, Subpath: "skills/commitlint/SKILL.md",
			Repo: "rshade/agent-skills",
		},
		{
			ID: "decide", Name: "decide",
			Description: "Three-agent adversarial debate protocol for strategic decisions",
			Category:    CategoryArchitecture, Tags: []string{"decision-making", "debate", "tradeoffs"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: rshadeAgentSkillsRef, Subpath: "skills/decide/SKILL.md",
			Repo: "rshade/agent-skills",
		},
		{
			ID: "dep-upgrade", Name: "dep-upgrade",
			Description: "Safe systematic dependency upgrades with vulnerability scanning and rollback",
			Category:    CategoryWorkflow, Tags: []string{"dependencies", "upgrades", "security", "vulnerability"},
			Risk: RiskCritical, Tier: TierVerified,
			Uses: rshadeAgentSkillsRef, Subpath: "skills/dep-upgrade/SKILL.md",
			Repo: "rshade/agent-skills",
		},
		{
			ID: "design-principles", Name: "design-principles",
			Description: "Analyze codebases against SOLID, DRY, YAGNI, KISS, and other design principles",
			Category:    CategoryArchitecture, Tags: []string{"design", "solid", "code-quality", "audit"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: rshadeAgentSkillsRef, Subpath: "skills/design-principles/SKILL.md",
			Repo: "rshade/agent-skills",
		},
		{
			ID: "go-nolint-audit", Name: "go-nolint-audit",
			Description: "Audit nolint directives in Go codebases for stale or unjustified suppressions",
			Category:    CategoryLanguage, Tags: []string{"go", "linting", "nolint", "audit"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: rshadeAgentSkillsRef, Subpath: "skills/go-nolint-audit/SKILL.md",
			Repo: "rshade/agent-skills",
		},
		{
			ID: "lint-fix", Name: "lint-fix",
			Description: "Detect linting tools, run them to zero errors, and fix issues atomically",
			Category:    CategoryWorkflow, Tags: []string{"linting", "fixing", "code-quality"},
			Risk: RiskCritical, Tier: TierVerified,
			Uses: rshadeAgentSkillsRef, Subpath: "skills/lint-fix/SKILL.md",
			Repo: "rshade/agent-skills",
		},
		{
			ID: "markdownlint", Name: "markdownlint",
			Description: "Validate markdown files against formatting standards with auto-fix",
			Category:    CategoryWorkflow, Tags: []string{"markdown", "linting", "documentation"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: rshadeAgentSkillsRef, Subpath: "skills/markdownlint/SKILL.md",
			Repo: "rshade/agent-skills",
		},
		{
			ID: "pull-request-msg", Name: "pull-request-msg",
			Description: "Generate structured PR descriptions from session context using GitHub CLI",
			Category:    CategoryWorkflow, Tags: []string{"github", "pull-request", "automation"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: rshadeAgentSkillsRef, Subpath: "skills/pull-request-msg-with-gh/SKILL.md",
			Repo: "rshade/agent-skills",
		},
		{
			ID: "roadmap", Name: "roadmap",
			Description: "Strategic roadmap management synced with GitHub Issues and labels",
			Category:    CategoryWorkflow, Tags: []string{"roadmap", "planning", "github", "issues"},
			Risk: RiskCritical, Tier: TierVerified,
			Uses: rshadeAgentSkillsRef, Subpath: "skills/roadmap/SKILL.md",
			Repo: "rshade/agent-skills",
		},
		{
			ID: "scout", Name: "scout",
			Description: "Identify top improvement opportunities in files you are touching",
			Category:    CategoryWorkflow, Tags: []string{"code-quality", "improvement", "review"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: rshadeAgentSkillsRef, Subpath: "skills/scout/SKILL.md",
			Repo: "rshade/agent-skills",
		},
		{
			ID: "security-audit", Name: "security-audit",
			Description: "Comprehensive vulnerability assessment with OWASP Top 10 and threat modeling",
			Category:    CategorySecurity, Tags: []string{"security", "owasp", "vulnerability", "audit"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: rshadeAgentSkillsRef, Subpath: "skills/security-audit/SKILL.md",
			Repo: "rshade/agent-skills",
		},
		{
			ID: "tailscale-install", Name: "tailscale-install",
			Description: "Install and configure Tailscale across platforms including WSL2 and containers",
			Category:    CategoryInfrastructure, Tags: []string{"tailscale", "vpn", "networking", "install"},
			Risk: RiskCritical, Tier: TierVerified,
			Uses: rshadeAgentSkillsRef, Subpath: "skills/tailscale-install/SKILL.md",
			Repo: "rshade/agent-skills",
		},
		{
			ID: "tech-debt", Name: "tech-debt",
			Description: "Systematic technical debt analysis across 9 categories with health scoring",
			Category:    CategoryArchitecture, Tags: []string{"tech-debt", "audit", "code-quality", "architecture"},
			Risk: RiskSafe, Tier: TierVerified,
			Uses: rshadeAgentSkillsRef, Subpath: "skills/tech-debt/SKILL.md",
			Repo: "rshade/agent-skills",
		},
	}
}

func defaultBundles() []Bundle {
	return []Bundle{
		{
			ID: "go-dev", Name: "Go Development",
			Description: "Full Go development stack",
			SkillIDs:    []string{"golang-pro", "pulumi-go", "cli-developer", "agent-ready-go", "go-nolint-audit"},
		},
		{
			ID: "pulumi-core", Name: "Pulumi Core",
			Description: "Essential Pulumi IaC skills",
			SkillIDs:    []string{"pulumi-best-practices", "pulumi-component", "pulumi-automation-api", "pulumi-esc", "pulumi-cli"},
		},
		{
			ID: "pulumi-migration", Name: "Pulumi Migration",
			Description: "Cloud migration toolkit for moving to Pulumi",
			SkillIDs:    []string{"pulumi-terraform-to-pulumi", "pulumi-cdk-to-pulumi", "cloudformation-to-pulumi", "pulumi-arm-to-pulumi"},
		},
		{
			ID: "pulumi-languages", Name: "Pulumi Languages",
			Description: "Pulumi language-specific IaC skills",
			SkillIDs:    []string{"pulumi-typescript", "pulumi-go", "pulumi-python", "pulumi-neo", "pulumi-cli"},
		},
		{
			ID: "security", Name: "Security",
			Description: "Security-focused review and analysis skills",
			SkillIDs:    []string{"review", "code-reviewer", "the-fool", "security-audit"},
		},
		{
			ID: "full-stack", Name: "Full Stack",
			Description: "Language skills for full-stack web development",
			SkillIDs:    []string{"typescript-pro", "javascript-pro", "python-pro", "csharp-developer", "spring-boot-engineer"},
		},
		{
			ID: "devops", Name: "DevOps & SRE",
			Description: "Infrastructure, operations, and reliability skills",
			SkillIDs:    []string{"devops-engineer", "sre-engineer", "kubernetes-specialist", "flux-cli", "flux-operator-cli"},
		},
		{
			ID: "rust", Name: "Rust Development",
			Description: "Rust ecosystem skills",
			SkillIDs:    []string{"rust-best-practices", "rust-async-patterns", "rust-engineer"},
		},
		{
			ID: "code-quality", Name: "Code Quality",
			Description: "Linting, auditing, and code improvement skills",
			SkillIDs:    []string{"lint-fix", "commitlint", "markdownlint", "scout", "design-principles", "tech-debt", "dep-upgrade"},
		},
		{
			ID: "architecture", Name: "Architecture",
			Description: "System design, review, and critical thinking",
			SkillIDs:    []string{"architecture-designer", "code-reviewer", "the-fool"},
		},
		{
			ID: "workflow", Name: "Project Workflow",
			Description: "Structured project lifecycle with 29 workflow commands",
			SkillIDs: []string{
				"yaah/init", "yaah/discuss", "yaah/plan", "yaah/execute", "yaah/verify", "yaah/docs", "yaah/next", "yaah/quick",
				"yaah/ship", "yaah/pause", "yaah/resume", "yaah/complete-milestone", "yaah/new-milestone", "yaah/settings",
				"yaah/add-phase", "yaah/insert-phase", "yaah/remove-phase", "yaah/health", "yaah/progress", "yaah/review",
				"yaah/secure", "yaah/todo", "yaah/note", "yaah/cleanup", "yaah/forensics", "yaah/explore", "yaah/scan", "yaah/import", "yaah/autonomous",
			},
		},
	}
}
