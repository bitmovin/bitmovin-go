package models

import (
	"github.com/bitmovin/bitmovin-go/bitmovintypes"
)

type Encoding struct {
	ID             *string                      `json:"id,omitempty"`
	Name           *string                      `json:"name,omitempty"`
	Description    *string                      `json:"description,omitempty"`
	CustomData     map[string]interface{}       `json:"customData,omitempty"`
	EncoderVersion bitmovintypes.EncoderVersion `json:"encoderVersion,omitempty"`
	CloudRegion    bitmovintypes.CloudRegion    `json:"cloudRegion,omitempty"`
}

type EncodingData struct {
	//Success fields
	Result   Encoding  `json:"result,omitempty"`
	Messages []Message `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type EncodingResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      EncodingData                 `json:"data,omitempty"`
}

type EncodingListResult struct {
	TotalCount *int64     `json:"totalCount,omitempty"`
	Previous   *string    `json:"previous,omitempty"`
	Next       *string    `json:"next,omitempty"`
	Items      []Encoding `json:"items,omitempty"`
}

type EncodingListData struct {
	Result EncodingListResult `json:"result,omitempty"`
}

type EncodingListResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      EncodingListData             `json:"data,omitempty"`
}

type InputStream struct {
	InputID       *string                     `json:"inputId,omitempty"`
	InputPath     *string                     `json:"inputPath,omitempty"`
	SelectionMode bitmovintypes.SelectionMode `json:"selectionMode,omitempty"`
	Position      *int64                      `json:"position,omitempty"`
}

type ACLItem struct {
	Scope      *string                     `json:"scope,omitempty"`
	Permission bitmovintypes.ACLPermission `json:"permission,omitempty"`
}

type Output struct {
	OutputID   *string   `json:"outputId,omitempty"`
	OutputPath *string   `json:"outputPath,omitempty"`
	ACL        []ACLItem `json:"acl,omitempty"`
}

type Stream struct {
	ID                   *string                `json:"id,omitempty"`
	Name                 *string                `json:"name,omitempty"`
	Description          *string                `json:"description,omitempty"`
	CustomData           map[string]interface{} `json:"customData,omitempty"`
	InputStreams         []InputStream          `json:"inputStreams,omitempty"`
	CodecConfigurationID *string                `json:"codecConfigId,omitempty"`
	Outputs              []Output               `json:"outputs,omitempty"`
}

type StreamData struct {
	//Success fields
	Result   Stream    `json:"result,omitempty"`
	Messages []Message `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type StreamResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      StreamData                   `json:"data,omitempty"`
}

type StreamListResult struct {
	TotalCount *int64   `json:"totalCount,omitempty"`
	Previous   *string  `json:"previous,omitempty"`
	Next       *string  `json:"next,omitempty"`
	Items      []Stream `json:"items,omitempty"`
}

type StreamListResponse struct {
	Result StreamListResult `json:"result,omitempty"`
}

type StreamItem struct {
	StreamID *string `json:"streamId,omitempty"`
}

type FMP4Muxing struct {
	ID              *string                `json:"id,omitempty"`
	Name            *string                `json:"name,omitempty"`
	Description     *string                `json:"description,omitempty"`
	CustomData      map[string]interface{} `json:"customData,omitempty"`
	Streams         []StreamItem           `json:"streams,omitempty"`
	Outputs         []Output               `json:"outputs,omitempty"`
	SegmentLength   *float64               `json:"segmentLength,omitempty"`
	SegmentNaming   *string                `json:"segmentNaming,omitempty"`
	InitSegmentName *string                `json:"initSegmentName,omitempty"`
}

