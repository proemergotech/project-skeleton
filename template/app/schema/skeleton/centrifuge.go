//%: {{ if .Centrifuge }}
//%: {{ `
package skeleton //%: ` | replace "skeleton" .SchemaPackage }}

import (
	"gitlab.com/proemergotech/uuid-go"
)

//%: {{ if .Examples }}
type CentrifugeData struct {
	Group string    `centrifuge:"group"`
	UUID  uuid.UUID `centrifuge:"uuid"`
} //%: {{ end }}
//%: {{ end }}
