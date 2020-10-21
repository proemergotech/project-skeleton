//%: {{ `
package skeleton //%: ` | replace "skeleton" .SchemaPackage }}

import "gitlab.com/proemergotech/uuid-go"

//%: {{ if .Examples }}

//%: {{ if .Yafuds }}
const (
	YcolBaseEditable = "BASE_EDITABLE"
) //%: {{ end }}

// todo: remove
//  Example struct with redis and yafuds tags
type DummyType struct {
	DummyUUID uuid.UUID `json:"dummy_uuid" redis:"dummy_uuid"`
	BaseEditable
}

type BaseEditable struct {
	Group    string    `json:"group" redis:"group" yafuds:"group"`
	TestUUID uuid.UUID `json:"test_uuid" redis:"test_uuid" yafuds:"test_uuid"`
}

//%: {{ end }}
