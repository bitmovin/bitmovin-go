package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type StreamInputData struct {
	Result StreamInputResult `json:"result"`
}

type StreamInputResult struct {
	FormatName   *string                `json:"formatName"`
	StartTime    *float64               `json:"startTime"`
	Duration     *float64               `json:"duration"`
	Size         *int64                 `json:"size"`
	Bitrate      *int64                 `json:"bitrate"`
	AudioStreams []StreamInputAudio     `json:"audioStreams"`
	VideoStreams []StreamInputVideo     `json:"videoStreams"`
	Tags         map[string]interface{} `json:"tags"`
}

type StreamInputAudio struct {
	ID              *string  `json:"id"`
	Position        *int64   `json:"position"`
	Duration        *float64 `json:"duration"`
	Codec           *string  `json:"codec"`
	SampleRate      *int64   `json:"sampleRate"`
	Bitrate         *int64   `json:"bitrate,string"`
	ChannelFormat   *string  `json:"channelFormat"`
	Language        *string  `json:"language"`
	HearingImpaired *bool    `json:"hearingImpaired"`
}

type StreamInputVideo struct {
	ID       *string  `json:"id"`
	Position *int64   `json:"position"`
	Duration *float64 `json:"duration"`
	Codec    *string  `json:"codec"`
	FPS      *string  `json:"fps"`
	Bitrate  *int64   `json:"bitrate,string"`
	Width    *int64   `json:"width"`
	Height   *int64   `json:"height"`
}

type StreamInputResponse struct {
	RequestID *string                      `json:"requestId"`
	Status    bitmovintypes.ResponseStatus `json:"status"`
	Data      StreamInputData              `json:"data"`
}
