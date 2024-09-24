package model

import "encoding/json"

type Result struct {
	Result json.RawMessage `json:"result"`
}
