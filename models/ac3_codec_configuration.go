package models

import "github.com/streamco/bitmovin-go/bitmovintypes"

type AC3CodecConfiguration struct {
	ID              string                      `json:"id,omitempty"`
	Name            string                      `json:"name"`
	Description     string                      `json:"description,omitempty"`
	CustomData      map[string]interface{}      `json:"customData,omitempty"`
	Bitrate         int64                       `json:"bitrate"`
	SamplingRate    float64                     `json:"rate,omitempty"`
	ChannelLayout   bitmovintypes.ChannelLayout `json:"channelLayout,omitempty"`
	CutoffFrequency *int64                      `json:"cutoffFrequency,omitempty"`
}

type AC3CodecConfigurationResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      struct {
		Result AC3CodecConfiguration `json:"result,omitempty"`
	} `json:"data,omitempty"`
}
