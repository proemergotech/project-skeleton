//%: {{ `
package skeleton //%: ` | replace "skeleton" .SchemaPackage }}

import (
	"gitlab.com/proemergotech/uuid-go"
)

//%: {{ if .Examples }}
type DummyRequest struct {
	DummyParam string    `param:"dummy_param"`
	DummyGroup string    `json:"dummy_group" validate:"required"`
	DummyUUID  uuid.UUID `json:"dummy_uuid"`
} //%: {{ end }}
