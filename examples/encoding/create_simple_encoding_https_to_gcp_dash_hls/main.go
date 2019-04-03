package main

import (
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
	bitmovin := bitmovin.NewBitmovin("<YOUR API KEY>", "https://api.bitmovin.com/v1/", 5)

	dateComponent := time.Now().Format(time.RFC3339)

	// Creating the HTTP Input
	httpsIS := services.NewHTTPSInputService(bitmovin)
	httpsInput := &models.HTTPSInput{
		Host: stringToPtr("<YOUR HTTPS HOSTNAME>"), // eg. storage.googleapis.com
	}
	httpsResp, err := httpsIS.Create(httpsInput)
	errorHandler(err)

	inputFilePath := "/test.mp4" // eg. /path-to-your-file/test.mp4

	fmt.Printf("Creating GCS Output")
	gcsOS := services.NewGCSOutputService(bitmovin)
	gcsOutput := &models.GCSOutput{
		AccessKey:  stringToPtr("<YOUR ACCESS KEY>"),
		SecretKey:  stringToPtr("<YOUR GCP SECRET KEY>"),
		BucketName: stringToPtr("<YOUR BUCKET NAME>"),
	}
	outputResponse, err := gcsOS.Create(gcsOutput)

	errorHandler(err)

	encodingS := services.NewEncodingService(bitmovin)
	encoding := &models.Encoding{
		Name:        stringToPtr("example encoding"),
		CloudRegion: bitmovintypes.CloudRegionGoogleUSEast1,
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
	aacConfig := &models.AACCodecConfiguration{
		Name:         stringToPtr("example_audio_codec_configuration"),
		Bitrate:      intToPtr(128000),
		SamplingRate: floatToPtr(48000.0),
	}
	aacResp, err := aacS.Create(aacConfig)
	errorHandler(err)

	videoInputStream := models.InputStream{
		InputID:       httpsResp.Data.Result.ID,
		InputPath:     stringToPtr(inputFilePath),
		SelectionMode: bitmovintypes.SelectionModeAuto,
	}
	audioInputStream := models.InputStream{
		InputID:       httpsResp.Data.Result.ID,
		InputPath:     stringToPtr(inputFilePath),
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
	audioStream := &models.Stream{
		CodecConfigurationID: aacResp.Data.Result.ID,
		InputStreams:         ais,
	}
	aacStreamResp, err := encodingS.AddStream(*encodingResp.Data.Result.ID, audioStream)
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
		OutputID:   outputResponse.Data.Result.ID,
		OutputPath: stringToPtr(fmt.Sprintf("golang_test/%s/video/1080p", dateComponent)),
		ACL:        acl,
	}
	videoMuxing720pOutput := models.Output{
		OutputID:   outputResponse.Data.Result.ID,
		OutputPath: stringToPtr(fmt.Sprintf("golang_test/%s/video/720p", dateComponent)),
		ACL:        acl,
	}
	audioMuxingOutput := models.Output{
		OutputID:   outputResponse.Data.Result.ID,
		OutputPath: stringToPtr(fmt.Sprintf("golang_test/%s/audio", dateComponent)),
		ACL:        acl,
	}

	videoFMP4Muxing1080p := &models.FMP4Muxing{
		SegmentLength:   floatToPtr(4.0),
		SegmentNaming:   stringToPtr("seg_%number%.m4s"),
		InitSegmentName: stringToPtr("init.mp4"),
		Streams:         []models.StreamItem{videoMuxingStream1080p},
		Outputs:         []models.Output{videoMuxing1080pOutput},
	}
	videoFMP4Muxing1080pResp, err := encodingS.AddFMP4Muxing(*encodingResp.Data.Result.ID, videoFMP4Muxing1080p)
	errorHandler(err)

	videoFMP4Muxing720p := &models.FMP4Muxing{
		SegmentLength:   floatToPtr(4.0),
		SegmentNaming:   stringToPtr("seg_%number%.m4s"),
		InitSegmentName: stringToPtr("init.mp4"),
		Streams:         []models.StreamItem{videoMuxingStream720p},
		Outputs:         []models.Output{videoMuxing720pOutput},
	}
	videoFMP4Muxing720pResp, err := encodingS.AddFMP4Muxing(*encodingResp.Data.Result.ID, videoFMP4Muxing720p)
	errorHandler(err)

	audioFMP4Muxing := &models.FMP4Muxing{
		SegmentLength:   floatToPtr(4.0),
		SegmentNaming:   stringToPtr("seg_%number%.m4s"),
		InitSegmentName: stringToPtr("init.mp4"),
		Streams:         []models.StreamItem{audioMuxingStream},
		Outputs:         []models.Output{audioMuxingOutput},
	}
	audioFMP4MuxingResp, err := encodingS.AddFMP4Muxing(*encodingResp.Data.Result.ID, audioFMP4Muxing)
	errorHandler(err)

	videoTSMuxing1080p := &models.TSMuxing{
		SegmentLength: floatToPtr(4.0),
		SegmentNaming: stringToPtr("seg_%number%.ts"),
		Streams:       []models.StreamItem{videoMuxingStream1080p},
		Outputs:       []models.Output{videoMuxing1080pOutput},
	}
	videoTSMuxing1080pResp, err := encodingS.AddTSMuxing(*encodingResp.Data.Result.ID, videoTSMuxing1080p)
	errorHandler(err)

	videoTSMuxing720p := &models.TSMuxing{
		SegmentLength: floatToPtr(4.0),
		SegmentNaming: stringToPtr("seg_%number%.ts"),
		Streams:       []models.StreamItem{videoMuxingStream720p},
		Outputs:       []models.Output{videoMuxing720pOutput},
	}
	videoTSMuxing720pResp, err := encodingS.AddTSMuxing(*encodingResp.Data.Result.ID, videoTSMuxing720p)
	errorHandler(err)

	audioTSMuxing := &models.TSMuxing{
		SegmentLength: floatToPtr(4.0),
		SegmentNaming: stringToPtr("seg_%number%.ts"),
		Streams:       []models.StreamItem{audioMuxingStream},
		Outputs:       []models.Output{audioMuxingOutput},
	}
	audioTSMuxingResp, err := encodingS.AddTSMuxing(*encodingResp.Data.Result.ID, audioTSMuxing)
	errorHandler(err)

	_, err = encodingS.Start(*encodingResp.Data.Result.ID)
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
		fmt.Printf("%+v\n", statusResp)
		status = *statusResp.Data.Result.Status
		if status == "ERROR" {
			fmt.Println("error in Encoding Status")
			fmt.Printf("%+v\n", statusResp)
			return
		}
	}
	fmt.Println("Encoding Finished - writing Manifests now")

	manifestOutput := models.Output{
		OutputID:   outputResponse.Data.Result.ID,
		OutputPath: stringToPtr(fmt.Sprintf("golang_test/%s/manifest", dateComponent)),
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
		MuxingID:    videoFMP4Muxing1080pResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: stringToPtr("../video/1080p"),
	}
	_, err = dashService.AddFMP4Representation(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *vasResp.Data.Result.ID, fmp4Rep1080)
	errorHandler(err)

	fmp4Rep720 := &models.FMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    videoFMP4Muxing720pResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: stringToPtr("../video/720p"),
	}
	_, err = dashService.AddFMP4Representation(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *vasResp.Data.Result.ID, fmp4Rep720)
	errorHandler(err)

	fmp4RepAudio := &models.FMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    audioFMP4MuxingResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: stringToPtr("../audio"),
	}
	_, err = dashService.AddFMP4Representation(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *aasResp.Data.Result.ID, fmp4RepAudio)
	errorHandler(err)

	_, err = dashService.Start(*dashManifestResp.Data.Result.ID)
	errorHandler(err)

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

	hlsManifest := &models.HLSManifest{
		ManifestName: stringToPtr("your_manifest_name.m3u8"),
		Outputs:      []models.Output{manifestOutput},
	}
	hlsService := services.NewHLSManifestService(bitmovin)
	hlsManifestResp, err := hlsService.Create(hlsManifest)
	errorHandler(err)

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
		SegmentPath:     stringToPtr("../audio"),
		EncodingID:      encodingResp.Data.Result.ID,
		StreamID:        aacStreamResp.Data.Result.ID,
		MuxingID:        audioTSMuxingResp.Data.Result.ID,
	}
	_, err = hlsService.AddMediaInfo(*hlsManifestResp.Data.Result.ID, audioMediaInfo)
	errorHandler(err)

	video1080pStreamInfo := &models.StreamInfo{
		Audio:       stringToPtr("audio_group"),
		SegmentPath: stringToPtr("../video/1080p"),
		URI:         stringToPtr("video_hi.m3u8"),
		EncodingID:  encodingResp.Data.Result.ID,
		StreamID:    videoStream1080pResp.Data.Result.ID,
		MuxingID:    videoTSMuxing1080pResp.Data.Result.ID,
	}
	_, err = hlsService.AddStreamInfo(*hlsManifestResp.Data.Result.ID, video1080pStreamInfo)
	errorHandler(err)

	video720pStreamInfo := &models.StreamInfo{
		Audio:       stringToPtr("audio_group"),
		SegmentPath: stringToPtr("../video/720p"),
		URI:         stringToPtr("video_lo.m3u8"),
		EncodingID:  encodingResp.Data.Result.ID,
		StreamID:    videoStream720pResp.Data.Result.ID,
		MuxingID:    videoTSMuxing720pResp.Data.Result.ID,
	}
	_, err = hlsService.AddStreamInfo(*hlsManifestResp.Data.Result.ID, video720pStreamInfo)
	errorHandler(err)

	hlsService.Start(*hlsManifestResp.Data.Result.ID)

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
		fmt.Printf("%+v\n", statusResp)
		status = *statusResp.Data.Result.Status
		if status == "ERROR" {
			fmt.Println("error in Manifest Status")
			fmt.Printf("%+v\n", statusResp)
			return
		}
	}

	// Delete Encoding
	_, err = encodingS.Delete(*encodingResp.Data.Result.ID)
	errorHandler(err)
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
