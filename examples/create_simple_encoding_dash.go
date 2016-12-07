package main

import (
	"fmt"
	"time"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/bitmovintypes"
	"github.com/bitmovin/bitmovin-go/models"
	"github.com/bitmovin/bitmovin-go/services"
)

func CreateDash() {
	// Creating Bitmovin object
	bitmovin := bitmovin.NewBitmovin("YOUR API KEY", "https://api.bitmovin.com/v1/", 5)

	// Creating the HTTP Input
	httpIS := services.NewHTTPInputService(bitmovin)
	httpInput := &models.HTTPInput{
		Host: StringToPtr("YOUR HTTP HOST"),
	}
	httpResp, err := httpIS.Create(httpInput)
	ErrorHandler(httpResp.Status, err)

	s3OS := services.NewS3OutputService(bitmovin)
	s3Output := &models.S3Output{
		AccessKey:   StringToPtr("YOUR_ACCESS_KEY"),
		SecretKey:   StringToPtr("YOUR_SECRET_KEY"),
		BucketName:  StringToPtr("YOUR_BUCKET_NAME"),
		CloudRegion: bitmovintypes.AWSCloudRegionEUWest1,
	}
	s3OutputResp, err := s3OS.Create(s3Output)
	ErrorHandler(s3OutputResp.Status, err)

	encodingS := services.NewEncodingService(bitmovin)
	encoding := &models.Encoding{
		Name:        StringToPtr("example encoding"),
		CloudRegion: bitmovintypes.CloudRegionGoogleEuropeWest1,
	}
	encodingResp, err := encodingS.Create(encoding)
	ErrorHandler(encodingResp.Status, err)

	h264S := services.NewH264CodecConfigurationService(bitmovin)
	video1080pConfig := &models.H264CodecConfiguration{
		Name:      StringToPtr("example_video_codec_configuration_1080p"),
		Bitrate:   IntToPtr(4800000),
		FrameRate: FloatToPtr(25.0),
		Width:     IntToPtr(1920),
		Height:    IntToPtr(1080),
		Profile:   bitmovintypes.H264ProfileHigh,
	}
	video720Config := &models.H264CodecConfiguration{
		Name:      StringToPtr("example_video_codec_configuration_720p"),
		Bitrate:   IntToPtr(2400000),
		FrameRate: FloatToPtr(25.0),
		Width:     IntToPtr(1280),
		Height:    IntToPtr(720),
		Profile:   bitmovintypes.H264ProfileHigh,
	}
	video1080pResp, err := h264S.Create(video1080pConfig)
	ErrorHandler(video1080pResp.Status, err)
	video720Resp, err := h264S.Create(video720Config)
	ErrorHandler(video720Resp.Status, err)

	aacS := services.NewAACCodecConfigurationService(bitmovin)
	aacConfig := &models.AACCodecConfiguration{
		Name:         StringToPtr("example_audio_codec_configuration"),
		Bitrate:      IntToPtr(128000),
		SamplingRate: FloatToPtr(48000.0),
	}
	aacResp, err := aacS.Create(aacConfig)
	ErrorHandler(aacResp.Status, err)

	videoInputStream := models.InputStream{
		InputID:       httpResp.Data.Result.ID,
		InputPath:     StringToPtr("YOUR INPUT FILE PATH AND LOCATION"),
		SelectionMode: bitmovintypes.SelectionModeAuto,
	}
	audioInputStream := models.InputStream{
		InputID:       httpResp.Data.Result.ID,
		InputPath:     StringToPtr("YOUR INPUT FILE PATH AND LOCATION"),
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
	ErrorHandler(videoStream1080pResp.Status, err)
	videoStream720pResp, err := encodingS.AddStream(*encodingResp.Data.Result.ID, videoStream720p)
	ErrorHandler(videoStream720pResp.Status, err)

	ais := []models.InputStream{audioInputStream}
	audioStream := &models.Stream{
		CodecConfigurationID: aacResp.Data.Result.ID,
		InputStreams:         ais,
	}
	aacStreamResp, err := encodingS.AddStream(*encodingResp.Data.Result.ID, audioStream)
	ErrorHandler(aacStreamResp.Status, err)

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
		OutputPath: StringToPtr("golang_test/video/1080p"),
		ACL:        acl,
	}
	videoMuxing720pOutput := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: StringToPtr("golang_test/video/720p"),
		ACL:        acl,
	}
	audioMuxingOutput := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: StringToPtr("golang_test/audio"),
		ACL:        acl,
	}

	videoMuxing1080p := &models.FMP4Muxing{
		SegmentLength:   FloatToPtr(4.0),
		SegmentNaming:   StringToPtr("seg_%number%.m4s"),
		InitSegmentName: StringToPtr("init.mp4"),
		Streams:         []models.StreamItem{videoMuxingStream1080p},
		Outputs:         []models.Output{videoMuxing1080pOutput},
	}
	videoMuxing1080pResp, err := encodingS.AddFMP4Muxing(*encodingResp.Data.Result.ID, videoMuxing1080p)
	ErrorHandler(videoMuxing1080pResp.Status, err)

	videoMuxing720p := &models.FMP4Muxing{
		SegmentLength:   FloatToPtr(4.0),
		SegmentNaming:   StringToPtr("seg_%number%.m4s"),
		InitSegmentName: StringToPtr("init.mp4"),
		Streams:         []models.StreamItem{videoMuxingStream720p},
		Outputs:         []models.Output{videoMuxing720pOutput},
	}
	videoMuxing720pResp, err := encodingS.AddFMP4Muxing(*encodingResp.Data.Result.ID, videoMuxing720p)
	ErrorHandler(videoMuxing720pResp.Status, err)

	audioMuxing := &models.FMP4Muxing{
		SegmentLength:   FloatToPtr(4.0),
		SegmentNaming:   StringToPtr("seg_%number%.m4s"),
		InitSegmentName: StringToPtr("init.mp4"),
		Streams:         []models.StreamItem{audioMuxingStream},
		Outputs:         []models.Output{audioMuxingOutput},
	}
	audioMuxingResp, err := encodingS.AddFMP4Muxing(*encodingResp.Data.Result.ID, audioMuxing)
	ErrorHandler(audioMuxingResp.Status, err)

	startResp, err := encodingS.Start(*encodingResp.Data.Result.ID)
	ErrorHandler(startResp.Status, err)

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
		fmt.Printf("%+v\n", statusResp)
		status = *statusResp.Data.Result.Status
		if status == "ERROR" {
			fmt.Println("error in Encoding Status")
			fmt.Printf("%+v\n", statusResp)
			return
		}
	}

	manifestOutput := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: StringToPtr("golang_test/manifest"),
		ACL:        acl,
	}
	dashManifest := &models.DashManifest{
		ManifestName: StringToPtr("your_manifest_name.mpd"),
		Outputs:      []models.Output{manifestOutput},
	}
	dashService := services.NewDashManifestService(bitmovin)
	dashManifestResp, err := dashService.Create(dashManifest)
	ErrorHandler(dashManifestResp.Status, err)

	period := &models.Period{}
	periodResp, err := dashService.AddPeriod(*dashManifestResp.Data.Result.ID, period)
	ErrorHandler(periodResp.Status, err)

	vas := &models.VideoAdaptationSet{}
	vasResp, err := dashService.AddVideoAdaptationSet(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, vas)
	ErrorHandler(vasResp.Status, err)

	aas := &models.AudioAdaptationSet{
		Language: StringToPtr("en"),
	}
	aasResp, err := dashService.AddAudioAdaptationSet(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, aas)
	ErrorHandler(aasResp.Status, err)

	fmp4Rep1080 := &models.FMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    videoMuxing1080pResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: StringToPtr("../video/1080p"),
	}
	fmp4Rep1080Resp, err := dashService.AddFMP4Representation(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *vasResp.Data.Result.ID, fmp4Rep1080)
	ErrorHandler(fmp4Rep1080Resp.Status, err)

	fmp4Rep720 := &models.FMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    videoMuxing720pResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: StringToPtr("../video/720p"),
	}
	fmp4Rep720Resp, err := dashService.AddFMP4Representation(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *vasResp.Data.Result.ID, fmp4Rep720)
	ErrorHandler(fmp4Rep720Resp.Status, err)

	fmp4RepAudio := &models.FMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    audioMuxingResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: StringToPtr("../audio"),
	}
	fmp4RepAudioResp, err := dashService.AddFMP4Representation(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *aasResp.Data.Result.ID, fmp4RepAudio)
	ErrorHandler(fmp4RepAudioResp.Status, err)

	startResp, err = dashService.Start(*dashManifestResp.Data.Result.ID)
	ErrorHandler(startResp.Status, err)

	status = ""
	for status != "FINISHED" {
		time.Sleep(5 * time.Second)
		statusResp, err := dashService.RetrieveStatus(*dashManifestResp.Data.Result.ID)
		if err != nil {
			fmt.Println("error in Manifest Status")
			fmt.Println(err)
			return
		}
		// Polling and Printing out the response
		fmt.Printf("%+v\n", statusResp)
		status = *statusResp.Data.Result.Status
		if status == "ERROR" {
			fmt.Println("error in Manifest Status")
			fmt.Printf("%+v\n", statusResp)
			return
		}
	}

	// Delete Encoding
	deleteResp, err := encodingS.Delete(*encodingResp.Data.Result.ID)
	ErrorHandler(deleteResp.Status, err)
}
