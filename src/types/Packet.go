package types

import "encoding/json"

type Packet struct {
	Opcode Opcode          `json:"op,omitempty"`
	Event  string          `json:"event,omitempty"`
	Seq    int             `json:"seq,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
}
