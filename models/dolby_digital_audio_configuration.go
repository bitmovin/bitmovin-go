package models

import (
	"github.com/streamco/bitmovin-go/bitmovintypes"
)

// DolbyDigitalAudioConfiguration model
type DolbyDigitalAudioConfiguration struct {
	// Id of the resource (required)
	Id *string `json:"id,omitempty"`
	// Name of the resource. Can be freely chosen by the user. (required)
	Name *string `json:"name,omitempty"`
	// Description of the resource. Can be freely chosen by the user.
	Description *string `json:"description,omitempty"`
	// Creation timestamp, returned as UTC expressed in ISO 8601 format: YYYY-MM-DDThh:mm:ssZ
	CreatedAt *string `json:"createdAt,omitempty"`
	// Modified timestamp, returned as UTC expressed in ISO 8601 format: YYYY-MM-DDThh:mm:ssZ
	ModifiedAt *string `json:"modifiedAt,omitempty"`
	// User-specific meta data. This can hold anything.
	CustomData *map[string]interface{} `json:"customData,omitempty"`
	// Target bitrate for the encoded audio in bps (required)
	Bitrate *int64 `json:"bitrate,omitempty"`
	// Audio sampling rate in Hz
	Rate *float64 `json:"rate,omitempty"`
	// BitstreamInfo defines metadata parameters contained in the Dolby Digital audio bitstream
	BitstreamInfo *DolbyDigitalBitstreamInfo `json:"bitstreamInfo,omitempty"`
	// Channel layout of the audio codec configuration.
	ChannelLayout bitmovintypes.DolbyDigitalChannelLayout `json:"channelLayout,omitempty"`
	Downmixing    *DolbyDigitalDownmixing                 `json:"downmixing,omitempty"`
	// It provides a framework for signaling new evolution framework applications, such as Intelligent Loudness, in each Dolby codec.
	EvolutionFrameworkControl bitmovintypes.DolbyDigitalEvolutionFrameworkControl `json:"evolutionFrameworkControl,omitempty"`
	// Settings for loudness control (required)
	LoudnessControl *DolbyDigitalLoudnessControl `json:"loudnessControl,omitempty"`
	Preprocessing   *DolbyDigitalPreprocessing   `json:"preprocessing,omitempty"`
}
