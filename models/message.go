package models

import (
	"github.com/streamco/bitmovin-go/bitmovintypes"
)

type Message struct {
	Date *string                   `json:"date"`
	ID   *string                   `json:"id"`
	Type bitmovintypes.MessageType `json:"type"`
	Text *string                   `json:"text"`
}
