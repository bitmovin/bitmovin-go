package models

import "github.com/streamco/bitmovin-go/bitmovintypes"

type AV1CodecConfiguration struct {
	ID                           *string                           `json:"id,omitempty"`
	Name                         *string                           `json:"name,omitempty"`
	Description                  *string                           `json:"description,omitempty"`
	CustomData                   map[string]interface{}            `json:"customData,omitempty"`
	Width                        *int64                            `json:"width,omitempty"`
	Height                       *int64                            `json:"height,omitempty"`
	Bitrate                      *int64                            `json:"bitrate,omitempty"`
	FrameRate                    *float64                          `json:"rate,omitempty"`
	PixelFormat                  bitmovintypes.PixelFormat         `json:"pixelFormat,omitempty"`
	ColorConfig                  ColorConfig                       `json:"colorConfig,omitempty"`
	SampleAspectRatioNumerator   *int64                            `json:"sampleAspectRatioNumerator,omitempty"`
	SampleAspectRatioDenominator *int64                            `json:"sampleAspectRatioDenominator,omitempty"`
	DisplayAspectRatio           *bitmovintypes.AspectRatio        `json:"displayAspectRatio,omitempty"`
	EncodingMode                 bitmovintypes.EncodingMode        `json:"encodingMode,omitempty"`
	PresetConfiguration          bitmovintypes.PresetConfiguration `json:"presetConfiguration,omitempty"`
}

type AV1CodecConfigurationData struct {
	//Success fields
	Result   H265CodecConfiguration `json:"result,omitempty"`
	Messages []Message              `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type AV1CodecConfigurationResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      H265CodecConfigurationData   `json:"data,omitempty"`
}

type AV1CodecConfigurationListResult struct {
	TotalCount *int64                  `json:"totalCount,omitempty"`
	Previous   *string                 `json:"previous,omitempty"`
	Next       *string                 `json:"next,omitempty"`
	Items      []AV1CodecConfiguration `json:"items,omitempty"`
}

type AV1CodecConfigurationListData struct {
	Result AV1CodecConfigurationListResult `json:"result,omitempty"`
}

type AV1CodecConfigurationListResponse struct {
	RequestID *string                       `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus  `json:"status,omitempty"`
	Data      AV1CodecConfigurationListData `json:"data,omitempty"`
}
