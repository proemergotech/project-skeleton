//%: {{ if .Centrifuge }}
//%: {{ `
package skeleton //%: ` | replace "skeleton" .SchemaPackage }}

import (
	"gitlab.com/proemergotech/uuid-go"
)

//%: {{ if .Examples }}
type CentrifugeData struct {
	DummyUUID uuid.UUID `centrifuge:"dummy_uuid"`
	Group     string    `centrifuge:"group"`
} //%: {{ end }}
//%: {{ end }}
