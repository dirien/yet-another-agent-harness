package builtins

// TodoCommand manages lightweight task capture via /yaah:todo.
// It is lightweight and does not implement CommandWithAdvanced.
type TodoCommand struct{}

// NewTodoCommand creates a new TodoCommand.
func NewTodoCommand() *TodoCommand { return &TodoCommand{} }

func (c *TodoCommand) Name() string        { return "yaah/todo" }
func (c *TodoCommand) Description() string { return "Capture, list, or complete quick todo items" }
func (c *TodoCommand) ArgumentHint() string {
	return "[add <text> | list | done <number>]"
}
func (c *TodoCommand) Model() string        { return "" }
func (c *TodoCommand) AllowedTools() string { return "Read, Write, Glob" }

func (c *TodoCommand) Content() string {
	return `# /yaah:todo — Quick Todo Management

## When to use
When the user runs ` + "`/yaah:todo`" + ` to manage lightweight task items outside the full planning workflow.

## Argument forms
- ` + "`/yaah:todo add <text>`" + ` — append a new todo item
- ` + "`/yaah:todo list`" + ` — show all items with pending/completed status
- ` + "`/yaah:todo done <number>`" + ` — mark item ` + "`#<number>`" + ` as completed
- ` + "`/yaah:todo`" + ` (no argument) — same as ` + "`list`"  + `

## File location
All todos are stored in ` + "`.planning/TODOS.md`" + `. Create the file and ` + "`.planning/`" + ` directory if they do not exist.

## Format
` + "```" + `markdown
# TODOs
- [ ] #1 (2026-04-06) Add rate limiting to API endpoints
- [x] #2 (2026-04-06) Fix typo in README — done 2026-04-06
- [ ] #3 (2026-04-06) Consider adding WebSocket support
` + "```" + `

Each line follows this structure:
- ` + "`- [ ]`" + ` for pending, ` + "`- [x]`" + ` for completed
- ` + "`#{sequential-number}`" + ` — number assigned in order of creation, never reused
- ` + "`({YYYY-MM-DD})`" + ` — date the item was added
- Text of the todo
- For completed items, append ` + "`— done {YYYY-MM-DD}`" + ` at the end of the line

## Steps by subcommand

### add <text>
1. Read ` + "`.planning/TODOS.md`" + ` to determine the next sequential number
2. Append a new pending item with today's date and the provided text
3. Write the updated file
4. Print: "Added #N: {text}"

### list (or no argument)
1. Read ` + "`.planning/TODOS.md`" + `
2. Display all items grouped: pending first, then completed
3. Show totals: "N pending, M completed"

### done <number>
1. Read ` + "`.planning/TODOS.md`" + `
2. Find the item with ` + "`#{number}`" + `
3. Change ` + "`[ ]`" + ` to ` + "`[x]`" + ` and append ` + "`— done {YYYY-MM-DD}`" + `
4. Write the updated file
5. Print: "Marked #N done."

## Rules
- Sequential numbers are permanent — never renumber existing items
- Create ` + "`.planning/TODOS.md`" + ` if it does not exist (write the header ` + "`# TODOs`" + ` and an empty body)
- No git commit — todo updates are intentionally lightweight
- If ` + "`done <number>`" + ` targets a non-existent or already-completed item, report the issue and stop
`
}
