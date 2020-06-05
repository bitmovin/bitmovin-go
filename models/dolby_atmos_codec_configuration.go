package models

import "github.com/streamco/bitmovin-go/bitmovintypes"

type DolbyAtmosLoudnessControl struct {
	MeteringMode         string `json:"meteringMode"`
	DialogueIntelligence string `json:"dialogueIntelligence"`
	SpeechThreshold      int    `json:"speechThreshold"`
}

type DolbyAtmosCodecConfiguration struct {
	ID              string                    `json:"id,omitempty"`
	Name            string                    `json:"name"`
	Description     string                    `json:"description,omitempty"`
	CustomData      map[string]interface{}    `json:"customData,omitempty"`
	Bitrate         int64                     `json:"bitrate"`
	SamplingRate    float64                   `json:"rate,omitempty"`
	LoudnessControl DolbyAtmosLoudnessControl `json:"loudnessControl"`
}

type DolbyAtmosCodecConfigurationResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      struct {
		Result DolbyAtmosCodecConfiguration `json:"result,omitempty"`
	} `json:"data,omitempty"`
}
