//%: {{ if .Geb }}
//%: {{ `
package skeleton //%: ` | replace "skeleton" .SchemaPackage }}

import "gitlab.com/proemergotech/uuid-go"

type AddedEvent struct {
	Group    string    `geb:"group" validate:"required"`
	UUID     uuid.UUID `geb:"uuid" validate:"required,uuid"`
	TestUUID uuid.UUID `geb:"test_uuid" validate:"required,uuid"`
}

type UpdatedEvent struct {
	UUID          uuid.UUID `geb:"uuid" validate:"required,uuid"`
	ChangedFields []string  `json:"changed_fields" geb:"changed_fields"`
	OldVersion    string    `json:"old_version" geb:"old_version"`
	NewVersion    string    `json:"new_version" geb:"new_version"`
}

//%: {{ end }}
