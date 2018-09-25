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
	inputFilePath  = "/path/to/your/input/file.mov"
	outputBasePath = "golang_example"
	outputFilePath = "myOutputFile.mp4"
)

func main() {
	// Creating Bitmovin object
	bitmovin := bitmovin.NewBitmovin("<INSERT YOUR API KEY>", "https://api.bitmovin.com/v1/", 5)

	// Creating the GCS Input
	gcsIS := services.NewGCSInputService(bitmovin)
	gcsInput := &models.GCSInput{
		AccessKey:  stringToPtr("<INSERT YOUR ACCESS KEY>"),
		SecretKey:  stringToPtr("<INSERT YOUR SECRET KEY>"),
		BucketName: stringToPtr("<INSERT YOUR BUCKET NAME>"),
	}

	inputResp, err := gcsIS.Create(gcsInput)
	errorHandler(err)

	// Creating the GCS Output
	gcsOS := services.NewGCSOutputService(bitmovin)
	gcsOutput := &models.GCSOutput{
		AccessKey:   stringToPtr("<INSERT YOUR ACCESS KEY>"),
		SecretKey:   stringToPtr("<INSERT YOUR SECRET KEY>"),
		BucketName:  stringToPtr("<INSERT YOUR BUCKET NAME>"),
		CloudRegion: bitmovintypes.GoogleCloudRegionEuropeWest1,
	}
	gcsOutputResp, err := gcsOS.Create(gcsOutput)
	errorHandler(err)

	t := time.Now()
	outputBasePath := outputBasePath + "/" + t.Format("2006-01-02-15-04-05") + "/"

	encodingS := services.NewEncodingService(bitmovin)
	encoding := &models.Encoding{
		Name:           stringToPtr("HDR Master Display Information Encoding"),
		CloudRegion:    bitmovintypes.CloudRegionGoogleEuropeWest1,
		EncoderVersion: "STABLE",
	}
	encodingResp, err := encodingS.Create(encoding)
	errorHandler(err)

	h265S := services.NewH265CodecConfigurationService(bitmovin)
	videoConfig := newHdrVideoConfig()
	videoResp, err := h265S.Create(videoConfig)
	errorHandler(err)

	aacS := services.NewAACCodecConfigurationService(bitmovin)
	audioConfig := newAudioConfig()
	aacResp, err := aacS.Create(audioConfig)
	errorHandler(err)

	videoInputStream := models.InputStream{
		InputID:       inputResp.Data.Result.ID,
		InputPath:     stringToPtr(inputFilePath),
		SelectionMode: bitmovintypes.SelectionModeAuto,
	}
	audioInputStream := models.InputStream{
		InputID:       inputResp.Data.Result.ID,
		InputPath:     stringToPtr(inputFilePath),
		SelectionMode: bitmovintypes.SelectionModeAuto,
	}

	vis := []models.InputStream{videoInputStream}

	videoStream := &models.Stream{
		CodecConfigurationID: videoResp.Data.Result.ID,
		InputStreams:         vis,
	}

	videoStreamResp, err := encodingS.AddStream(*encodingResp.Data.Result.ID, videoStream)
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

	videoMuxingStream := models.StreamItem{
		StreamID: videoStreamResp.Data.Result.ID,
	}
	audioMuxingStream := models.StreamItem{
		StreamID: aacStreamResp.Data.Result.ID,
	}

	combinedStreams := []models.StreamItem{videoMuxingStream, audioMuxingStream}

	encodingOutput := models.Output{
		OutputID:   gcsOutputResp.Data.Result.ID,
		OutputPath: stringToPtr(outputBasePath),
		ACL:        acl,
	}

	/*
		MP4 MUXINGS
	*/
	mp4Muxing := &models.MP4Muxing{
		Name:        stringToPtr("HDR MP4 muxing"),
		Description: stringToPtr("HDR MP4 muxing"),
		Outputs:     []models.Output{encodingOutput},
		Filename:    stringToPtr(outputFilePath),
		Streams:     combinedStreams,
	}

	encodingS.AddMP4Muxing(*encodingResp.Data.Result.ID, mp4Muxing)

	/*
		START ENCODING AND WAIT TO FOR IT TO BE FINISHED
	*/
	fmt.Printf("Starting encoding with id %s...\n", *encodingResp.Data.Result.ID)

	_, err = encodingS.Start(*encodingResp.Data.Result.ID)
	errorHandler(err)

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

func newHdrVideoConfig() *models.H265CodecConfiguration {

	colorConfig := models.ColorConfig{
		ColorTransfer:  bitmovintypes.ColorTransferSMPTE2084,
		ColorPrimaries: bitmovintypes.ColorPrimariesBT2020,
		ColorSpace:     bitmovintypes.ColorSpaceBT2020_NCL,
	}

	return &models.H265CodecConfiguration{
		Name:         stringToPtr("HEVC_HDR10"),
		Description:  stringToPtr("HEVC_HDR_10"),
		Profile:      bitmovintypes.H265ProfileMain10,
		Width:        intToPtr(3840),
		Height:       intToPtr(2160),
		FrameRate:    floatToPtr(25.0),
		CRF:          floatToPtr(19.0),
		BAdapt:       bitmovintypes.BAdaptFull,
		RCLookahead:  intToPtr(40),
		RefFrames:    intToPtr(5),
		MotionSearch: bitmovintypes.MotionSearchStar,
		SubMe:        intToPtr(4),
		TUInterDepth: bitmovintypes.TUInterDepth3,
		TUIntraDepth: bitmovintypes.TUIntraDepth3,
		MaxCTUSize:   bitmovintypes.MaxCTUSize64,
		BFrames:      intToPtr(4),
		SAO:          boolToPtr(true),
		WeightPredictionOnPSlice:    boolToPtr(true),
		WeightPredictionOnBSlice:    boolToPtr(false),
		ColorConfig:                 colorConfig,
		MasterDisplay:               stringToPtr("G(8500,39850)B(6550,2300)R(35400,14600)WP(15635,16450)L(100000000000,0)"),
		MaxContentLightLevel:        intToPtr(1000),
		MaxPictureAverageLightLevel: intToPtr(400),
		HDR:         boolToPtr(true),
		PixelFormat: bitmovintypes.PixelFormatYUV420P10LE,
	}
}

func newAudioConfig() *models.AACCodecConfiguration {
	return &models.AACCodecConfiguration{
		Name:         stringToPtr("AAC Audio Config"),
		Description:  stringToPtr("AAC Audio Config"),
		Bitrate:      intToPtr(96000),
		SamplingRate: floatToPtr(48000),
	}
}
