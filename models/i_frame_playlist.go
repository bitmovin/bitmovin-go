package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type IFramePlaylist struct {
	ID          *string                `json:"id,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	CustomData  map[string]interface{} `json:"customData,omitempty"`
	Filename    *string                `json:"filename,omitempty"`
}

type IFramePlaylistData struct {
	//Success fields
	Result   IFramePlaylist `json:"result,omitempty"`
	Messages []Message      `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type IFramePlaylistResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      IFramePlaylistData           `json:"data,omitempty"`
}
