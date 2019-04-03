package models

import "github.com/streamco/bitmovin-go/bitmovintypes"

type CustomTag struct {
	ID           *string                    `json:"id,omitempty"`
	Name         *string                    `json:"name,omitempty"`
	Description  *string                    `json:"description,omitempty"`
	CustomData   map[string]interface{}     `json:"customData,omitempty"`
	PositionMode bitmovintypes.PositionMode `json:"positionMode,omitempty"`
	KeyframeID   *string                    `json:"keyframeId,omitempty"`
	Time         *float64                   `json:"time,omitempty"`
	Segment      *int64                     `json:"segment,omitempty"`
	Data         *string                    `json:"data,omitempty"`
}

type CustomTagData struct {
	//Success fields
	Result   CustomTag `json:"result,omitempty"`
	Messages []Message `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type CustomTagResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      CustomTagData                `json:"data,omitempty"`
}
