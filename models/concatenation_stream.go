package models

type ConcatenationStream struct {
	StreamID   string `json:"inputStreamId"`
	IsMain     bool   `json:"isMain"`
	Position   int    `json:"position"`
	AspectMode string `json:"aspectMode,omitempty"`
}
