package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/bitmovintypes"
	"github.com/bitmovin/bitmovin-go/models"
	"github.com/bitmovin/bitmovin-go/services"
)

const (
	MaxRetries int = 500
)

func main() {
	// Creating Bitmovin object
	bitmovin := bitmovin.NewBitmovin("YOUR API KEY", "https://api.bitmovin.com/v1/", 5)

	// Creating the RTMP Input
	rtmpIS := services.NewRTMPInputService(bitmovin)
	rtmpInputListResp, err := rtmpIS.List(0, 10)
	errorHandler(err)

	if len(rtmpInputListResp.Data.Result.Items) < 1 {
		fmt.Println("No RTMP inputs on account!")
		return
	}

	s3OS := services.NewS3OutputService(bitmovin)
	s3Output := &models.S3Output{
		AccessKey:   stringToPtr("YOUR_ACCESS_KEY"),
		SecretKey:   stringToPtr("YOUR_SECRET_KEY"),
		BucketName:  stringToPtr("YOUR_BUCKET_NAME"),
		CloudRegion: bitmovintypes.AWSCloudRegionUSEast1,
	}
	s3OutputResp, err := s3OS.Create(s3Output)
	errorHandler(err)

	encodingS := services.NewEncodingService(bitmovin)
	encoding := &models.Encoding{
		Name:        stringToPtr("example golang live encoding"),
		CloudRegion: bitmovintypes.CloudRegionGoogleUSEast1,
	}
	encodingResp, err := encodingS.Create(encoding)
	errorHandler(err)

	encodingID := *encodingResp.Data.Result.ID

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
	aacConfig := &models.AACCodecConfiguration{
		Name:         stringToPtr("example_audio_codec_configuration"),
		Bitrate:      intToPtr(128000),
		SamplingRate: floatToPtr(48000.0),
	}
	aacResp, err := aacS.Create(aacConfig)
	errorHandler(err)

	videoInputStream := models.InputStream{
		InputID:       rtmpInputListResp.Data.Result.Items[0].ID,
		InputPath:     stringToPtr("live"),
		Position:      intToPtr(0),
		SelectionMode: bitmovintypes.SelectionModeAuto,
	}
	audioInputStream := models.InputStream{
		InputID:       rtmpInputListResp.Data.Result.Items[0].ID,
		InputPath:     stringToPtr("live"),
		Position:      intToPtr(1),
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

	videoStream1080pResp, err := encodingS.AddStream(encodingID, videoStream1080p)
	errorHandler(err)
	videoStream720pResp, err := encodingS.AddStream(encodingID, videoStream720p)
	errorHandler(err)

	ais := []models.InputStream{audioInputStream}
	audioStream := &models.Stream{
		CodecConfigurationID: aacResp.Data.Result.ID,
		InputStreams:         ais,
	}
	aacStreamResp, err := encodingS.AddStream(encodingID, audioStream)
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
	audioMuxingStream := models.StreamItem{
		StreamID: aacStreamResp.Data.Result.ID,
	}

	videoMuxing1080pOutput := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: stringToPtr("golang_live_test/video/1080p"),
		ACL:        acl,
	}
	videoMuxing720pOutput := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: stringToPtr("golang_live_test/video/720p"),
		ACL:        acl,
	}
	audioMuxingOutput := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: stringToPtr("golang_live_test/audio"),
		ACL:        acl,
	}

	videoMuxing1080p := &models.FMP4Muxing{
		SegmentLength:   floatToPtr(4.0),
		SegmentNaming:   stringToPtr("seg_%number%.m4s"),
		InitSegmentName: stringToPtr("init.mp4"),
		Streams:         []models.StreamItem{videoMuxingStream1080p},
		Outputs:         []models.Output{videoMuxing1080pOutput},
	}
	videoMuxing1080pResp, err := encodingS.AddFMP4Muxing(encodingID, videoMuxing1080p)
	errorHandler(err)

	videoMuxing720p := &models.FMP4Muxing{
		SegmentLength:   floatToPtr(4.0),
		SegmentNaming:   stringToPtr("seg_%number%.m4s"),
		InitSegmentName: stringToPtr("init.mp4"),
		Streams:         []models.StreamItem{videoMuxingStream720p},
		Outputs:         []models.Output{videoMuxing720pOutput},
	}
	videoMuxing720pResp, err := encodingS.AddFMP4Muxing(encodingID, videoMuxing720p)
	errorHandler(err)

	audioMuxing := &models.FMP4Muxing{
		SegmentLength:   floatToPtr(4.0),
		SegmentNaming:   stringToPtr("seg_%number%.m4s"),
		InitSegmentName: stringToPtr("init.mp4"),
		Streams:         []models.StreamItem{audioMuxingStream},
		Outputs:         []models.Output{audioMuxingOutput},
	}
	audioMuxingResp, err := encodingS.AddFMP4Muxing(encodingID, audioMuxing)
	errorHandler(err)

	manifestOutput := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: stringToPtr("golang_live_test/manifest"),
		ACL:        acl,
	}
	dashManifest := &models.DashManifest{
		ManifestName: stringToPtr("your_manifest_name.mpd"),
		Outputs:      []models.Output{manifestOutput},
	}
	dashService := services.NewDashManifestService(bitmovin)
	dashManifestResp, err := dashService.Create(dashManifest)
	errorHandler(err)

	period := &models.Period{}
	periodResp, err := dashService.AddPeriod(*dashManifestResp.Data.Result.ID, period)
	errorHandler(err)

	vas := &models.VideoAdaptationSet{}
	vasResp, err := dashService.AddVideoAdaptationSet(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, vas)
	errorHandler(err)

	aas := &models.AudioAdaptationSet{
		Language: stringToPtr("en"),
	}
	aasResp, err := dashService.AddAudioAdaptationSet(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, aas)
	errorHandler(err)

	fmp4Rep1080 := &models.FMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    videoMuxing1080pResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: stringToPtr("../video/1080p"),
	}
	_, err = dashService.AddFMP4Representation(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *vasResp.Data.Result.ID, fmp4Rep1080)
	errorHandler(err)

	fmp4Rep720 := &models.FMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    videoMuxing720pResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: stringToPtr("../video/720p"),
	}
	_, err = dashService.AddFMP4Representation(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *vasResp.Data.Result.ID, fmp4Rep720)
	errorHandler(err)

	fmp4RepAudio := &models.FMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    audioMuxingResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: stringToPtr("../audio"),
	}
	_, err = dashService.AddFMP4Representation(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *aasResp.Data.Result.ID, fmp4RepAudio)
	errorHandler(err)

	liveDashManifest := models.LiveDashManifest{
		ManifestID:     dashManifestResp.Data.Result.ID,
		LiveEdgeOffset: floatToPtr(45.0),
	}

	liveStreamConfig := &models.LiveStreamConfiguration{
		StreamKey:     stringToPtr("YOUR STREAM KEY"),
		DashManifests: []models.LiveDashManifest{liveDashManifest},
	}

	_, err = encodingS.StartLive(encodingID, liveStreamConfig)
	errorHandler(err)

	for numRetries := 0; numRetries < MaxRetries; numRetries++ {
		time.Sleep(10 * time.Second)
		statusResp, err := encodingS.RetrieveLiveStatus(encodingID)
		if err != nil {
			be, ok := err.(models.BitmovinError)
			if ok {
				if be.DataEnvelope.Data.Code != 2023 {
					fmt.Println("Error in starting live encoding")
					fmt.Println(err)
					return
				}
			} else {
				fmt.Println("General Error, exiting.")
				fmt.Println(err)
				return
			}
			fmt.Println("Encoding details not ready yet.")
			continue
		}
		if statusResp != nil {
			if statusResp.Data.Result.EncoderIP == nil {
				fmt.Println("Encoder IP detail empty, encoding failed")
				return
			}
			if statusResp.Data.Result.StreamKey == nil {
				fmt.Println("Stream Key detail empty, encoding failed")
				return
			}
			fmt.Println("---------------")
			fmt.Println("Live Stream set up successfully:")
			fmt.Printf("Encoding ID ... %v \n", encodingID)
			fmt.Printf("Encoder IP .... %v \n", *statusResp.Data.Result.EncoderIP)
			fmt.Printf("Stream Key .... %v \n", *statusResp.Data.Result.StreamKey)
			fmt.Printf("Stream URL: ... rtmp://%v/live \n", *statusResp.Data.Result.EncoderIP)
			fmt.Println("---------------")
			return
		}
	}
	fmt.Println("Maximum number of retries reached.")
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
