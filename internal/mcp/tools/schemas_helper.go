package tools

import (
	"encoding/json"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
)

// mustParseSchema unmarshals a JSON schema string into a *jsonschema.Schema and panics on error.
func mustParseSchema(s string) *jsonschema.Schema {
	var sch jsonschema.Schema
	if err := json.Unmarshal([]byte(s), &sch); err != nil {
		panic(fmt.Sprintf("invalid JSON schema: %v\nschema: %s", err, s))
	}
	return &sch
}
