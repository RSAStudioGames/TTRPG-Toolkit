package services

import (
	"embed"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

//go:embed schemas/system.json
var systemSchemaFS embed.FS

func validateSystemImportJSON(data []byte) error {
	schemaBytes, err := systemSchemaFS.ReadFile("schemas/system.json")
	if err != nil {
		return fmt.Errorf("load schema: %w", err)
	}
	schemaLoader := gojsonschema.NewBytesLoader(schemaBytes)
	documentLoader := gojsonschema.NewBytesLoader(data)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("schema validation: %w", err)
	}
	if !result.Valid() {
		var msgs []string
		for _, e := range result.Errors() {
			msgs = append(msgs, e.String())
		}
		if len(msgs) > 0 {
			return fmt.Errorf("%s", msgs[0])
		}
		return fmt.Errorf("imported JSON does not match system schema")
	}
	return nil
}
