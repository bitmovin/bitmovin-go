package main

import (
	"fmt"
	"time"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/bitmovintypes"
	"github.com/bitmovin/bitmovin-go/models"
	"github.com/bitmovin/bitmovin-go/services"
)

func main() {
	// Creating Bitmovin object
	bitmovin := bitmovin.NewBitmovin("YOUR API KEY", "https://api.bitmovin.com/v1/", 5)

	// Creating the HTTPS Input
	httpIS := services.NewHTTPSInputService(bitmovin)
	httpInput := &models.HTTPSInput{
		Host: stringToPtr("YOUR HTTP HOST"),
	}
	httpResp, err := httpIS.Create(httpInput)
	errorHandler(httpResp.Status, err)

	gcsOS := services.NewGCSOutputService(bitmovin)
	gcsOutput := &models.GCSOutput{
		AccessKey:   stringToPtr("YOUR_ACCESS_KEY"),
		SecretKey:   stringToPtr("YOUR_SECRET_KEY"),
		BucketName:  stringToPtr("YOUR_BUCKET_NAME"),
		CloudRegion: bitmovintypes.GoogleCloudRegionEuropeWest1,
	}
	gcsOutputResp, err := gcsOS.Create(gcsOutput)
	errorHandler(gcsOutputResp.Status, err)

	t := time.Now()
	outputBasePath := "golang_live_drm_test_" + t.Format("2006-01-02-15-04-05")

	encodingS := services.NewEncodingService(bitmovin)
	encoding := &models.Encoding{
		Name:        stringToPtr("Presence Condition Encoding"),
		CloudRegion: bitmovintypes.CloudRegionGoogleEuropeWest1,
	}
	encodingResp, err := encodingS.Create(encoding)
	errorHandler(encodingResp.Status, err)

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
	errorHandler(video1080pResp.Status, err)
	video720Resp, err := h264S.Create(video720Config)
	errorHandler(video720Resp.Status, err)

	aacS := services.NewAACCodecConfigurationService(bitmovin)
	aacConfig := &models.AACCodecConfiguration{
		Name:         stringToPtr("example_audio_codec_configuration"),
		Bitrate:      intToPtr(128000),
		SamplingRate: floatToPtr(48000.0),
	}
	aacResp, err := aacS.Create(aacConfig)
	errorHandler(aacResp.Status, err)

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
	errorHandler(videoStream1080pResp.Status, err)
	videoStream720pResp, err := encodingS.AddStream(*encodingResp.Data.Result.ID, videoStream720p)
	errorHandler(videoStream720pResp.Status, err)

	ais := []models.InputStream{audioInputStream}
	audioStream := &models.Stream{
		CodecConfigurationID: aacResp.Data.Result.ID,
		InputStreams:         ais,
		Conditions:           models.NewAttributeCondition(bitmovintypes.ConditionAttributeInputStream, "==", "true"),
	}
	aacStreamResp, err := encodingS.AddStream(*encodingResp.Data.Result.ID, audioStream)
	errorHandler(aacStreamResp.Status, err)

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

	/*
	FMP4 MUXINGS
	 */
	videoFMP4Muxing1080pOutput := models.Output{
		OutputID:   gcsOutputResp.Data.Result.ID,
		OutputPath: stringToPtr(outputBasePath + "/video/dash/1080p"),
		ACL:        acl,
	}
	videoFMP4Muxing720pOutput := models.Output{
		OutputID:   gcsOutputResp.Data.Result.ID,
		OutputPath: stringToPtr(outputBasePath + "/video/dash/720p"),
		ACL:        acl,
	}
	audioFMP4MuxingOutput := models.Output{
		OutputID:   gcsOutputResp.Data.Result.ID,
		OutputPath: stringToPtr(outputBasePath + "/dash/audio"),
		ACL:        acl,
	}

	videoFMP4Muxing1080p := &models.FMP4Muxing{
		SegmentLength:   floatToPtr(4.0),
		SegmentNaming:   stringToPtr("seg_%number%.m4s"),
		InitSegmentName: stringToPtr("init.mp4"),
		Streams:         []models.StreamItem{videoMuxingStream1080p},
		Outputs:         []models.Output{videoFMP4Muxing1080pOutput},
	}
	videoFMP4Muxing1080pResp, err := encodingS.AddFMP4Muxing(*encodingResp.Data.Result.ID, videoFMP4Muxing1080p)
	errorHandler(videoFMP4Muxing1080pResp.Status, err)

	videoFMP4Muxing720p := &models.FMP4Muxing{
		SegmentLength:   floatToPtr(4.0),
		SegmentNaming:   stringToPtr("seg_%number%.m4s"),
		InitSegmentName: stringToPtr("init.mp4"),
		Streams:         []models.StreamItem{videoMuxingStream720p},
		Outputs:         []models.Output{videoFMP4Muxing720pOutput},
	}
	videoFMP4Muxing720pResp, err := encodingS.AddFMP4Muxing(*encodingResp.Data.Result.ID, videoFMP4Muxing720p)
	errorHandler(videoFMP4Muxing720pResp.Status, err)

	audioFMP4Muxing := &models.FMP4Muxing{
		SegmentLength:   floatToPtr(4.0),
		SegmentNaming:   stringToPtr("seg_%number%.m4s"),
		InitSegmentName: stringToPtr("init.mp4"),
		Streams:         []models.StreamItem{audioMuxingStream},
		Outputs:         []models.Output{audioFMP4MuxingOutput},
	}
	audioFMP4MuxingResp, err := encodingS.AddFMP4Muxing(*encodingResp.Data.Result.ID, audioFMP4Muxing)
	errorHandler(audioFMP4MuxingResp.Status, err)

	/*
	TS MUXINGS
	 */
	videoTSMuxing1080pOutput := models.Output{
		OutputID:   gcsOutputResp.Data.Result.ID,
		OutputPath: stringToPtr(outputBasePath + "/video/hls/1080p"),
		ACL:        acl,
	}
	videoTSMuxing720pOutput := models.Output{
		OutputID:   gcsOutputResp.Data.Result.ID,
		OutputPath: stringToPtr(outputBasePath + "/video/hls/720p"),
		ACL:        acl,
	}
	audioTSMuxingOutput := models.Output{
		OutputID:   gcsOutputResp.Data.Result.ID,
		OutputPath: stringToPtr(outputBasePath + "/hls/audio"),
		ACL:        acl,
	}

	videoMuxing1080p := &models.TSMuxing{
		SegmentLength: floatToPtr(4.0),
		SegmentNaming: stringToPtr("seg_%number%.ts"),
		Streams:       []models.StreamItem{videoMuxingStream1080p},
		Outputs:       []models.Output{videoTSMuxing1080pOutput},
	}
	videoMuxing1080pResp, err := encodingS.AddTSMuxing(*encodingResp.Data.Result.ID, videoMuxing1080p)
	errorHandler(videoMuxing1080pResp.Status, err)

	videoMuxing720p := &models.TSMuxing{
		SegmentLength: floatToPtr(4.0),
		SegmentNaming: stringToPtr("seg_%number%.ts"),
		Streams:       []models.StreamItem{videoMuxingStream720p},
		Outputs:       []models.Output{videoTSMuxing720pOutput},
	}
	videoMuxing720pResp, err := encodingS.AddTSMuxing(*encodingResp.Data.Result.ID, videoMuxing720p)
	errorHandler(videoMuxing720pResp.Status, err)

	audioMuxing := &models.TSMuxing{
		SegmentLength: floatToPtr(4.0),
		SegmentNaming: stringToPtr("seg_%number%.ts"),
		Streams:       []models.StreamItem{audioMuxingStream},
		Outputs:       []models.Output{audioTSMuxingOutput},
	}
	audioMuxingResp, err := encodingS.AddTSMuxing(*encodingResp.Data.Result.ID, audioMuxing)
	errorHandler(audioMuxingResp.Status, err)


	/*
	START ENCODING AND WAIT TO FOR IT TO BE FINISHED
	 */
	fmt.Printf("Starting encoding with id %s...\n", *encodingResp.Data.Result.ID)

	startResp, err := encodingS.Start(*encodingResp.Data.Result.ID)
	errorHandler(startResp.Status, err)

	var status string
	status = ""
	fmt.Println("Waiting for encoding to be FINISHED...")
	for status != "FINISHED" {
		time.Sleep(10 * time.Second)
		statusResp, err := encodingS.RetrieveStatus(*encodingResp.Data.Result.ID)
		if err != nil {
			fmt.Println("error in Encoding Status")
			fmt.Println(err)
			return
		}
		// Polling and Printing out the response
		status = *statusResp.Data.Result.Status
		fmt.Printf("STATUS: %s\n", status)
		if status == "ERROR" {
			fmt.Println("error in Encoding Status")
			fmt.Printf("STATUS: %s\n", status)
			return
		}
	}
	fmt.Println("Encoding finished successfully!")

	/*
	MANIFEST GENERATION
	 */
	manifestOutput := models.Output{
		OutputID:   gcsOutputResp.Data.Result.ID,
		OutputPath: stringToPtr(outputBasePath + "/manifest"),
		ACL:        acl,
	}

	/*
	DASH MANIFEST
	 */
	dashManifest := &models.DashManifest{
		ManifestName: stringToPtr("your_manifest_name.mpd"),
		Outputs:      []models.Output{manifestOutput},
	}
	dashService := services.NewDashManifestService(bitmovin)
	dashManifestResp, err := dashService.Create(dashManifest)
	errorHandler(dashManifestResp.Status, err)

	period := &models.Period{}
	periodResp, err := dashService.AddPeriod(*dashManifestResp.Data.Result.ID, period)
	errorHandler(periodResp.Status, err)

	vas := &models.VideoAdaptationSet{}
	vasResp, err := dashService.AddVideoAdaptationSet(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, vas)
	errorHandler(vasResp.Status, err)

	aas := &models.AudioAdaptationSet{
		Language: stringToPtr("en"),
	}
	aasResp, err := dashService.AddAudioAdaptationSet(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, aas)
	errorHandler(aasResp.Status, err)

	fmp4Rep1080 := &models.FMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    videoFMP4Muxing1080pResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: stringToPtr("../video/dash/1080p"),
	}
	fmp4Rep1080Resp, err := dashService.AddFMP4Representation(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *vasResp.Data.Result.ID, fmp4Rep1080)
	errorHandler(fmp4Rep1080Resp.Status, err)

	fmp4Rep720 := &models.FMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    videoFMP4Muxing720pResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: stringToPtr("../video/dash/720p"),
	}
	fmp4Rep720Resp, err := dashService.AddFMP4Representation(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *vasResp.Data.Result.ID, fmp4Rep720)
	errorHandler(fmp4Rep720Resp.Status, err)

	fmp4RepAudio := &models.FMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    audioFMP4MuxingResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: stringToPtr("../audio/dash"),
	}
	fmp4RepAudioResp, err := dashService.AddFMP4Representation(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *aasResp.Data.Result.ID, fmp4RepAudio)
	errorHandler(fmp4RepAudioResp.Status, err)

	fmt.Printf("Starting DASH manifest generation with manifest id %s...\n", *dashManifestResp.Data.Result.ID)

	startResp, err = dashService.Start(*dashManifestResp.Data.Result.ID)
	errorHandler(startResp.Status, err)

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
		status = *statusResp.Data.Result.Status
		fmt.Printf("STATUS: %s\n", status)
		if status == "ERROR" {
			fmt.Println("error in Manifest Status")
			fmt.Printf("STATUS: %s\n", status)
			return
		}
	}
	fmt.Println("DASH manifest created successfully!")

	/*
	HLS MANIFEST
	 */
	hlsManifest := &models.HLSManifest{
		ManifestName: stringToPtr("your_manifest_name.m3u8"),
		Outputs:      []models.Output{manifestOutput},
	}
	hlsService := services.NewHLSManifestService(bitmovin)
	hlsManifestResp, err := hlsService.Create(hlsManifest)
	errorHandler(hlsManifestResp.Status, err)

	audioMediaInfo := &models.MediaInfo{
		Type:            bitmovintypes.MediaTypeAudio,
		URI:             stringToPtr("audio.m3u8"),
		GroupID:         stringToPtr("audio_group"),
		Language:        stringToPtr("en"),
		Name:            stringToPtr("Rendition Description"),
		IsDefault:       boolToPtr(false),
		Autoselect:      boolToPtr(false),
		Forced:          boolToPtr(false),
		Characteristics: []string{"public.accessibility.describes-video"},
		SegmentPath:     stringToPtr("../audio/hls"),
		EncodingID:      encodingResp.Data.Result.ID,
		StreamID:        aacStreamResp.Data.Result.ID,
		MuxingID:        audioMuxingResp.Data.Result.ID,
	}
	audioMediaInfoResp, err := hlsService.AddMediaInfo(*hlsManifestResp.Data.Result.ID, audioMediaInfo)
	errorHandler(audioMediaInfoResp.Status, err)

	video1080pStreamInfo := &models.StreamInfo{
		Audio:       stringToPtr("audio_group"),
		SegmentPath: stringToPtr("../video/hls/1080p"),
		URI:         stringToPtr("video_hi.m3u8"),
		EncodingID:  encodingResp.Data.Result.ID,
		StreamID:    videoStream1080pResp.Data.Result.ID,
		MuxingID:    videoMuxing1080pResp.Data.Result.ID,
	}
	video1080pStreamInfoResponse, err := hlsService.AddStreamInfo(*hlsManifestResp.Data.Result.ID, video1080pStreamInfo)
	errorHandler(video1080pStreamInfoResponse.Status, err)

	video720pStreamInfo := &models.StreamInfo{
		Audio:       stringToPtr("audio_group"),
		SegmentPath: stringToPtr("../video/hls/720p"),
		URI:         stringToPtr("video_lo.m3u8"),
		EncodingID:  encodingResp.Data.Result.ID,
		StreamID:    videoStream720pResp.Data.Result.ID,
		MuxingID:    videoMuxing720pResp.Data.Result.ID,
	}
	video720pStreamInfoResponse, err := hlsService.AddStreamInfo(*hlsManifestResp.Data.Result.ID, video720pStreamInfo)
	errorHandler(video720pStreamInfoResponse.Status, err)

	fmt.Printf("Starting HLS manifest generation with manifest id %s...\n", *dashManifestResp.Data.Result.ID)

	startResp, err = hlsService.Start(*hlsManifestResp.Data.Result.ID)
	errorHandler(startResp.Status, err)

	status = ""
	for status != "FINISHED" {
		time.Sleep(5 * time.Second)
		statusResp, err := hlsService.RetrieveStatus(*hlsManifestResp.Data.Result.ID)
		if err != nil {
			fmt.Println("error in Manifest Status")
			fmt.Println(err)
			return
		}
		// Polling and Printing out the response
		status = *statusResp.Data.Result.Status
		fmt.Printf("STATUS: %s\n", status)
		if status == "ERROR" {
			fmt.Println("error in Manifest Status")
			fmt.Printf("STATUS: %s\n", status)
			return
		}
	}
	fmt.Println("HLS manifest created successfully!")

	// Delete Encoding
	deleteResp, err := encodingS.Delete(*encodingResp.Data.Result.ID)
	errorHandler(deleteResp.Status, err)
}

func errorHandler(responseStatus bitmovintypes.ResponseStatus, err error) {
	if err != nil {
		fmt.Println("go error")
		fmt.Println(err)
	} else if responseStatus == "ERROR" {
		fmt.Println("api error")
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