type FMP4MuxingData struct {
	//Success fields
	Result   FMP4Muxing `json:"result,omitempty"`
	Messages []Message  `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type FMP4MuxingResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      FMP4MuxingData               `json:"data,omitempty"`
}

type FMP4MuxingListResult struct {
	TotalCount *int64       `json:"totalCount,omitempty"`
	Previous   *string      `json:"previous,omitempty"`
	Next       *string      `json:"next,omitempty"`
	Items      []FMP4Muxing `json:"items,omitempty"`
}

type FMP4MuxingListData struct {
	Result FMP4MuxingListResult `json:"result,omitempty"`
}

type FMP4MuxingListResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      FMP4MuxingListData           `json:"data,omitempty"`
}

type TSMuxing struct {
	ID            *string                `json:"id,omitempty"`
	Name          *string                `json:"name,omitempty"`
	Description   *string                `json:"description,omitempty"`
	CustomData    map[string]interface{} `json:"customData,omitempty"`
	Streams       []StreamItem           `json:"streams,omitempty"`
	Outputs       []Output               `json:"outputs,omitempty"`
	SegmentLength *float64               `json:"segmentLength,omitempty"`
	SegmentNaming *string                `json:"segmentNaming,omitempty"`
}

type TSMuxingData struct {
	//Success fields
	Result   TSMuxing  `json:"result,omitempty"`
	Messages []Message `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type TSMuxingResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      TSMuxingData                 `json:"data,omitempty"`
}

type TSMuxingListResult struct {
	TotalCount *int64     `json:"totalCount,omitempty"`
	Previous   *string    `json:"previous,omitempty"`
	Next       *string    `json:"next,omitempty"`
	Items      []TSMuxing `json:"items,omitempty"`
}

type TSMuxingListData struct {
	Result TSMuxingListResult `json:"result,omitempty"`
}

type TSMuxingListResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      TSMuxingListData             `json:"data,omitempty"`
}

type MP4Muxing struct {
	ID               *string                `json:"id,omitempty"`
	Name             *string                `json:"name,omitempty"`
	Description      *string                `json:"description,omitempty"`
	CustomData       map[string]interface{} `json:"customData,omitempty"`
	Streams          []StreamItem           `json:"streams,omitempty"`
	Outputs          []Output               `json:"outputs,omitempty"`
	Filename         *string                `json:"filename,omitempty"`
	FragmentDuration *int64                 `json:"fragmentDuration,omitempty"`
}

type MP4MuxingData struct {
	//Success fields
	Result   MP4Muxing `json:"result,omitempty"`
	Messages []Message `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type MP4MuxingResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      MP4MuxingData                `json:"data,omitempty"`
}

type MP4MuxingListResult struct {
	TotalCount *int64      `json:"totalCount,omitempty"`
	Previous   *string     `json:"previous,omitempty"`
	Next       *string     `json:"next,omitempty"`
	Items      []MP4Muxing `json:"items,omitempty"`
}

type MP4MuxingListData struct {
	Result MP4MuxingListResult `json:"result,omitempty"`
}

type MP4MuxingListResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      MP4MuxingListData            `json:"data,omitempty"`
}

type StartResult struct {
	ID *string `json:"id,omitempty"`
}

type StartData struct {
	//Success fields
	Result   StartResult `json:"result,omitempty"`
	Messages []Message   `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type StartStopResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      StartData                    `json:"data,omitempty"`
}

type Subtask struct {
	Status   *string  `json:"status,omitempty"`
	Name     *string  `json:"name,omitempty"`
	ETA      *float64 `json:"eta,omitempty"`
	Progress *float64 `json:"progress,omitempty"`
}

type StatusResult struct {
	Status   *string   `json:"status,omitempty"`
	ETA      *float64  `json:"eta,omitempty"`
	Progress *float64  `json:"progress,omitempty"`
	Messages []Message `json:"messages,omitempty"`
	Subtasks []Subtask `json:"subtasks,omitempty"`
}

type StatusData struct {
	//Success fields
	Result   StatusResult `json:"result,omitempty"`
	Messages []Message    `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type StatusResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      StatusData                   `json:"data,omitempty"`
}

type LiveStreamConfiguration struct {
	StreamKey     *string            `json:"streamKey,omitempty"`
	HLSManifests  []LiveHLSManifest  `json:"hlsManifests,omitempty"`
	DashManifests []LiveDashManifest `json:"dashManifests,omitempty"`
}

type LiveStatusResult struct {
	StreamKey *string `json:"streamKey,omitempty"`
	EncoderIP *string `json:"encoderIp,omitempty"`
}

type LiveStatusData struct {
	//Success fields
	Result   LiveStatusResult `json:"result,omitempty"`
	Messages []Message        `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type LiveStatusResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      LiveStatusData               `json:"data,omitempty"`
}
