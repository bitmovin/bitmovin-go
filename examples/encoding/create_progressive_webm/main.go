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

	vp8S := services.NewVP8CodecConfigurationService(bitmovin)
	videoConfig := &models.VP8CodecConfiguration{
		Name:      stringToPtr("example_vp8_codec_configuration"),
		Bitrate:   intToPtr(1000000),
		FrameRate: floatToPtr(25.0),
		Width:     intToPtr(640),
		Height:    intToPtr(360),
	}
	videoResp, err := vp8S.Create(videoConfig)
	errorHandler(err)

	vorbisS := services.NewVorbisCodecConfigurationService(bitmovin)
	vorbisConfig := &models.VorbisCodecConfiguration{
		Name:         stringToPtr("example_audio_codec_configuration"),
		Bitrate:      intToPtr(128000),
		SamplingRate: floatToPtr(48000.0),
	}
	vorbisResp, err := vorbisS.Create(vorbisConfig)
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
	videoStream := &models.Stream{
		CodecConfigurationID: videoResp.Data.Result.ID,
		InputStreams:         vis,
	}
	videoStreamResp, err := encodingS.AddStream(*encodingResp.Data.Result.ID, videoStream)
	errorHandler(err)

	ais := []models.InputStream{audioInputStream}
	audioStream := &models.Stream{
		CodecConfigurationID: vorbisResp.Data.Result.ID,
		InputStreams:         ais,
	}
	vorbisStreamResp, err := encodingS.AddStream(*encodingResp.Data.Result.ID, audioStream)
	errorHandler(err)

	aclEntry := models.ACLItem{
		Permission: bitmovintypes.ACLPermissionPublicRead,
	}
	acl := []models.ACLItem{aclEntry}

	videoMuxingStream := models.StreamItem{
		StreamID: videoStreamResp.Data.Result.ID,
	}
	audioMuxingStream := models.StreamItem{
		StreamID: vorbisStreamResp.Data.Result.ID,
	}

	videoMuxingOutput := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: stringToPtr("golang_test/video"),
		ACL:        acl,
	}

	webmMuxing := &models.ProgressiveWebMMuxing{
		Streams:  []models.StreamItem{videoMuxingStream, audioMuxingStream},
		Outputs:  []models.Output{videoMuxingOutput},
		Filename: stringToPtr("yourfilename.webm"),
	}
	_, err = encodingS.AddProgressiveWebMMuxing(*encodingResp.Data.Result.ID, webmMuxing)
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
