package builtins

// SettingsCommand configures workflow behavior via /yaah:settings.
// It is lightweight and does not implement CommandWithAdvanced.
type SettingsCommand struct{}

// NewSettingsCommand creates a new SettingsCommand.
func NewSettingsCommand() *SettingsCommand { return &SettingsCommand{} }

func (c *SettingsCommand) Name() string        { return "yaah/settings" }
func (c *SettingsCommand) Description() string { return "View or update workflow configuration" }
func (c *SettingsCommand) ArgumentHint() string { return "[key] [value]" }
func (c *SettingsCommand) Model() string        { return "" }
func (c *SettingsCommand) AllowedTools() string { return "Read, Write" }

func (c *SettingsCommand) Content() string {
	return `# /yaah:settings — View or Update Workflow Configuration

## When to use
When the user runs ` + "`/yaah:settings`" + ` to view current settings, or ` + "`/yaah:settings <key> <value>`" + ` to update a setting.

## Steps

### 1. Parse arguments
- If no arguments: display current settings (read-only mode)
- If one argument (key only): display the current value of that key
- If two arguments (key + value): update that setting in config.json

### 2. Read config.json
- Read ` + "`.planning/config.json`" + `
- If it does not exist, report "No config found — run ` + "`/yaah:init`" + ` first" and stop

### 3. If displaying settings
Print the current configuration in a readable table format:

| Setting | Value | Description |
|---------|-------|-------------|
| mode | {value} | interactive or autonomous |
| granularity | {value} | coarse, standard, or fine |
| model_profile | {value} | quality, balanced, or budget |
| research | {value} | enable research phase (true/false) |
| plan_check | {value} | require plan review before execute (true/false) |
| discuss | {value} | enable discuss phase (true/false) |
| auto_advance | {value} | automatically advance phases (true/false) |

### 4. If updating a setting
Validate the key and value:

**Valid keys and accepted values:**
- ` + "`mode`" + `: ` + "`interactive`" + ` | ` + "`autonomous`" + `
- ` + "`granularity`" + `: ` + "`coarse`" + ` | ` + "`standard`" + ` | ` + "`fine`" + `
- ` + "`model_profile`" + `: ` + "`quality`" + ` | ` + "`balanced`" + ` | ` + "`budget`" + `
- ` + "`research`" + `: ` + "`true`" + ` | ` + "`false`" + `
- ` + "`plan_check`" + `: ` + "`true`" + ` | ` + "`false`" + `
- ` + "`discuss`" + `: ` + "`true`" + ` | ` + "`false`" + `
- ` + "`auto_advance`" + `: ` + "`true`" + ` | ` + "`false`" + `

If the key is unknown or value is invalid, print the valid options and stop — do NOT write.

Update the field in config.json and write it back.

### 5. Print current settings after any change
Always display the full settings table after a successful update so the user can confirm.

## Rules
- NEVER write config.json if the key or value is invalid
- NEVER create config.json — only ` + "`/yaah:init`" + ` should create it
- Display settings in a human-readable table, not raw JSON
`
}
