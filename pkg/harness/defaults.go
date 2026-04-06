package harness

import (
	"context"
	"regexp"
	"slices"

	agentpkg "github.com/dirien/yet-another-agent-harness/pkg/agents"
	"github.com/dirien/yet-another-agent-harness/pkg/catalog"
	cmdbuiltins "github.com/dirien/yet-another-agent-harness/pkg/commands/builtins"
	"github.com/dirien/yet-another-agent-harness/pkg/hooks"
	"github.com/dirien/yet-another-agent-harness/pkg/hooks/handlers"
	lspproviders "github.com/dirien/yet-another-agent-harness/pkg/lsp/providers"
	mcpproviders "github.com/dirien/yet-another-agent-harness/pkg/mcp/providers"
	"github.com/dirien/yet-another-agent-harness/pkg/plugins"
	"github.com/dirien/yet-another-agent-harness/pkg/schema"
	"github.com/dirien/yet-another-agent-harness/pkg/skills"
	"github.com/dirien/yet-another-agent-harness/pkg/skills/builtins"
)

// DefaultOptions controls which built-in components to enlist.
type DefaultOptions struct {
	// Handlers
	EnableCommandGuard      bool
	EnableCommentChecker    bool
	EnableSecretScanner     bool
	EnableSecretRemediation bool
	EnableSessionLogger     bool
	LintProfiles            []handlers.Profile

	// Providers (MCP)
	EnableContext7  bool
	EnablePulumiMCP bool
	EnableYaahMCP   bool
	NotionToken     string

	// Skills
	EnableCommitSkill bool
	EnablePRSkill     bool
	EnableReviewSkill bool

	// Remote skills — pulumi/agent-skills
	EnablePulumiBestPractices    bool
	EnablePulumiComponent        bool
	EnablePulumiAutomationAPI    bool
	EnablePulumiESC              bool
	EnablePulumiTerraformMigrate bool
	EnablePulumiCDKMigrate       bool
	EnablePulumiCFNMigrate       bool
	EnablePulumiARMMigrate       bool

	// Remote skills — dirien/claude-skills
	EnablePulumiTypeScript bool
	EnablePulumiGo         bool
	EnablePulumiPython     bool
	EnablePulumiNeo        bool
	EnablePulumiCLI        bool
	EnableFluxCLI          bool
	EnableFluxOperatorCLI  bool

	// Remote skills — jeffallan/claude-skills
	EnableGolangPro            bool
	EnableKubernetesSpecialist bool
	EnableDevOpsEngineer       bool
	EnablePythonPro            bool
	EnableTypeScriptPro        bool
	EnableCSharpDeveloper      bool
	EnableJavaScriptPro        bool
	EnableCLIDeveloper         bool
	EnableSREEngineer          bool
	EnableTheFool              bool
	EnableArchitectureDesigner bool
	EnableSpringBootEngineer   bool
	EnableCodeReviewer         bool

	// Remote skills — Rust
	EnableRustBestPractices bool
	EnableRustAsyncPatterns bool
	EnableRustEngineer      bool

	// Remote skills — netresearch/agent-rules-skill
	EnableAgentRules bool

	// Remote skills — rshade/agent-skills
	EnableAgentReadyGo     bool
	EnableCommitlint       bool
	EnableDecide           bool
	EnableDepUpgrade       bool
	EnableDesignPrinciples bool
	EnableGoNolintAudit    bool
	EnableLintFix          bool
	EnableMarkdownlint     bool
	EnablePullRequestMsg   bool
	EnableRoadmap          bool
	EnableScout            bool
	EnableSecurityAudit    bool
	EnableTailscaleInstall bool
	EnableTechDebt         bool

	// Catalog-based skill selection (overrides individual Enable* flags for skills when set).
	SkillIDs   []string // Register only these skills from the catalog.
	BundleIDs  []string // Resolve bundles and register their skills.
	ExcludeIDs []string // Exclude specific skills.

	// Plugins (marketplace-backed)
	EnableCodexPlugin bool

	// LSP servers (marketplace-backed)
	EnableGopls      bool
	EnablePyright    bool
	EnableTypeScript bool
	EnableCSharp     bool

	// Workflow commands
	EnableInitCommand              bool
	EnableDiscussCommand           bool
	EnablePlanCommand              bool
	EnableExecuteCommand           bool
	EnableVerifyCommand            bool
	EnableDocsCommand              bool
	EnableNextCommand              bool
	EnableQuickCommand             bool
	EnableShipCommand              bool
	EnablePauseCommand             bool
	EnableResumeCommand            bool
	EnableCompleteMilestoneCommand bool
	EnableNewMilestoneCommand      bool
	EnableSettingsCommand          bool
	EnableAddPhaseCommand          bool
	EnableInsertPhaseCommand       bool
	EnableRemovePhaseCommand       bool
	EnableHealthCommand            bool
	EnableProgressCommand          bool
	EnableCodeReviewCommand        bool
	EnableSecureCommand            bool
	EnableTodoCommand              bool
	EnableNoteCommand              bool
	EnableCleanupCommand           bool
	EnableForensicsCommand         bool
	EnableExploreCommand           bool
	EnableScanCommand              bool
	EnableImportCommand            bool
	EnableAutonomousCommand        bool

	// Agents
	EnableExecutor  bool
	EnableLibrarian bool
	EnableReviewer  bool

	// Workflow agents
	EnableResearcher bool
	EnablePlanner    bool
	EnableDocWriter  bool
	EnableVerifier   bool

	// Remote agents — msitarzewski/agency-agents
	EnableAgencyAIEngineer             bool
	EnableAgencyBackendArchitect       bool
	EnableAgencySecurityEngineer       bool
	EnableAgencyCodeReviewerAgent      bool
	EnableAgencySoftwareArchitect      bool
	EnableAgencyDevOpsAutomator        bool
	EnableAgencySRE                    bool
	EnableAgencyAPITester              bool
	EnableAgencyPerformanceBenchmarker bool

	// Settings
	Settings *schema.Settings
}

