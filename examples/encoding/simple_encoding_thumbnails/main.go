package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/bitmovintypes"
	"github.com/bitmovin/bitmovin-go/models"
	"github.com/bitmovin/bitmovin-go/services"
)

func main() {
	// Creating Bitmovin object
	bitmovin := bitmovin.NewBitmovin("<YOUR BITMOVIN API KEY>", "https://api.bitmovin.com/v1/", 5)

	dateComponent := time.Now().Format(time.RFC3339)

	// Creating the HTTP Input
	httpsIS := services.NewHTTPSInputService(bitmovin)
	httpsInput := &models.HTTPSInput{
		Host: stringToPtr("<YOUR HTTP HOST>"), // eg. storage.googleapis.com
	}
	httpsResp, err := httpsIS.Create(httpsInput)
	errorHandler(httpsResp.Status, err)

	inputFilePath := "test.mp4" // eg. /path-to-your-file/test.mp4

	fmt.Printf("Creating GCS Outpput")
	gcsOS := services.NewGCSOutputService(bitmovin)
	gcsOutput := &models.GCSOutput{
		AccessKey:  stringToPtr("<YOUR GCP ACCESS KEY>"),
		SecretKey:  stringToPtr("<YOUR SECRET KEY>"),
		BucketName: stringToPtr("<YOUR BUCKET NAME>"),
	}
	outputResponse, err := gcsOS.Create(gcsOutput)

	errorHandler(outputResponse.Status, err)

	encodingS := services.NewEncodingService(bitmovin)
	encoding := &models.Encoding{
		Name:        stringToPtr("example encoding"),
		CloudRegion: bitmovintypes.CloudRegionGoogleUSEast1,
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
	errorHandler(videoStream1080pResp.Status, err)
	videoStream720pResp, err := encodingS.AddStream(*encodingResp.Data.Result.ID, videoStream720p)
	errorHandler(videoStream720pResp.Status, err)

	ais := []models.InputStream{audioInputStream}
	audioStream := &models.Stream{
		CodecConfigurationID: aacResp.Data.Result.ID,
		InputStreams:         ais,
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

	// Add Thumbnails to the Stream
	thumbOutput := &models.Output{
		OutputID:   outputResponse.Data.Result.ID,
		OutputPath: stringToPtr(fmt.Sprintf("golang_test/%s/thumbs", dateComponent)),
		ACL:        acl,
	}
	thumb1080 := models.NewThumbnail(400, []float64{3, 5, 30}, []models.Output{*thumbOutput}).Builder().
		Pattern("thumbnail-%number%.png").Build()

	if _, err := encodingS.AddThumbnail(*encodingResp.Data.Result.ID, *videoStream1080pResp.Data.Result.ID, thumb1080); err != nil {
		log.Fatalf("Error creating 1080p Thumbnail resource")
	}

	videoTSMuxing1080p := &models.TSMuxing{
		SegmentLength: floatToPtr(4.0),
		SegmentNaming: stringToPtr("seg_%number%.ts"),
		Streams:       []models.StreamItem{videoMuxingStream1080p},
		Outputs:       []models.Output{videoMuxing1080pOutput},
	}
	videoTSMuxing1080pResp, err := encodingS.AddTSMuxing(*encodingResp.Data.Result.ID, videoTSMuxing1080p)
	errorHandler(videoTSMuxing1080pResp.Status, err)

	videoTSMuxing720p := &models.TSMuxing{
		SegmentLength: floatToPtr(4.0),
		SegmentNaming: stringToPtr("seg_%number%.ts"),
		Streams:       []models.StreamItem{videoMuxingStream720p},
		Outputs:       []models.Output{videoMuxing720pOutput},
	}
	videoTSMuxing720pResp, err := encodingS.AddTSMuxing(*encodingResp.Data.Result.ID, videoTSMuxing720p)
	errorHandler(videoTSMuxing720pResp.Status, err)

	audioTSMuxing := &models.TSMuxing{
		SegmentLength: floatToPtr(4.0),
		SegmentNaming: stringToPtr("seg_%number%.ts"),
		Streams:       []models.StreamItem{audioMuxingStream},
		Outputs:       []models.Output{audioMuxingOutput},
	}
	audioTSMuxingResp, err := encodingS.AddTSMuxing(*encodingResp.Data.Result.ID, audioTSMuxing)
	errorHandler(audioTSMuxingResp.Status, err)

	startResp, err := encodingS.Start(*encodingResp.Data.Result.ID)
	errorHandler(startResp.Status, err)

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
		SegmentPath:     stringToPtr("../audio"),
		EncodingID:      encodingResp.Data.Result.ID,
		StreamID:        aacStreamResp.Data.Result.ID,
		MuxingID:        audioTSMuxingResp.Data.Result.ID,
	}
	audioMediaInfoResp, err := hlsService.AddMediaInfo(*hlsManifestResp.Data.Result.ID, audioMediaInfo)
	errorHandler(audioMediaInfoResp.Status, err)

	video1080pStreamInfo := &models.StreamInfo{
		Audio:       stringToPtr("audio_group"),
		SegmentPath: stringToPtr("../video/1080p"),
		URI:         stringToPtr("video_hi.m3u8"),
		EncodingID:  encodingResp.Data.Result.ID,
		StreamID:    videoStream1080pResp.Data.Result.ID,
		MuxingID:    videoTSMuxing1080pResp.Data.Result.ID,
	}
	video1080pStreamInfoResponse, err := hlsService.AddStreamInfo(*hlsManifestResp.Data.Result.ID, video1080pStreamInfo)
	errorHandler(video1080pStreamInfoResponse.Status, err)

	video720pStreamInfo := &models.StreamInfo{
		Audio:       stringToPtr("audio_group"),
		SegmentPath: stringToPtr("../video/720p"),
		URI:         stringToPtr("video_lo.m3u8"),
		EncodingID:  encodingResp.Data.Result.ID,
		StreamID:    videoStream720pResp.Data.Result.ID,
		MuxingID:    videoTSMuxing720pResp.Data.Result.ID,
	}
	video720pStreamInfoResponse, err := hlsService.AddStreamInfo(*hlsManifestResp.Data.Result.ID, video720pStreamInfo)
	errorHandler(video720pStreamInfoResponse.Status, err)

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

	//// Delete Encoding
	//deleteResp, err := encodingS.Delete(*encodingResp.Data.Result.ID)
	//errorHandler(deleteResp.Status, err)
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
