//%: {{ if .Geb }}
//%: {{ `
package skeleton //%: ` | replace "skeleton" .SchemaPackage }}

import "gitlab.com/proemergotech/uuid-go"

type AddedEvent struct {
	Group    string    `geb:"group" validate:"required"`
	UUID     uuid.UUID `geb:"uuid" validate:"required,uuid"`
	TestUUID uuid.UUID `geb:"test_uuid" validate:"required,uuid"`
}

//%: {{ end }}