// AllDefaults returns options with everything enabled (except Notion which needs a token).
func AllDefaults() DefaultOptions {
	thinking := true
	return DefaultOptions{
		EnableCommandGuard:      true,
		EnableCommentChecker:    true,
		EnableSecretScanner:     true,
		EnableSecretRemediation: true,
		EnableSessionLogger:     true,
		LintProfiles: []handlers.Profile{
			handlers.GolangCILint(),
			handlers.Ruff(),
			handlers.Prettier(),
			handlers.TypeScript(),
			handlers.RustFmt(),
		},
		EnableContext7:                     true,
		EnablePulumiMCP:                    true,
		EnableYaahMCP:                      true,
		EnableCommitSkill:                  true,
		EnablePRSkill:                      true,
		EnableReviewSkill:                  true,
		EnablePulumiBestPractices:          true,
		EnablePulumiComponent:              true,
		EnablePulumiAutomationAPI:          true,
		EnablePulumiESC:                    true,
		EnablePulumiTerraformMigrate:       true,
		EnablePulumiCDKMigrate:             true,
		EnablePulumiCFNMigrate:             true,
		EnablePulumiARMMigrate:             true,
		EnablePulumiTypeScript:             true,
		EnablePulumiGo:                     true,
		EnablePulumiPython:                 true,
		EnablePulumiNeo:                    true,
		EnablePulumiCLI:                    true,
		EnableFluxCLI:                      true,
		EnableFluxOperatorCLI:              true,
		EnableGolangPro:                    true,
		EnableKubernetesSpecialist:         true,
		EnableDevOpsEngineer:               true,
		EnablePythonPro:                    true,
		EnableTypeScriptPro:                true,
		EnableCSharpDeveloper:              true,
		EnableJavaScriptPro:                true,
		EnableCLIDeveloper:                 true,
		EnableSREEngineer:                  true,
		EnableTheFool:                      true,
		EnableArchitectureDesigner:         true,
		EnableSpringBootEngineer:           true,
		EnableCodeReviewer:                 true,
		EnableRustBestPractices:            true,
		EnableRustAsyncPatterns:            true,
		EnableRustEngineer:                 true,
		EnableAgentRules:                   true,
		EnableAgentReadyGo:                 true,
		EnableCommitlint:                   true,
		EnableDecide:                       true,
		EnableDepUpgrade:                   true,
		EnableDesignPrinciples:             true,
		EnableGoNolintAudit:                true,
		EnableLintFix:                      true,
		EnableMarkdownlint:                 true,
		EnablePullRequestMsg:               true,
		EnableRoadmap:                      true,
		EnableScout:                        true,
		EnableSecurityAudit:                true,
		EnableTailscaleInstall:             true,
		EnableTechDebt:                     true,
		EnableInitCommand:                  true,
		EnableDiscussCommand:               true,
		EnablePlanCommand:                  true,
		EnableExecuteCommand:               true,
		EnableVerifyCommand:                true,
		EnableDocsCommand:                  true,
		EnableNextCommand:                  true,
		EnableQuickCommand:                 true,
		EnableShipCommand:                  true,
		EnablePauseCommand:                 true,
		EnableResumeCommand:                true,
		EnableCompleteMilestoneCommand:     true,
		EnableNewMilestoneCommand:          true,
		EnableSettingsCommand:              true,
		EnableAddPhaseCommand:              true,
		EnableInsertPhaseCommand:           true,
		EnableRemovePhaseCommand:           true,
		EnableHealthCommand:                true,
		EnableProgressCommand:              true,
		EnableCodeReviewCommand:            true,
		EnableSecureCommand:                true,
		EnableTodoCommand:                  true,
		EnableNoteCommand:                  true,
		EnableCleanupCommand:               true,
		EnableForensicsCommand:             true,
		EnableExploreCommand:               true,
		EnableScanCommand:                  true,
		EnableImportCommand:                true,
		EnableAutonomousCommand:            true,
		EnableCodexPlugin:                  true,
		EnableGopls:                        true,
		EnablePyright:                      true,
		EnableTypeScript:                   true,
		EnableCSharp:                       true,
		EnableExecutor:                     true,
		EnableLibrarian:                    true,
		EnableReviewer:                     true,
		EnableResearcher:                   true,
		EnablePlanner:                      true,
		EnableDocWriter:                    true,
		EnableVerifier:                     true,
		EnableAgencyAIEngineer:             true,
		EnableAgencyBackendArchitect:       true,
		EnableAgencySecurityEngineer:       true,
		EnableAgencyCodeReviewerAgent:      true,
		EnableAgencySoftwareArchitect:      true,
		EnableAgencyDevOpsAutomator:        true,
		EnableAgencySRE:                    true,
		EnableAgencyAPITester:              true,
		EnableAgencyPerformanceBenchmarker: true,
		Settings: &schema.Settings{
			Model:                 "opus[1m]",
			AlwaysThinkingEnabled: &thinking,
			EffortLevel:           "high",
			AutoUpdatesChannel:    "latest",
		},
	}
}

