package centrifuge

import "encoding/json"

type PublishRequest struct {
	Channel string          `json:"channel" validate:"required"`
	Data    json.RawMessage `json:"data"`
}
