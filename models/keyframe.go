package models

import "github.com/streamco/bitmovin-go/bitmovintypes"

type Keyframe struct {
	ID          *string                `json:"id,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	CustomData  map[string]interface{} `json:"customData,omitempty"`
	Time        *float64               `json:"time,omitempty"`
	SegmentCut  *bool                  `json:"segmentCut,omitempty"`
}

type KeyframeData struct {
	//Success fields
	Result   Keyframe  `json:"result,omitempty"`
	Messages []Message `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type KeyframeResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      KeyframeData                 `json:"data,omitempty"`
}

type KeyframeListResponse struct {
	Data KeyframeListData `json:"data,omitempty"`
}

type KeyframeListData struct {
	Result KeyframeListResult `json:"result,omitempty"`
}

type KeyframeListResult struct {
	TotalCount *int64     `json:"totalCount,omitempty"`
	Previous   *string    `json:"previous,omitempty"`
	Next       *string    `json:"next,omitempty"`
	Items      []Keyframe `json:"items,omitempty"`
}