// NewWithDefaults creates a Harness pre-loaded with built-in components.
func NewWithDefaults(opts DefaultOptions) *Harness {
	p := New()

	if opts.Settings != nil {
		p.SetSettings(opts.Settings)
	}

	// Handlers.
	if len(opts.LintProfiles) > 0 {
		p.Hooks().Register(handlers.NewLinterWith(opts.LintProfiles...))
	}
	if opts.EnableCommandGuard {
		p.Hooks().Register(handlers.NewCommandGuard())
	}
	if opts.EnableCommentChecker {
		p.Hooks().Register(handlers.NewCommentChecker())
	}
	if opts.EnableSecretRemediation {
		chain := hooks.NewChain(
			"secret-remediation",
			[]schema.HookEvent{schema.HookPostToolUse},
			regexp.MustCompile(`(?i)^(Edit|Write|MultiEdit)$`),
			hooks.HandlerLink(handlers.NewSecretScanner()),
			hooks.OnBlock(func(_ context.Context, _ *hooks.Input, prev *hooks.Result) (*hooks.Result, error) {
				prev.Output += "\n\nRemediation: Move the detected secret to an environment variable or a secrets manager."
				return prev, nil
			}),
		)
		p.Hooks().Register(chain)
	} else if opts.EnableSecretScanner {
		p.Hooks().Register(handlers.NewSecretScanner())
	}
	if opts.EnableSessionLogger {
		p.Hooks().Register(handlers.NewSessionLogger(""))
	}

	// Providers.
	if opts.EnableContext7 {
		p.MCP().Register(mcpproviders.NewContext7())
	}
	if opts.EnablePulumiMCP {
		p.MCP().Register(mcpproviders.NewPulumi())
	}
	if opts.EnableYaahMCP {
		p.MCP().Register(mcpproviders.NewYaah())
	}
	if opts.NotionToken != "" {
		p.MCP().Register(mcpproviders.NewNotion(opts.NotionToken))
	}

	// Workflow commands.
	if opts.EnableInitCommand {
		p.Commands().Register(cmdbuiltins.NewInitCommand())
	}
	if opts.EnableDiscussCommand {
		p.Commands().Register(cmdbuiltins.NewDiscussCommand())
	}
	if opts.EnablePlanCommand {
		p.Commands().Register(cmdbuiltins.NewPlanCommand())
	}
	if opts.EnableExecuteCommand {
		p.Commands().Register(cmdbuiltins.NewExecuteCommand())
	}
	if opts.EnableVerifyCommand {
		p.Commands().Register(cmdbuiltins.NewVerifyCommand())
	}
	if opts.EnableDocsCommand {
		p.Commands().Register(cmdbuiltins.NewDocsCommand())
	}
	if opts.EnableNextCommand {
		p.Commands().Register(cmdbuiltins.NewNextCommand())
	}
	if opts.EnableQuickCommand {
		p.Commands().Register(cmdbuiltins.NewQuickCommand())
	}
	if opts.EnableShipCommand {
		p.Commands().Register(cmdbuiltins.NewShipCommand())
	}
	if opts.EnablePauseCommand {
		p.Commands().Register(cmdbuiltins.NewPauseCommand())
	}
	if opts.EnableResumeCommand {
		p.Commands().Register(cmdbuiltins.NewResumeCommand())
	}
	if opts.EnableCompleteMilestoneCommand {
		p.Commands().Register(cmdbuiltins.NewCompleteMilestoneCommand())
	}
	if opts.EnableNewMilestoneCommand {
		p.Commands().Register(cmdbuiltins.NewNewMilestoneCommand())
	}
	if opts.EnableSettingsCommand {
		p.Commands().Register(cmdbuiltins.NewSettingsCommand())
	}
	if opts.EnableAddPhaseCommand {
		p.Commands().Register(cmdbuiltins.NewAddPhaseCommand())
	}
	if opts.EnableInsertPhaseCommand {
		p.Commands().Register(cmdbuiltins.NewInsertPhaseCommand())
	}
	if opts.EnableRemovePhaseCommand {
		p.Commands().Register(cmdbuiltins.NewRemovePhaseCommand())
	}
	if opts.EnableHealthCommand {
		p.Commands().Register(cmdbuiltins.NewHealthCommand())
	}
	if opts.EnableProgressCommand {
		p.Commands().Register(cmdbuiltins.NewProgressCommand())
	}
	if opts.EnableCodeReviewCommand {
		p.Commands().Register(cmdbuiltins.NewCodeReviewCommand())
	}
	if opts.EnableSecureCommand {
		p.Commands().Register(cmdbuiltins.NewSecureCommand())
	}
	if opts.EnableTodoCommand {
		p.Commands().Register(cmdbuiltins.NewTodoCommand())
	}
	if opts.EnableNoteCommand {
		p.Commands().Register(cmdbuiltins.NewNoteCommand())
	}
	if opts.EnableCleanupCommand {
		p.Commands().Register(cmdbuiltins.NewCleanupCommand())
	}
	if opts.EnableForensicsCommand {
		p.Commands().Register(cmdbuiltins.NewForensicsCommand())
	}
	if opts.EnableExploreCommand {
		p.Commands().Register(cmdbuiltins.NewExploreCommand())
	}
	if opts.EnableScanCommand {
		p.Commands().Register(cmdbuiltins.NewScanCommand())
	}
	if opts.EnableImportCommand {
		p.Commands().Register(cmdbuiltins.NewImportCommand())
	}
	if opts.EnableAutonomousCommand {
		p.Commands().Register(cmdbuiltins.NewAutonomousCommand())
	}

	// Skills.
	if opts.EnableCommitSkill {
		p.Skills().Register(builtins.NewCommitSkill())
	}
	if opts.EnablePRSkill {
		p.Skills().Register(builtins.NewPRSkill())
	}
	if opts.EnableReviewSkill {
		p.Skills().Register(builtins.NewReviewSkill())
	}

	// Remote skills — pulumi/agent-skills
	if opts.EnablePulumiBestPractices {
		p.Skills().Register(skills.NewRemoteSkill(
			"pulumi-best-practices", "Pulumi best practices for reliable programs",
			"github.com/pulumi/agent-skills@b6b942fc6e34517e2bbc52d6db04ca529baf3ad4", "authoring/skills/pulumi-best-practices/SKILL.md",
		))
	}
	if opts.EnablePulumiComponent {
		p.Skills().Register(skills.NewRemoteSkill(
			"pulumi-component", "Pulumi ComponentResource authoring",
			"github.com/pulumi/agent-skills@b6b942fc6e34517e2bbc52d6db04ca529baf3ad4", "authoring/skills/pulumi-component/SKILL.md",
		))
	}
	if opts.EnablePulumiAutomationAPI {
		p.Skills().Register(skills.NewRemoteSkill(
			"pulumi-automation-api", "Pulumi Automation API best practices",
			"github.com/pulumi/agent-skills@b6b942fc6e34517e2bbc52d6db04ca529baf3ad4", "authoring/skills/pulumi-automation-api/SKILL.md",
		))
	}
	if opts.EnablePulumiESC {
		p.Skills().Register(skills.NewRemoteSkill(
			"pulumi-esc", "Pulumi ESC environments, secrets, and configuration",
			"github.com/pulumi/agent-skills@b6b942fc6e34517e2bbc52d6db04ca529baf3ad4", "authoring/skills/pulumi-esc/SKILL.md",
		))
	}
	if opts.EnablePulumiTerraformMigrate {
		p.Skills().Register(skills.NewRemoteSkill(
			"pulumi-terraform-to-pulumi", "Convert Terraform to Pulumi",
			"github.com/pulumi/agent-skills@b6b942fc6e34517e2bbc52d6db04ca529baf3ad4", "migration/skills/pulumi-terraform-to-pulumi/SKILL.md",
		))
	}
	if opts.EnablePulumiCDKMigrate {
		p.Skills().Register(skills.NewRemoteSkill(
			"pulumi-cdk-to-pulumi", "Convert AWS CDK to Pulumi",
			"github.com/pulumi/agent-skills@b6b942fc6e34517e2bbc52d6db04ca529baf3ad4", "migration/skills/pulumi-cdk-to-pulumi/SKILL.md",
		))
	}
	if opts.EnablePulumiCFNMigrate {
		p.Skills().Register(skills.NewRemoteSkill(
			"cloudformation-to-pulumi", "Convert CloudFormation to Pulumi",
			"github.com/pulumi/agent-skills@b6b942fc6e34517e2bbc52d6db04ca529baf3ad4", "migration/skills/cloudformation-to-pulumi/SKILL.md",
		))
	}
	if opts.EnablePulumiARMMigrate {
		p.Skills().Register(skills.NewRemoteSkill(
			"pulumi-arm-to-pulumi", "Convert Azure ARM/Bicep to Pulumi",
			"github.com/pulumi/agent-skills@b6b942fc6e34517e2bbc52d6db04ca529baf3ad4", "migration/skills/pulumi-arm-to-pulumi/SKILL.md",
		))
	}

	// Remote skills — dirien/claude-skills
	if opts.EnablePulumiTypeScript {
		p.Skills().Register(skills.NewRemoteSkill(
			"pulumi-typescript", "Pulumi TypeScript IaC with ESC and OIDC",
			"github.com/dirien/claude-skills@85b0ee2a07cb1e3420d445d3f2336eadca45cde5", "pulumi-typescript/SKILL.md",
		))
	}
	if opts.EnablePulumiGo {
		p.Skills().Register(skills.NewRemoteSkill(
			"pulumi-go", "Pulumi Go IaC with ESC and OIDC",
			"github.com/dirien/claude-skills@85b0ee2a07cb1e3420d445d3f2336eadca45cde5", "pulumi-go/SKILL.md",
		))
	}
	if opts.EnablePulumiPython {
		p.Skills().Register(skills.NewRemoteSkill(
			"pulumi-python", "Pulumi Python IaC with ESC and OIDC",
			"github.com/dirien/claude-skills@85b0ee2a07cb1e3420d445d3f2336eadca45cde5", "pulumi-python/SKILL.md",
		))
	}
	if opts.EnablePulumiNeo {
		p.Skills().Register(skills.NewRemoteSkill(
			"pulumi-neo", "Pulumi Neo conversational infrastructure management",
			"github.com/dirien/claude-skills@85b0ee2a07cb1e3420d445d3f2336eadca45cde5", "pulumi-neo/SKILL.md",
		))
	}
	if opts.EnablePulumiCLI {
		p.Skills().Register(skills.NewRemoteSkill(
			"pulumi-cli", "Pulumi CLI command reference for deployments",
			"github.com/dirien/claude-skills@85b0ee2a07cb1e3420d445d3f2336eadca45cde5", "pulumi-cli/SKILL.md",
		))
	}
	if opts.EnableFluxCLI {
		p.Skills().Register(skills.NewRemoteSkill(
			"flux-cli", "GitOps for Kubernetes using Flux CD CLI",
			"github.com/dirien/claude-skills@85b0ee2a07cb1e3420d445d3f2336eadca45cde5", "flux-cli/SKILL.md",
		))
	}
	if opts.EnableFluxOperatorCLI {
		p.Skills().Register(skills.NewRemoteSkill(
			"flux-operator-cli", "Flux Operator CLI for managing Flux CD deployments on Kubernetes",
			"github.com/dirien/claude-skills@85b0ee2a07cb1e3420d445d3f2336eadca45cde5", "flux-operator-cli/SKILL.md",
		))
	}

	// Remote skills — jeffallan/claude-skills
	if opts.EnableGolangPro {
		p.Skills().Register(skills.NewRemoteSkill(
			"golang-pro", "Go concurrent patterns, microservices, gRPC, and performance optimization",
			"github.com/jeffallan/claude-skills@3bf9a24b76a7c122f1fc05e83929fbc84e1c207a", "skills/golang-pro/SKILL.md",
		))
	}
	if opts.EnableKubernetesSpecialist {
		p.Skills().Register(skills.NewRemoteSkill(
			"kubernetes-specialist", "Kubernetes deployments, Helm, RBAC, NetworkPolicies, and multi-cluster",
			"github.com/jeffallan/claude-skills@3bf9a24b76a7c122f1fc05e83929fbc84e1c207a", "skills/kubernetes-specialist/SKILL.md",
		))
	}
	if opts.EnableDevOpsEngineer {
		p.Skills().Register(skills.NewRemoteSkill(
			"devops-engineer", "CI/CD pipelines, Docker, Kubernetes, Terraform, and GitOps",
			"github.com/jeffallan/claude-skills@3bf9a24b76a7c122f1fc05e83929fbc84e1c207a", "skills/devops-engineer/SKILL.md",
		))
	}
	if opts.EnablePythonPro {
		p.Skills().Register(skills.NewRemoteSkill(
			"python-pro", "Python 3.11+ with type safety, async, pytest, and ruff",
			"github.com/jeffallan/claude-skills@3bf9a24b76a7c122f1fc05e83929fbc84e1c207a", "skills/python-pro/SKILL.md",
		))
	}
	if opts.EnableTypeScriptPro {
		p.Skills().Register(skills.NewRemoteSkill(
			"typescript-pro", "Advanced TypeScript types, generics, tRPC, and monorepo setup",
			"github.com/jeffallan/claude-skills@3bf9a24b76a7c122f1fc05e83929fbc84e1c207a", "skills/typescript-pro/SKILL.md",
		))
	}
	if opts.EnableCSharpDeveloper {
		p.Skills().Register(skills.NewRemoteSkill(
			"csharp-developer", "C# .NET 8+, ASP.NET Core, Blazor, EF Core, and MediatR",
			"github.com/jeffallan/claude-skills@3bf9a24b76a7c122f1fc05e83929fbc84e1c207a", "skills/csharp-developer/SKILL.md",
		))
	}
	if opts.EnableJavaScriptPro {
		p.Skills().Register(skills.NewRemoteSkill(
			"javascript-pro", "Modern ES2023+ JavaScript, async/await, ESM, and Node.js",
			"github.com/jeffallan/claude-skills@3bf9a24b76a7c122f1fc05e83929fbc84e1c207a", "skills/javascript-pro/SKILL.md",
		))
	}
	if opts.EnableCLIDeveloper {
		p.Skills().Register(skills.NewRemoteSkill(
			"cli-developer", "CLI tools with argument parsing, completions, and cross-platform support",
			"github.com/jeffallan/claude-skills@3bf9a24b76a7c122f1fc05e83929fbc84e1c207a", "skills/cli-developer/SKILL.md",
		))
	}
	if opts.EnableSREEngineer {
		p.Skills().Register(skills.NewRemoteSkill(
			"sre-engineer", "SLOs, error budgets, incident response, and capacity planning",
			"github.com/jeffallan/claude-skills@3bf9a24b76a7c122f1fc05e83929fbc84e1c207a", "skills/sre-engineer/SKILL.md",
		))
	}
	if opts.EnableTheFool {
		p.Skills().Register(skills.NewRemoteSkill(
			"the-fool", "Devil's advocate, pre-mortems, red teaming, and assumption auditing",
			"github.com/jeffallan/claude-skills@3bf9a24b76a7c122f1fc05e83929fbc84e1c207a", "skills/the-fool/SKILL.md",
		))
	}
	if opts.EnableArchitectureDesigner {
		p.Skills().Register(skills.NewRemoteSkill(
			"architecture-designer", "System architecture, ADRs, trade-offs, and scalability planning",
			"github.com/jeffallan/claude-skills@3bf9a24b76a7c122f1fc05e83929fbc84e1c207a", "skills/architecture-designer/SKILL.md",
		))
	}
	if opts.EnableSpringBootEngineer {
		p.Skills().Register(skills.NewRemoteSkill(
			"spring-boot-engineer", "Spring Boot 3.x, Spring Security 6, JPA, WebFlux, and Spring Cloud",
			"github.com/jeffallan/claude-skills@3bf9a24b76a7c122f1fc05e83929fbc84e1c207a", "skills/spring-boot-engineer/SKILL.md",
		))
	}
	if opts.EnableCodeReviewer {
		p.Skills().Register(skills.NewRemoteSkill(
			"code-reviewer", "Code review for bugs, security, performance, and maintainability",
			"github.com/jeffallan/claude-skills@3bf9a24b76a7c122f1fc05e83929fbc84e1c207a", "skills/code-reviewer/SKILL.md",
		))
	}

	// Remote skills — Rust
	if opts.EnableRustBestPractices {
		p.Skills().Register(skills.NewRemoteSkill(
			"rust-best-practices", "Idiomatic Rust code, borrowing, error handling, and performance optimization",
			"github.com/apollographql/skills@e1979d2f1e7c38cef58753b2bfd6fc9509101bdc", "skills/rust-best-practices/SKILL.md",
		))
	}
	if opts.EnableRustAsyncPatterns {
		p.Skills().Register(skills.NewRemoteSkill(
			"rust-async-patterns", "Rust async programming with Tokio, async traits, and concurrent patterns",
			"github.com/wshobson/agents@1ad2f007d5e9ec822a2d79e727ac1dcdf5f66f11", "plugins/systems-programming/skills/rust-async-patterns/SKILL.md",
		))
	}
	if opts.EnableRustEngineer {
		p.Skills().Register(skills.NewRemoteSkill(
			"rust-engineer", "Idiomatic Rust with ownership, lifetimes, traits, tokio, and error handling",
			"github.com/jeffallan/claude-skills@3bf9a24b76a7c122f1fc05e83929fbc84e1c207a", "skills/rust-engineer/SKILL.md",
		))
	}

	// Remote skills — netresearch/agent-rules-skill
	if opts.EnableAgentRules {
		p.Skills().Register(skills.NewRemoteSkill(
			"agent-rules", "Generate and maintain AGENTS.md files following the agents.md convention",
			"github.com/netresearch/agent-rules-skill@9b67bf594a52b1a7d38d8b0ec0a076a31f8d3d7e", "skills/agent-rules/SKILL.md",
		))
	}

	// Remote skills — rshade/agent-skills
	const rshadeRef = "github.com/rshade/agent-skills@4aff11fe89bb156337c2c7c303bb2db234cc9740"
	if opts.EnableAgentReadyGo {
		p.Skills().Register(skills.NewRemoteSkill(
			"agent-ready-go", "Prepare Go apps for AI agent interaction with structured logging and CLI design",
			rshadeRef, "skills/agent-ready-go/SKILL.md",
		))
	}
	if opts.EnableCommitlint {
		p.Skills().Register(skills.NewRemoteSkill(
			"commitlint", "Validate commit messages against Conventional Commits specification",
			rshadeRef, "skills/commitlint/SKILL.md",
		))
	}
	if opts.EnableDecide {
		p.Skills().Register(skills.NewRemoteSkill(
			"decide", "Three-agent adversarial debate protocol for strategic decisions",
			rshadeRef, "skills/decide/SKILL.md",
		))
	}
	if opts.EnableDepUpgrade {
		p.Skills().Register(skills.NewRemoteSkill(
			"dep-upgrade", "Safe systematic dependency upgrades with vulnerability scanning and rollback",
			rshadeRef, "skills/dep-upgrade/SKILL.md",
		))
	}
	if opts.EnableDesignPrinciples {
		p.Skills().Register(skills.NewRemoteSkill(
			"design-principles", "Analyze codebases against SOLID, DRY, YAGNI, KISS, and other design principles",
			rshadeRef, "skills/design-principles/SKILL.md",
		))
	}
	if opts.EnableGoNolintAudit {
		p.Skills().Register(skills.NewRemoteSkill(
			"go-nolint-audit", "Audit nolint directives in Go codebases for stale or unjustified suppressions",
			rshadeRef, "skills/go-nolint-audit/SKILL.md",
		))
	}
	if opts.EnableLintFix {
		p.Skills().Register(skills.NewRemoteSkill(
			"lint-fix", "Detect linting tools, run them to zero errors, and fix issues atomically",
			rshadeRef, "skills/lint-fix/SKILL.md",
		))
	}
	if opts.EnableMarkdownlint {
		p.Skills().Register(skills.NewRemoteSkill(
			"markdownlint", "Validate markdown files against formatting standards with auto-fix",
			rshadeRef, "skills/markdownlint/SKILL.md",
		))
	}
	if opts.EnablePullRequestMsg {
		p.Skills().Register(skills.NewRemoteSkill(
			"pull-request-msg", "Generate structured PR descriptions from session context using GitHub CLI",
			rshadeRef, "skills/pull-request-msg-with-gh/SKILL.md",
		))
	}
	if opts.EnableRoadmap {
		p.Skills().Register(skills.NewRemoteSkill(
			"roadmap", "Strategic roadmap management synced with GitHub Issues and labels",
			rshadeRef, "skills/roadmap/SKILL.md",
		))
	}
	if opts.EnableScout {
		p.Skills().Register(skills.NewRemoteSkill(
			"scout", "Identify top improvement opportunities in files you are touching",
			rshadeRef, "skills/scout/SKILL.md",
		))
	}
	if opts.EnableSecurityAudit {
		p.Skills().Register(skills.NewRemoteSkill(
			"security-audit", "Comprehensive vulnerability assessment with OWASP Top 10 and threat modeling",
			rshadeRef, "skills/security-audit/SKILL.md",
		))
	}
	if opts.EnableTailscaleInstall {
		p.Skills().Register(skills.NewRemoteSkill(
			"tailscale-install", "Install and configure Tailscale across platforms including WSL2 and containers",
			rshadeRef, "skills/tailscale-install/SKILL.md",
		))
	}
	if opts.EnableTechDebt {
		p.Skills().Register(skills.NewRemoteSkill(
			"tech-debt", "Systematic technical debt analysis across 9 categories with health scoring",
			rshadeRef, "skills/tech-debt/SKILL.md",
		))
	}

	// Plugins.
	if opts.EnableCodexPlugin {
		p.Plugins().Register(plugins.NewCodex())
	}

	// LSP servers.
	if opts.EnableGopls {
		p.LSP().Register(lspproviders.NewGopls())
	}
	if opts.EnablePyright {
		p.LSP().Register(lspproviders.NewPyright())
	}
	if opts.EnableTypeScript {
		p.LSP().Register(lspproviders.NewTypeScript())
	}
	if opts.EnableCSharp {
		p.LSP().Register(lspproviders.NewCSharp())
	}

	// Agents.
	if opts.EnableExecutor {
		p.Agents().Register(agentpkg.NewExecutor())
	}
	if opts.EnableLibrarian {
		p.Agents().Register(agentpkg.NewLibrarian())
	}
	if opts.EnableReviewer {
		p.Agents().Register(agentpkg.NewReviewer())
	}

	// Workflow agents.
	if opts.EnableResearcher {
		p.Agents().Register(agentpkg.NewResearcher())
	}
	if opts.EnablePlanner {
		p.Agents().Register(agentpkg.NewPlanner())
	}
	if opts.EnableDocWriter {
		p.Agents().Register(agentpkg.NewDocWriter())
	}
	if opts.EnableVerifier {
		p.Agents().Register(agentpkg.NewVerifier())
	}

	// Remote agents — msitarzewski/agency-agents
	const agencyAgentsRef = "github.com/msitarzewski/agency-agents@6254154899f510eb4a4de10561fecfc1f32ff17f"
	if opts.EnableAgencyAIEngineer {
		p.Agents().Register(agentpkg.NewRemoteAgent(
			"agency-ai-engineer", "AI/ML engineering, model integration, LLM pipelines, and AI system design",
			agencyAgentsRef, "engineering/engineering-ai-engineer.md",
			agentpkg.WithModel("sonnet"),
		))
	}
	if opts.EnableAgencyBackendArchitect {
		p.Agents().Register(agentpkg.NewRemoteAgent(
			"agency-backend-architect", "Backend system design, API architecture, and scalability planning",
			agencyAgentsRef, "engineering/engineering-backend-architect.md",
			agentpkg.WithModel("sonnet"),
		))
	}
	if opts.EnableAgencySecurityEngineer {
		p.Agents().Register(agentpkg.NewRemoteAgent(
			"agency-security-engineer", "Security analysis, threat modeling, and vulnerability assessment",
			agencyAgentsRef, "engineering/engineering-security-engineer.md",
			agentpkg.WithModel("sonnet"),
		))
	}
	if opts.EnableAgencyCodeReviewerAgent {
		p.Agents().Register(agentpkg.NewRemoteAgent(
			"agency-code-reviewer", "Structured code review with quality, security, and performance analysis",
			agencyAgentsRef, "engineering/engineering-code-reviewer.md",
			agentpkg.WithModel("sonnet"),
			agentpkg.WithTools("Read, Grep, Glob"),
		))
	}
	if opts.EnableAgencySoftwareArchitect {
		p.Agents().Register(agentpkg.NewRemoteAgent(
			"agency-software-architect", "System architecture, design patterns, and technical decision-making",
			agencyAgentsRef, "engineering/engineering-software-architect.md",
			agentpkg.WithModel("opus"),
		))
	}
	if opts.EnableAgencyDevOpsAutomator {
		p.Agents().Register(agentpkg.NewRemoteAgent(
			"agency-devops-automator", "CI/CD pipelines, infrastructure automation, and deployment workflows",
			agencyAgentsRef, "engineering/engineering-devops-automator.md",
			agentpkg.WithModel("sonnet"),
		))
	}
	if opts.EnableAgencySRE {
		p.Agents().Register(agentpkg.NewRemoteAgent(
			"agency-sre", "Site reliability engineering, observability, and incident response",
			agencyAgentsRef, "engineering/engineering-sre.md",
			agentpkg.WithModel("sonnet"),
		))
	}
	if opts.EnableAgencyAPITester {
		p.Agents().Register(agentpkg.NewRemoteAgent(
			"agency-api-tester", "API testing, contract validation, and endpoint coverage analysis",
			agencyAgentsRef, "testing/testing-api-tester.md",
			agentpkg.WithModel("sonnet"),
			agentpkg.WithTools("Read, Grep, Glob, Bash(*), WebFetch"),
		))
	}
	if opts.EnableAgencyPerformanceBenchmarker {
		p.Agents().Register(agentpkg.NewRemoteAgent(
			"agency-performance-benchmarker", "Performance profiling, load testing, and optimization recommendations",
			agencyAgentsRef, "testing/testing-performance-benchmarker.md",
			agentpkg.WithModel("sonnet"),
			agentpkg.WithTools("Read, Grep, Glob, Bash(*)"),
		))
	}

	return p
}

