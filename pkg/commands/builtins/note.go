package builtins

// NoteCommand captures zero-friction ideas via /yaah:note.
// It is lightweight and does not implement CommandWithAdvanced.
type NoteCommand struct{}

// NewNoteCommand creates a new NoteCommand.
func NewNoteCommand() *NoteCommand { return &NoteCommand{} }

func (c *NoteCommand) Name() string         { return "yaah/note" }
func (c *NoteCommand) Description() string  { return "Capture an idea or observation for later" }
func (c *NoteCommand) ArgumentHint() string { return "<text>" }
func (c *NoteCommand) Model() string        { return "" }
func (c *NoteCommand) AllowedTools() string { return "Read, Write" }

func (c *NoteCommand) Content() string {
	return `# /yaah:note — Zero-Friction Idea Capture

## When to use
When the user runs ` + "`/yaah:note <text>`" + ` to capture a quick idea, observation, or reminder without interrupting current work.

## Steps

### 1. Determine the target file
- Target file: ` + "`.planning/notes/{YYYY-MM-DD}.md`" + ` where ` + "`{YYYY-MM-DD}`" + ` is today's date
- Create ` + "`.planning/notes/`" + ` directory if it does not exist
- Create the daily file if it does not exist (write the header ` + "`# Notes — {YYYY-MM-DD}`" + ` as the first line)

### 2. Append the note
- Read the current file to count existing entries (determines ` + "`{count}`" + `)
- Append a new line in this format: ` + "`- [{HH:MM}] {text}`"  + `
- Use 24-hour time (local clock)
- Write the updated file

### 3. Confirm
- Print: "Noted. ({count} notes today)"
- ` + "`{count}`" + ` is the total number of notes in today's file after appending

## File format
` + "```" + `markdown
# Notes — 2026-04-06
- [09:14] Consider caching the token refresh response to reduce latency
- [11:32] The retry logic in pkg/client may need exponential backoff
- [14:05] Check if the new rate limiter interacts badly with the test suite
` + "```" + `

## Rules
- No git commit — notes are ephemeral scratch space
- Do not create any file other than the daily notes file
- Do not modify any existing source file
- If no argument text is provided, ask the user what to note
`
}
