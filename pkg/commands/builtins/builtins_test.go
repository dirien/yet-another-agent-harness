package builtins_test

import (
	"strings"
	"testing"

	"github.com/dirien/yet-another-agent-harness/pkg/commands"
	"github.com/dirien/yet-another-agent-harness/pkg/commands/builtins"
)

// allCommands returns every workflow command for testing.
func allCommands() []commands.Command {
	return []commands.Command{
		builtins.NewInitCommand(),
		builtins.NewDiscussCommand(),
		builtins.NewPlanCommand(),
		builtins.NewExecuteCommand(),
		builtins.NewVerifyCommand(),
		builtins.NewDocsCommand(),
		builtins.NewNextCommand(),
		builtins.NewQuickCommand(),
		builtins.NewShipCommand(),
		builtins.NewPauseCommand(),
		builtins.NewResumeCommand(),
		builtins.NewCompleteMilestoneCommand(),
		builtins.NewNewMilestoneCommand(),
		builtins.NewSettingsCommand(),
		builtins.NewAddPhaseCommand(),
		builtins.NewInsertPhaseCommand(),
		builtins.NewRemovePhaseCommand(),
		builtins.NewHealthCommand(),
		builtins.NewProgressCommand(),
		builtins.NewCodeReviewCommand(),
		builtins.NewSecureCommand(),
		builtins.NewTodoCommand(),
		builtins.NewNoteCommand(),
		builtins.NewCleanupCommand(),
		builtins.NewForensicsCommand(),
		builtins.NewExploreCommand(),
		builtins.NewScanCommand(),
		builtins.NewImportCommand(),
		builtins.NewAutonomousCommand(),
	}
}

func TestAllCommandsImplementInterface(t *testing.T) {
	cmds := allCommands()
	if len(cmds) != 29 {
		t.Fatalf("expected 29 commands, got %d", len(cmds))
	}

	for _, c := range cmds {
		t.Run(c.Name(), func(t *testing.T) {
			if c.Name() == "" {
				t.Error("Name() must be non-empty")
			}
			if c.Description() == "" {
				t.Error("Description() must be non-empty")
			}
			if c.Content() == "" {
				t.Error("Content() must be non-empty")
			}
		})
	}
}

func TestWorkflowCommandsHaveForkContext(t *testing.T) {
	// These commands should run as subagents (context: fork).
	heavy := []commands.Command{
		builtins.NewInitCommand(),
		builtins.NewDiscussCommand(),
		builtins.NewPlanCommand(),
		builtins.NewExecuteCommand(),
		builtins.NewVerifyCommand(),
		builtins.NewDocsCommand(),
		builtins.NewQuickCommand(),
		builtins.NewShipCommand(),
		builtins.NewPauseCommand(),
		builtins.NewResumeCommand(),
		builtins.NewCompleteMilestoneCommand(),
		builtins.NewNewMilestoneCommand(),
		builtins.NewAddPhaseCommand(),
		builtins.NewInsertPhaseCommand(),
		builtins.NewRemovePhaseCommand(),
		builtins.NewCodeReviewCommand(),
		builtins.NewSecureCommand(),
		builtins.NewCleanupCommand(),
		builtins.NewForensicsCommand(),
		builtins.NewExploreCommand(),
		builtins.NewScanCommand(),
		builtins.NewImportCommand(),
		builtins.NewAutonomousCommand(),
	}

	for _, c := range heavy {
		t.Run(c.Name(), func(t *testing.T) {
			adv, ok := c.(commands.CommandWithAdvanced)
			if !ok {
				t.Fatalf("command %q must implement CommandWithAdvanced", c.Name())
			}
			if adv.Advanced().Context != "fork" {
				t.Errorf("command %q: Advanced().Context = %q, want %q", c.Name(), adv.Advanced().Context, "fork")
			}
		})
	}
}