// NewFromCatalog creates a Harness using the catalog to resolve skills.
// Skills are selected by SkillIDs and BundleIDs from the default catalog.
// ExcludeIDs removes specific skills from the resolved set.
// Non-skill components (hooks, MCP, LSP, agents) are configured via the same
// boolean flags as NewWithDefaults.
func NewFromCatalog(opts DefaultOptions) *Harness {
	p := New()

	if opts.Settings != nil {
		p.SetSettings(opts.Settings)
	}

	// Handlers — same as NewWithDefaults.
	if len(opts.LintProfiles) > 0 {
		p.Hooks().Register(handlers.NewLinterWith(opts.LintProfiles...))
	}
	if opts.EnableCommandGuard {
		p.Hooks().Register(handlers.NewCommandGuard())
	}
	if opts.EnableCommentChecker {
		p.Hooks().Register(handlers.NewCommentChecker())
	}
	if opts.EnableSecretRemediation {
		chain := hooks.NewChain(
			"secret-remediation",
			[]schema.HookEvent{schema.HookPostToolUse},
			regexp.MustCompile(`(?i)^(Edit|Write|MultiEdit)$`),
			hooks.HandlerLink(handlers.NewSecretScanner()),
			hooks.OnBlock(func(_ context.Context, _ *hooks.Input, prev *hooks.Result) (*hooks.Result, error) {
				prev.Output += "\n\nRemediation: Move the detected secret to an environment variable or a secrets manager."
				return prev, nil
			}),
		)
		p.Hooks().Register(chain)
	} else if opts.EnableSecretScanner {
		p.Hooks().Register(handlers.NewSecretScanner())
	}
	if opts.EnableSessionLogger {
		p.Hooks().Register(handlers.NewSessionLogger(""))
	}

	// Providers — same as NewWithDefaults.
	if opts.EnableContext7 {
		p.MCP().Register(mcpproviders.NewContext7())
	}
	if opts.EnablePulumiMCP {
		p.MCP().Register(mcpproviders.NewPulumi())
	}
	if opts.EnableYaahMCP {
		p.MCP().Register(mcpproviders.NewYaah())
	}
	if opts.NotionToken != "" {
		p.MCP().Register(mcpproviders.NewNotion(opts.NotionToken))
	}

	// Workflow commands.
	if opts.EnableInitCommand {
		p.Commands().Register(cmdbuiltins.NewInitCommand())
	}
	if opts.EnableDiscussCommand {
		p.Commands().Register(cmdbuiltins.NewDiscussCommand())
	}
	if opts.EnablePlanCommand {
		p.Commands().Register(cmdbuiltins.NewPlanCommand())
	}
	if opts.EnableExecuteCommand {
		p.Commands().Register(cmdbuiltins.NewExecuteCommand())
	}
	if opts.EnableVerifyCommand {
		p.Commands().Register(cmdbuiltins.NewVerifyCommand())
	}
	if opts.EnableDocsCommand {
		p.Commands().Register(cmdbuiltins.NewDocsCommand())
	}
	if opts.EnableNextCommand {
		p.Commands().Register(cmdbuiltins.NewNextCommand())
	}
	if opts.EnableQuickCommand {
		p.Commands().Register(cmdbuiltins.NewQuickCommand())
	}
	if opts.EnableShipCommand {
		p.Commands().Register(cmdbuiltins.NewShipCommand())
	}
	if opts.EnablePauseCommand {
		p.Commands().Register(cmdbuiltins.NewPauseCommand())
	}
	if opts.EnableResumeCommand {
		p.Commands().Register(cmdbuiltins.NewResumeCommand())
	}
	if opts.EnableCompleteMilestoneCommand {
		p.Commands().Register(cmdbuiltins.NewCompleteMilestoneCommand())
	}
	if opts.EnableNewMilestoneCommand {
		p.Commands().Register(cmdbuiltins.NewNewMilestoneCommand())
	}
	if opts.EnableSettingsCommand {
		p.Commands().Register(cmdbuiltins.NewSettingsCommand())
	}
	if opts.EnableAddPhaseCommand {
		p.Commands().Register(cmdbuiltins.NewAddPhaseCommand())
	}
	if opts.EnableInsertPhaseCommand {
		p.Commands().Register(cmdbuiltins.NewInsertPhaseCommand())
	}
	if opts.EnableRemovePhaseCommand {
		p.Commands().Register(cmdbuiltins.NewRemovePhaseCommand())
	}
	if opts.EnableHealthCommand {
		p.Commands().Register(cmdbuiltins.NewHealthCommand())
	}
	if opts.EnableProgressCommand {
		p.Commands().Register(cmdbuiltins.NewProgressCommand())
	}
	if opts.EnableCodeReviewCommand {
		p.Commands().Register(cmdbuiltins.NewCodeReviewCommand())
	}
	if opts.EnableSecureCommand {
		p.Commands().Register(cmdbuiltins.NewSecureCommand())
	}
	if opts.EnableTodoCommand {
		p.Commands().Register(cmdbuiltins.NewTodoCommand())
	}
	if opts.EnableNoteCommand {
		p.Commands().Register(cmdbuiltins.NewNoteCommand())
	}
	if opts.EnableCleanupCommand {
		p.Commands().Register(cmdbuiltins.NewCleanupCommand())
	}
	if opts.EnableForensicsCommand {
		p.Commands().Register(cmdbuiltins.NewForensicsCommand())
	}
	if opts.EnableExploreCommand {
		p.Commands().Register(cmdbuiltins.NewExploreCommand())
	}
	if opts.EnableScanCommand {
		p.Commands().Register(cmdbuiltins.NewScanCommand())
	}
	if opts.EnableImportCommand {
		p.Commands().Register(cmdbuiltins.NewImportCommand())
	}
	if opts.EnableAutonomousCommand {
		p.Commands().Register(cmdbuiltins.NewAutonomousCommand())
	}

	// Skills — catalog-based.
	cat := catalog.DefaultCatalog()
	wanted := resolveSkillIDs(cat, opts.SkillIDs, opts.BundleIDs, opts.ExcludeIDs)

	for _, id := range wanted {
		entry := cat.ByID(id)
		if entry == nil {
			continue
		}
		meta := skills.SkillMetadata{
			Category: string(entry.Category),
			Tags:     entry.Tags,
			Risk:     string(entry.Risk),
			Tier:     string(entry.Tier),
			Aliases:  entry.Aliases,
		}
		if entry.Uses == "builtin" {
			switch id {
			case "commit":
				p.Skills().Register(builtins.NewCommitSkill())
			case "pr":
				p.Skills().Register(builtins.NewPRSkill())
			case "review":
				p.Skills().Register(builtins.NewReviewSkill())
			}
			continue
		}
		p.Skills().Register(skills.NewRemoteSkill(
			entry.ID, entry.Description, entry.Uses, entry.Subpath,
			skills.WithMetadata(meta),
		))
	}

	// Plugins — same as NewWithDefaults.
	if opts.EnableCodexPlugin {
		p.Plugins().Register(plugins.NewCodex())
	}

	// LSP servers — same as NewWithDefaults.
	if opts.EnableGopls {
		p.LSP().Register(lspproviders.NewGopls())
	}
	if opts.EnablePyright {
		p.LSP().Register(lspproviders.NewPyright())
	}
	if opts.EnableTypeScript {
		p.LSP().Register(lspproviders.NewTypeScript())
	}
	if opts.EnableCSharp {
		p.LSP().Register(lspproviders.NewCSharp())
	}

	// Agents — same as NewWithDefaults.
	if opts.EnableExecutor {
		p.Agents().Register(agentpkg.NewExecutor())
	}
	if opts.EnableLibrarian {
		p.Agents().Register(agentpkg.NewLibrarian())
	}
	if opts.EnableReviewer {
		p.Agents().Register(agentpkg.NewReviewer())
	}

	// Workflow agents.
	if opts.EnableResearcher {
		p.Agents().Register(agentpkg.NewResearcher())
	}
	if opts.EnablePlanner {
		p.Agents().Register(agentpkg.NewPlanner())
	}
	if opts.EnableDocWriter {
		p.Agents().Register(agentpkg.NewDocWriter())
	}
	if opts.EnableVerifier {
		p.Agents().Register(agentpkg.NewVerifier())
	}

	const agencyCatalogRef = "github.com/msitarzewski/agency-agents@6254154899f510eb4a4de10561fecfc1f32ff17f"
	if opts.EnableAgencyAIEngineer {
		p.Agents().Register(agentpkg.NewRemoteAgent(
			"agency-ai-engineer", "AI/ML engineering, model integration, LLM pipelines, and AI system design",
			agencyCatalogRef, "engineering/engineering-ai-engineer.md",
			agentpkg.WithModel("sonnet"),
		))
	}
	if opts.EnableAgencyBackendArchitect {
		p.Agents().Register(agentpkg.NewRemoteAgent(
			"agency-backend-architect", "Backend system design, API architecture, and scalability planning",
			agencyCatalogRef, "engineering/engineering-backend-architect.md",
			agentpkg.WithModel("sonnet"),
		))
	}
	if opts.EnableAgencySecurityEngineer {
		p.Agents().Register(agentpkg.NewRemoteAgent(
			"agency-security-engineer", "Security analysis, threat modeling, and vulnerability assessment",
			agencyCatalogRef, "engineering/engineering-security-engineer.md",
			agentpkg.WithModel("sonnet"),
		))
	}
	if opts.EnableAgencyCodeReviewerAgent {
		p.Agents().Register(agentpkg.NewRemoteAgent(
			"agency-code-reviewer", "Structured code review with quality, security, and performance analysis",
			agencyCatalogRef, "engineering/engineering-code-reviewer.md",
			agentpkg.WithModel("sonnet"),
			agentpkg.WithTools("Read, Grep, Glob"),
		))
	}
	if opts.EnableAgencySoftwareArchitect {
		p.Agents().Register(agentpkg.NewRemoteAgent(
			"agency-software-architect", "System architecture, design patterns, and technical decision-making",
			agencyCatalogRef, "engineering/engineering-software-architect.md",
			agentpkg.WithModel("opus"),
		))
	}
	if opts.EnableAgencyDevOpsAutomator {
		p.Agents().Register(agentpkg.NewRemoteAgent(
			"agency-devops-automator", "CI/CD pipelines, infrastructure automation, and deployment workflows",
			agencyCatalogRef, "engineering/engineering-devops-automator.md",
			agentpkg.WithModel("sonnet"),
		))
	}
	if opts.EnableAgencySRE {
		p.Agents().Register(agentpkg.NewRemoteAgent(
			"agency-sre", "Site reliability engineering, observability, and incident response",
			agencyCatalogRef, "engineering/engineering-sre.md",
			agentpkg.WithModel("sonnet"),
		))
	}
	if opts.EnableAgencyAPITester {
		p.Agents().Register(agentpkg.NewRemoteAgent(
			"agency-api-tester", "API testing, contract validation, and endpoint coverage analysis",
			agencyCatalogRef, "testing/testing-api-tester.md",
			agentpkg.WithModel("sonnet"),
			agentpkg.WithTools("Read, Grep, Glob, Bash(*), WebFetch"),
		))
	}
	if opts.EnableAgencyPerformanceBenchmarker {
		p.Agents().Register(agentpkg.NewRemoteAgent(
			"agency-performance-benchmarker", "Performance profiling, load testing, and optimization recommendations",
			agencyCatalogRef, "testing/testing-performance-benchmarker.md",
			agentpkg.WithModel("sonnet"),
			agentpkg.WithTools("Read, Grep, Glob, Bash(*)"),
		))
	}

	return p
}

// resolveSkillIDs builds a deduplicated, ordered list of skill IDs from
// explicit IDs, bundle expansions, and exclusions.
func resolveSkillIDs(cat *catalog.Catalog, ids, bundles, exclude []string) []string {
	seen := make(map[string]bool)
	var result []string

	add := func(id string) {
		if !seen[id] && !slices.Contains(exclude, id) {
			seen[id] = true
			result = append(result, id)
		}
	}

	for _, id := range ids {
		add(id)
	}
	for _, bid := range bundles {
		for _, entry := range cat.ResolveBundle(bid) {
			add(entry.ID)
		}
	}
	return result
}
