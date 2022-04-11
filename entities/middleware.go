package entities

import "encoding/json"

type Middleware struct {
	Name     string          `json:"name"`
	Settings json.RawMessage `json:"settings"`
}
