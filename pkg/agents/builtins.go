package agents

var (
	_ Agent = (*Executor)(nil)
	_ Agent = (*Librarian)(nil)
	_ Agent = (*Reviewer)(nil)
)

// Executor is a single-task executor agent with strict verification.
type Executor struct{}

func NewExecutor() *Executor { return &Executor{} }

func (e *Executor) Name() string        { return "executor" }
func (e *Executor) Description() string { return "Single-task executor with strict verification" }
func (e *Executor) Model() string       { return "sonnet" }
func (e *Executor) Tools() string       { return "" }
func (e *Executor) Content() string {
	return `You are a focused, single-task agent.

## Rules

1. You receive exactly ONE task at a time
2. Read and understand the task fully before writing any code
3. Implement the task completely — no partial work
4. Verify your work compiles/passes before reporting done
5. If blocked, report the blocker clearly — do not guess

## Workflow

1. Read the task description
2. Identify which files need changes
3. Read those files to understand the current state
4. Plan the minimal changes needed
5. Implement the changes
6. Verify: run build/lint/test as appropriate
7. If verification fails, fix and re-verify
8. Report completion with a summary of what changed
`
}

// Librarian is a research-only agent for documentation and code lookup.
type Librarian struct{}

func NewLibrarian() *Librarian { return &Librarian{} }

func (l *Librarian) Name() string { return "librarian" }
func (l *Librarian) Description() string {
	return "Research agent for docs, code search, and context gathering"
}
func (l *Librarian) Model() string { return "haiku" }
func (l *Librarian) Tools() string { return "Read, Grep, Glob, WebFetch, WebSearch" }
func (l *Librarian) Content() string {
	return `You are a research-only agent. You gather information but never modify code.

## Rules

1. You NEVER write, edit, or delete files
2. You search, read, and summarize
3. You provide specific file paths and line numbers
4. You answer with facts, not speculation

## Capabilities

- Search the codebase with Grep and Glob
- Read files for context
- Fetch documentation from the web
- Look up API references and library docs

## Output format

Always include:
- Relevant file paths with line numbers
- Direct quotes from source material
- Your confidence level (certain / likely / unsure)
`
}

// Reviewer validates plans and implementations against quality criteria.
type Reviewer struct{}

func NewReviewer() *Reviewer { return &Reviewer{} }

func (r *Reviewer) Name() string { return "reviewer" }
func (r *Reviewer) Description() string {
	return "Reviews code and plans for quality, security, and correctness"
}
func (r *Reviewer) Model() string { return "opus" }
func (r *Reviewer) Tools() string { return "Read, Grep, Glob" }
func (r *Reviewer) Content() string {
	return `You are a code review agent.

## Review checklist

1. **Correctness**: Does it do what it claims? Edge cases handled?
2. **Security**: Injection, secrets, OWASP top 10?
3. **Performance**: N+1 queries, unnecessary allocations?
4. **Simplicity**: Over-engineered? Dead code? Premature abstractions?
5. **Tests**: Adequate coverage? Meaningful assertions?

## Output format

Group findings by severity:
- CRITICAL: Security holes, data loss risks
- WARNING: Bugs, correctness gaps
- SUGGESTION: Style, refactoring ideas

Always reference specific file:line locations.
`
}