func TestLightweightCommandsNoFork(t *testing.T) {
	// These commands are lightweight and should NOT fork.
	lightweight := []commands.Command{
		builtins.NewNextCommand(),
		builtins.NewSettingsCommand(),
		builtins.NewHealthCommand(),
		builtins.NewProgressCommand(),
		builtins.NewTodoCommand(),
		builtins.NewNoteCommand(),
	}

	for _, c := range lightweight {
		t.Run(c.Name(), func(t *testing.T) {
			if _, ok := c.(commands.CommandWithAdvanced); ok {
				t.Errorf("command %q should NOT implement CommandWithAdvanced (lightweight)", c.Name())
			}
		})
	}
}

func TestCommandContentContainsKeyInstructions(t *testing.T) {
	tests := []struct {
		name     string
		cmd      commands.Command
		keywords []string
	}{
		{"yaah/init", builtins.NewInitCommand(), []string{".planning/", "PROJECT.md", "ROADMAP.md", "STATE.md"}},
		{"yaah/discuss", builtins.NewDiscussCommand(), []string{"CONTEXT.md", "gray area", "decision"}},
		{"yaah/plan", builtins.NewPlanCommand(), []string{"PLAN.md", "wave"}},
		{"yaah/execute", builtins.NewExecuteCommand(), []string{"wave", "SUMMARY.md"}},
		{"yaah/verify", builtins.NewVerifyCommand(), []string{"VERIFICATION.md"}},
		{"yaah/docs", builtins.NewDocsCommand(), []string{"README.md", "ARCHITECTURE.md"}},
		{"yaah/next", builtins.NewNextCommand(), []string{".planning/", "STATE.md"}},
		{"yaah/quick", builtins.NewQuickCommand(), []string{"--discuss", "--research", "--validate"}},
		{"yaah/ship", builtins.NewShipCommand(), []string{"VERIFICATION.md", "gh pr create"}},
		{"yaah/pause", builtins.NewPauseCommand(), []string{"HANDOFF"}},
		{"yaah/resume", builtins.NewResumeCommand(), []string{"HANDOFF"}},
		{"yaah/complete-milestone", builtins.NewCompleteMilestoneCommand(), []string{"CHANGELOG", "tag"}},
		{"yaah/new-milestone", builtins.NewNewMilestoneCommand(), []string{"milestone", "REQUIREMENTS"}},
		{"yaah/settings", builtins.NewSettingsCommand(), []string{"config.json"}},
		{"yaah/add-phase", builtins.NewAddPhaseCommand(), []string{"ROADMAP.md", "phase"}},
		{"yaah/insert-phase", builtins.NewInsertPhaseCommand(), []string{"renumber", "phase"}},
		{"yaah/remove-phase", builtins.NewRemovePhaseCommand(), []string{"ROADMAP.md", "phase"}},
		{"yaah/health", builtins.NewHealthCommand(), []string{".planning/", "PROJECT.md"}},
		{"yaah/progress", builtins.NewProgressCommand(), []string{"STATE.md", "ROADMAP.md"}},
		{"yaah/review", builtins.NewCodeReviewCommand(), []string{"REVIEW.md", "CRITICAL"}},
		{"yaah/secure", builtins.NewSecureCommand(), []string{"STRIDE", "SECURITY.md"}},
		{"yaah/todo", builtins.NewTodoCommand(), []string{"TODOS.md"}},
		{"yaah/note", builtins.NewNoteCommand(), []string{".planning/notes/"}},
		{"yaah/cleanup", builtins.NewCleanupCommand(), []string{"HANDOFF", "dry-run"}},
		{"yaah/forensics", builtins.NewForensicsCommand(), []string{"FORENSICS", "Recovery"}},
		{"yaah/explore", builtins.NewExploreCommand(), []string{"structure", "entry"}},
		{"yaah/scan", builtins.NewScanCommand(), []string{"--security", "--quality"}},
		{"yaah/import", builtins.NewImportCommand(), []string{".planning/", "reverse-engineer"}},
		{"yaah/autonomous", builtins.NewAutonomousCommand(), []string{"autonomous", "dry-run"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := tt.cmd.Content()
			for _, kw := range tt.keywords {
				if !strings.Contains(content, kw) {
					t.Errorf("command %q: Content() does not contain %q", tt.name, kw)
				}
			}
		})
	}
}
