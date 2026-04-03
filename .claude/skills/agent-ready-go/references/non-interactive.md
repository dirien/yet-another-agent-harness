<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Non-Interactive Design for Agent-Ready Go

Agents execute commands programmatically. Interactive prompts that block waiting
for user input will deadlock an agent. The goal is not to eliminate all
interactivity — it is to ensure every blocking prompt has a machine-usable
bypass.

## The --yes Flag Pattern

Any command that asks for confirmation must support `--yes` / `-y`:

```go
var yes bool
cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation prompts")

func runDelete(cmd *cobra.Command, args []string) error {
    if !yes {
        confirmed, err := promptConfirm(fmt.Sprintf("Delete %d records?", count))
        if err != nil {
            return fmt.Errorf("prompt failed: %w", err)
        }
        if !confirmed {
            return nil
        }
    }
    return store.DeleteAll(cmd.Context())
}
```

Agents invoke: `myapp delete --yes`

## Environment Variable Overrides

For flags agents always want set, support env var fallbacks so they can be
configured once in the environment:

```go
// Accept both flag and env var (flag wins if both set)
apiKey := os.Getenv("MYAPP_API_KEY")
cmd.Flags().StringVar(&apiKey, "api-key", apiKey, "API key (or MYAPP_API_KEY)")
```

Standard env var conventions:

- Prefix with app name: `MYAPP_`
- All caps with underscores: `MYAPP_LOG_LEVEL`, `MYAPP_DRY_RUN`
- Booleans: `MYAPP_YES=1` or `MYAPP_YES=true`

## Avoiding Pure TUI Libraries

Some TUI libraries (bubbletea interactive menus, survey, promptui) render
full-screen UIs that agents cannot navigate.

**Use instead:**

- `--format json|text|csv` flag for format selection
- Positional arguments for required inputs
- Config files for complex multi-field inputs
- `--yes` for confirmations
- `--dry-run` to preview without executing

If you must include a TUI for human users, always provide a non-TUI path:

```go
// TUI path for humans
if term.IsTerminal(int(os.Stdin.Fd())) && !yes {
    return runInteractiveWizard()
}
// Non-TUI path for agents
return runWithFlags(flags)
```

## Dry Run Support

Agents benefit from being able to preview operations without side effects:

```go
var dryRun bool
cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview changes without applying")

if dryRun {
    fmt.Fprintf(os.Stderr, "[dry-run] would delete %d records\n", count)
    return nil
}
// proceed
```

## stdin Detection

Never block reading from stdin unless it is explicitly piped:

```go
stat, _ := os.Stdin.Stat()
if (stat.Mode() & os.ModeCharDevice) != 0 {
    // stdin is a terminal — do not block
    return fmt.Errorf("no input provided; pipe data or use --file flag")
}
data, _ := io.ReadAll(os.Stdin)
```

## CI/CD Environment Detection

CI environments set well-known environment variables. Use them to auto-enable
non-interactive mode so agents and pipelines never encounter blocking prompts:

```go
func isCI() bool {
    return os.Getenv("CI") != "" ||
        os.Getenv("GITHUB_ACTIONS") != "" ||
        os.Getenv("GITLAB_CI") != "" ||
        os.Getenv("JENKINS_URL") != ""
}

// In command flags:
var yes bool
cmd.Flags().BoolVarP(&yes, "yes", "y", isCI(), "skip confirmation prompts")
```

## NO\_COLOR and --no-color

Respect the `NO_COLOR` standard (<https://no-color.org/>) and the `dumb` terminal
type. Agents and log aggregators do not render ANSI escape codes — they appear
as noise in structured output:

```go
func colorEnabled() bool {
    if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
        return false
    }
    return term.IsTerminal(int(os.Stdout.Fd()))
}
```

## --quiet Flag

Suppress informational progress messages when running in scripts or pipelines.
Always write progress to stderr so it can be redirected independently:

```go
var quiet bool
cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "suppress informational output")

if !quiet {
    fmt.Fprintln(os.Stderr, "Processing 42 records...")
}
```

## --timeout Flag

Long-running commands must support a timeout so agents can bound execution time
and fail fast rather than waiting indefinitely:

```go
var timeout time.Duration
cmd.Flags().DurationVar(&timeout, "timeout", 30*time.Second, "maximum execution time")

ctx, cancel := context.WithTimeout(cmd.Context(), timeout)
defer cancel()
```

## Summary Checklist

- [ ] Destructive commands have `--yes`/`-y`
- [ ] CI environments auto-enable `--yes` (via `CI`/`GITHUB_ACTIONS` env vars)
- [ ] Sensitive config accepts env vars with `APPNAME_` prefix
- [ ] No interactive menus without a `--flag` alternative
- [ ] `--dry-run` on all state-changing commands
- [ ] stdin reads are guarded against blocking on a terminal
- [ ] `NO_COLOR` respected; color codes suppressed in non-TTY output
- [ ] `--quiet`/`-q` suppresses informational stderr
- [ ] Long-running commands support `--timeout`
