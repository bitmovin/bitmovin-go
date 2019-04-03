package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/streamco/bitmovin-go/bitmovin"
	"github.com/streamco/bitmovin-go/bitmovintypes"
	"github.com/streamco/bitmovin-go/models"
	"github.com/streamco/bitmovin-go/services"
)

func main() {
	// Creating Bitmovin object
	bitmovin := bitmovin.NewBitmovin("YOUR API KEY", "https://api.bitmovin.com/v1/", 5)

	// Creating the HTTP Input
	httpIS := services.NewHTTPInputService(bitmovin)
	httpInput := &models.HTTPInput{
		Host: stringToPtr("YOUR HTTP HOST"),
	}
	httpResp, err := httpIS.Create(httpInput)
	errorHandler(err)

	s3OS := services.NewS3OutputService(bitmovin)
	s3Output := &models.S3Output{
		AccessKey:   stringToPtr("YOUR_ACCESS_KEY"),
		SecretKey:   stringToPtr("YOUR_SECRET_KEY"),
		BucketName:  stringToPtr("YOUR_BUCKET_NAME"),
		CloudRegion: bitmovintypes.AWSCloudRegionEUWest1,
	}
	s3OutputResp, err := s3OS.Create(s3Output)
	errorHandler(err)

	encodingS := services.NewEncodingService(bitmovin)
	encoding := &models.Encoding{
		Name:        stringToPtr("example encoding"),
		CloudRegion: bitmovintypes.CloudRegionGoogleEuropeWest1,
	}
	encodingResp, err := encodingS.Create(encoding)
	errorHandler(err)

	h264S := services.NewH264CodecConfigurationService(bitmovin)
	video1080pConfig := &models.H264CodecConfiguration{
		Name:      stringToPtr("example_video_codec_configuration_1080p"),
		Bitrate:   intToPtr(4800000),
		FrameRate: floatToPtr(25.0),
		Width:     intToPtr(1920),
		Height:    intToPtr(1080),
		Profile:   bitmovintypes.H264ProfileHigh,
	}
	video720Config := &models.H264CodecConfiguration{
		Name:      stringToPtr("example_video_codec_configuration_720p"),
		Bitrate:   intToPtr(2400000),
		FrameRate: floatToPtr(25.0),
		Width:     intToPtr(1280),
		Height:    intToPtr(720),
		Profile:   bitmovintypes.H264ProfileHigh,
	}
	video1080pResp, err := h264S.Create(video1080pConfig)
	errorHandler(err)
	video720Resp, err := h264S.Create(video720Config)
	errorHandler(err)

	aacS := services.NewAACCodecConfigurationService(bitmovin)
	aacConfig128k := &models.AACCodecConfiguration{
		Name:         stringToPtr("example_audio_codec_configuration"),
		Bitrate:      intToPtr(128000),
		SamplingRate: floatToPtr(48000.0),
	}
	aac128kResp, err := aacS.Create(aacConfig128k)
	errorHandler(err)

	aacConfig98k := &models.AACCodecConfiguration{
		Name:         stringToPtr("example_audio_codec_configuration"),
		Bitrate:      intToPtr(98000),
		SamplingRate: floatToPtr(48000.0),
	}
	aac98kResp, err := aacS.Create(aacConfig98k)
	errorHandler(err)

	videoInputStream := models.InputStream{
		InputID:       httpResp.Data.Result.ID,
		InputPath:     stringToPtr("YOUR INPUT FILE PATH AND LOCATION"),
		SelectionMode: bitmovintypes.SelectionModeAuto,
	}
	audioInputStream := models.InputStream{
		InputID:       httpResp.Data.Result.ID,
		InputPath:     stringToPtr("YOUR INPUT FILE PATH AND LOCATION"),
		SelectionMode: bitmovintypes.SelectionModeAuto,
	}

	vis := []models.InputStream{videoInputStream}

	videoStream1080p := &models.Stream{
		CodecConfigurationID: video1080pResp.Data.Result.ID,
		InputStreams:         vis,
	}

	videoStream720p := &models.Stream{
		CodecConfigurationID: video720Resp.Data.Result.ID,
		InputStreams:         vis,
	}

	videoStream1080pResp, err := encodingS.AddStream(*encodingResp.Data.Result.ID, videoStream1080p)
	errorHandler(err)
	videoStream720pResp, err := encodingS.AddStream(*encodingResp.Data.Result.ID, videoStream720p)
	errorHandler(err)

	ais := []models.InputStream{audioInputStream}
	audioStream128k := &models.Stream{
		CodecConfigurationID: aac128kResp.Data.Result.ID,
		InputStreams:         ais,
		Conditions:           models.NewAttributeCondition(bitmovintypes.ConditionAttributeBitrate, ">=", "120000"),
	}
	aacStream128kResp, err := encodingS.AddStream(*encodingResp.Data.Result.ID, audioStream128k)
	errorHandler(err)

	audioStream98k := &models.Stream{
		CodecConfigurationID: aac98kResp.Data.Result.ID,
		InputStreams:         ais,
	}
	aacStream98kResp, err := encodingS.AddStream(*encodingResp.Data.Result.ID, audioStream98k)
	errorHandler(err)

	aclEntry := models.ACLItem{
		Permission: bitmovintypes.ACLPermissionPublicRead,
	}
	acl := []models.ACLItem{aclEntry}

	videoMuxingStream1080p := models.StreamItem{
		StreamID: videoStream1080pResp.Data.Result.ID,
	}
	videoMuxingStream720p := models.StreamItem{
		StreamID: videoStream720pResp.Data.Result.ID,
	}
	audioMuxing128kStream := models.StreamItem{
		StreamID: aacStream128kResp.Data.Result.ID,
	}
	audioMuxing98kStream := models.StreamItem{
		StreamID: aacStream98kResp.Data.Result.ID,
	}

	outputBasePath := "golang"

	videoMuxing1080pOutput := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: stringToPtr(fmt.Sprintf("%s/video/1080p", outputBasePath)),
		ACL:        acl,
	}
	videoMuxing720pOutput := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: stringToPtr(fmt.Sprintf("%s/video/720p", outputBasePath)),
		ACL:        acl,
	}
	audioMuxing128kOutput := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: stringToPtr(fmt.Sprintf("%s/audio/128k", outputBasePath)),
		ACL:        acl,
	}
	audioMuxing98kOutput := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: stringToPtr(fmt.Sprintf("%s/audio/98k", outputBasePath)),
		ACL:        acl,
	}

	videoMuxing1080p := &models.TSMuxing{
		SegmentLength: floatToPtr(4.0),
		SegmentNaming: stringToPtr("seg_%number%.ts"),
		Streams:       []models.StreamItem{videoMuxingStream1080p},
		Outputs:       []models.Output{videoMuxing1080pOutput},
	}
	videoMuxing1080pResp, err := encodingS.AddTSMuxing(*encodingResp.Data.Result.ID, videoMuxing1080p)
	errorHandler(err)

	videoMuxing720p := &models.TSMuxing{
		SegmentLength: floatToPtr(4.0),
		SegmentNaming: stringToPtr("seg_%number%.ts"),
		Streams:       []models.StreamItem{videoMuxingStream720p},
		Outputs:       []models.Output{videoMuxing720pOutput},
	}
	videoMuxing720pResp, err := encodingS.AddTSMuxing(*encodingResp.Data.Result.ID, videoMuxing720p)
	errorHandler(err)

	audioMuxing128k := &models.TSMuxing{
		SegmentLength: floatToPtr(4.0),
		SegmentNaming: stringToPtr("seg_%number%.ts"),
		Streams:       []models.StreamItem{audioMuxing128kStream},
		Outputs:       []models.Output{audioMuxing128kOutput},
	}

	audioMuxing128kResp, err := encodingS.AddTSMuxing(*encodingResp.Data.Result.ID, audioMuxing128k)
	errorHandler(err)

	audioMuxing98k := &models.TSMuxing{
		SegmentLength: floatToPtr(4.0),
		SegmentNaming: stringToPtr("seg_%number%.ts"),
		Streams:       []models.StreamItem{audioMuxing98kStream},
		Outputs:       []models.Output{audioMuxing98kOutput},
	}

	audioMuxing98kResp, err := encodingS.AddTSMuxing(*encodingResp.Data.Result.ID, audioMuxing98k)
	errorHandler(err)

	manifestOutput := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: stringToPtr(outputBasePath),
		ACL:        acl,
	}
	hlsManifest := &models.HLSManifest{
		ManifestName: stringToPtr("your_manifest_name.m3u8"),
		Outputs:      []models.Output{manifestOutput},
	}
	hlsService := services.NewHLSManifestService(bitmovin)
	hlsManifestResp, err := hlsService.Create(hlsManifest)
	errorHandler(err)

	audio128kMediaInfo := &models.MediaInfo{
		Type:               bitmovintypes.MediaTypeAudio,
		URI:                stringToPtr("audio_128k.m3u8"),
		GroupID:            stringToPtr("audio_128"),
		Language:           stringToPtr("en"),
		AssociatedLanguage: stringToPtr("en"),
		Name:               stringToPtr("audio_128"),
		IsDefault:          boolToPtr(false),
		Autoselect:         boolToPtr(false),
		Forced:             boolToPtr(false),
		Characteristics:    []string{"public.accessibility.describes-audio"},
		SegmentPath:        stringToPtr("audio/128k/"),
		EncodingID:         encodingResp.Data.Result.ID,
		StreamID:           aacStream128kResp.Data.Result.ID,
		MuxingID:           audioMuxing128kResp.Data.Result.ID,
	}
	_, err = hlsService.AddMediaInfo(*hlsManifestResp.Data.Result.ID, audio128kMediaInfo)
	errorHandler(err)

	audio98kMediaInfo := &models.MediaInfo{
		Type:               bitmovintypes.MediaTypeAudio,
		URI:                stringToPtr("audio_98.m3u8"),
		GroupID:            stringToPtr("audio_98"),
		Language:           stringToPtr("en"),
		AssociatedLanguage: stringToPtr("en"),
		Name:               stringToPtr("audio_98"),
		IsDefault:          boolToPtr(false),
		Autoselect:         boolToPtr(false),
		Forced:             boolToPtr(false),
		Characteristics:    []string{"public.accessibility.describes-audio"},
		SegmentPath:        stringToPtr("audio/98k/"),
		EncodingID:         encodingResp.Data.Result.ID,
		StreamID:           aacStream98kResp.Data.Result.ID,
		MuxingID:           audioMuxing98kResp.Data.Result.ID,
	}
	_, err = hlsService.AddMediaInfo(*hlsManifestResp.Data.Result.ID, audio98kMediaInfo)
	errorHandler(err)

	video1080pStreamInfo := &models.StreamInfo{
		AudioGroups: &models.HLSAudioGroupConfig{
			DroppingMode: bitmovintypes.HLSVariantStreamDroppingModeStream,
			Groups: []models.HLSAudioGroupDefinition{
				{
					Name:     "audio_98",
					Priority: intToPtr(1),
				},
				{
					Name:     "audio_128",
					Priority: intToPtr(10),
				},
			},
		},
		SegmentPath: stringToPtr("video/1080p/"),
		URI:         stringToPtr("video_hi.m3u8"),
		EncodingID:  encodingResp.Data.Result.ID,
		StreamID:    videoStream1080pResp.Data.Result.ID,
		MuxingID:    videoMuxing1080pResp.Data.Result.ID,
	}
	_, err = hlsService.AddStreamInfo(*hlsManifestResp.Data.Result.ID, video1080pStreamInfo)
	errorHandler(err)

	video720pStreamInfo := &models.StreamInfo{
		AudioGroups: &models.HLSAudioGroupConfig{
			DroppingMode: bitmovintypes.HLSVariantStreamDroppingModeStream,
			Groups: []models.HLSAudioGroupDefinition{
				{
					Name:     "audio_98",
					Priority: intToPtr(1),
				},
				{
					Name:     "audio_128",
					Priority: intToPtr(10),
				},
			},
		},
		SegmentPath: stringToPtr("video/720p/"),
		URI:         stringToPtr("video_lo.m3u8"),
		EncodingID:  encodingResp.Data.Result.ID,
		StreamID:    videoStream720pResp.Data.Result.ID,
		MuxingID:    videoMuxing720pResp.Data.Result.ID,
	}
	_, err = hlsService.AddStreamInfo(*hlsManifestResp.Data.Result.ID, video720pStreamInfo)
	errorHandler(err)

	// Start encoding with manifests for preview and vod
	vodHls := models.VodHlsManifest{
		ManifestID: *hlsManifestResp.Data.Result.ID,
	}

	options := &models.StartOptions{
		VodHlsManifests: []models.VodHlsManifest{vodHls},
	}

	_, err = encodingS.StartWithOptions(*encodingResp.Data.Result.ID, options)
	errorHandler(err)

	var status string
	status = ""
	for status != "FINISHED" {
		time.Sleep(10 * time.Second)
		statusResp, err := encodingS.RetrieveStatus(*encodingResp.Data.Result.ID)
		if err != nil {
			fmt.Println("error in Encoding Status")
			fmt.Println(err)
			return
		}
		// Polling and Printing out the response
		fmt.Printf("%s\n", getAsJsonString(*statusResp))
		status = *statusResp.Data.Result.Status
		if status == "ERROR" {
			fmt.Println("error in Encoding Status")
			fmt.Printf("%s\n", getAsJsonString(*statusResp))
			return
		}
	}
}

func errorHandler(err error) {
	if err != nil {
		switch err.(type) {
		case models.BitmovinError:
			fmt.Println("Bitmovin Error")
		default:
			fmt.Println("General Error")
		}
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func stringToPtr(s string) *string {
	return &s
}

func intToPtr(i int64) *int64 {
	return &i
}

func boolToPtr(b bool) *bool {
	return &b
}

func floatToPtr(f float64) *float64 {
	return &f
}

func getAsJsonString(v interface{}) string {
	j, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(j)
}
