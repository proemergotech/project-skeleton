//%: {{ `
package skeleton //%: ` | replace "skeleton" .SchemaPackage }}

import "github.com/proemergotech/uuid"

//%: {{ if .Examples }}

// todo: remove
//  Example struct with redis tags
type DummyType struct {
	DummyUUID uuid.UUID `json:"dummy_uuid" redis:"dummy_uuid"`
	BaseEditable
}

type BaseEditable struct {
	Group    string    `json:"group" redis:"group"`
	TestUUID uuid.UUID `json:"test_uuid" redis:"test_uuid"`
}

//%: {{ end }}
