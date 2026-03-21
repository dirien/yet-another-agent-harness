package harness

import (
	"context"
	"regexp"

	agentpkg "github.com/dirien/yet-another-agent-harness/pkg/agents"
	"github.com/dirien/yet-another-agent-harness/pkg/hooks"
	"github.com/dirien/yet-another-agent-harness/pkg/hooks/handlers"
	lspproviders "github.com/dirien/yet-another-agent-harness/pkg/lsp/providers"
	mcpproviders "github.com/dirien/yet-another-agent-harness/pkg/mcp/providers"
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

	// LSP servers (marketplace-backed)
	EnableGopls      bool
	EnablePyright    bool
	EnableTypeScript bool
	EnableCSharp     bool

	// Agents
	EnableExecutor  bool
	EnableLibrarian bool
	EnableReviewer  bool

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
		EnableGopls:                        true,
		EnablePyright:                      true,
		EnableTypeScript:                   true,
		EnableCSharp:                       true,
		EnableExecutor:                     true,
		EnableLibrarian:                    true,
		EnableReviewer:                     true,
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
			"github.com/dirien/claude-skills@073664d4f8e83fc1447a9e310dab6c51482f64bf", "pulumi-typescript/SKILL.md",
		))
	}
	if opts.EnablePulumiGo {
		p.Skills().Register(skills.NewRemoteSkill(
			"pulumi-go", "Pulumi Go IaC with ESC and OIDC",
			"github.com/dirien/claude-skills@073664d4f8e83fc1447a9e310dab6c51482f64bf", "pulumi-go/SKILL.md",
		))
	}
	if opts.EnablePulumiPython {
		p.Skills().Register(skills.NewRemoteSkill(
			"pulumi-python", "Pulumi Python IaC with ESC and OIDC",
			"github.com/dirien/claude-skills@073664d4f8e83fc1447a9e310dab6c51482f64bf", "pulumi-python/SKILL.md",
		))
	}
	if opts.EnablePulumiNeo {
		p.Skills().Register(skills.NewRemoteSkill(
			"pulumi-neo", "Pulumi Neo conversational infrastructure management",
			"github.com/dirien/claude-skills@073664d4f8e83fc1447a9e310dab6c51482f64bf", "pulumi-neo/SKILL.md",
		))
	}
	if opts.EnablePulumiCLI {
		p.Skills().Register(skills.NewRemoteSkill(
			"pulumi-cli", "Pulumi CLI command reference for deployments",
			"github.com/dirien/claude-skills@073664d4f8e83fc1447a9e310dab6c51482f64bf", "pulumi-cli/SKILL.md",
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
