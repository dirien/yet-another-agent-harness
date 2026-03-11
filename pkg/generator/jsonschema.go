package generator

import (
	"encoding/json"
	"fmt"

	"github.com/dirien/yet-another-agent-harness/pkg/schema"
	"github.com/invopop/jsonschema"
)

// GenerateJSONSchema produces a JSON Schema from the HarnessConfig struct.
func GenerateJSONSchema() ([]byte, error) {
	r := new(jsonschema.Reflector)
	r.ExpandedStruct = true

	s := r.Reflect(&schema.HarnessConfig{})
	s.Title = "Claude Code Harness Configuration"
	s.Description = "Schema for yet-another-agent-harness config files"
	s.ID = "https://github.com/dirien/yet-another-agent-harness/yaah.schema.json"

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal schema: %w", err)
	}
	return data, nil
}
