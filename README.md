# [![bitmovin](https://cloudfront-prod.bitmovin.com/wp-content/themes/Bitmovin-V-0.1/images/logo3.png)](http://www.bitmovin.com)
Golang-Client which enables you to seamlessly integrate the Bitmovin API into your existing projects

Installation 
------------
 
Run `go get github.com/bitmovin/bitmovin-go`

Also feel free to use your favorite go dependency manager such as glide.

Example
-----
The following example creates a simple transcoding job with a HTTP Input and a S3 Output ([create_simple_encoding.go](https://github.com/bitmovin/bitmovin-go/tree/master/examples/create_simple_encoding.go)):
```go
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

	// Creating the HTTP Input
	httpIS := services.NewHTTPInputService(bitmovin)
	httpInput := &models.HTTPInput{
		Host: "YOUR HTTP HOST",
	}
	httpResp, err := httpIS.Create(httpInput)
	errorHandler(httpResp.Status, err)

	s3OS := services.NewS3OutputService(bitmovin)
	s3Output := &models.S3Output{
		AccessKey:   "YOUR S3 ACCESS KEY",
		SecretKey:   "YOUR S3 SECRET KEY",
		BucketName:  "YOUR S3 BUCKET",
		CloudRegion: bitmovintypes.AWSCloudRegionEUWest1,
	}
	s3OutputResp, err := s3OS.Create(s3Output)
	errorHandler(s3OutputResp.Status, err)

	encodingS := services.NewEncodingService(bitmovin)
	encoding := &models.Encoding{
		Name:        "example encoding",
		CloudRegion: bitmovintypes.CloudRegionGoogleEuropeWest1,
	}
	encodingResp, err := encodingS.Create(encoding)
	errorHandler(encodingResp.Status, err)

	h264S := services.NewH264CodecConfigurationService(bitmovin)
	video1080pConfig := &models.H264CodecConfiguration{
		Name:      "example_video_codec_configuration_1080p",
		Bitrate:   4800000,
		FrameRate: 25.0,
		Width:     1920,
		Height:    1080,
		Profile:   bitmovintypes.H264ProfileHigh,
	}
	video720Config := &models.H264CodecConfiguration{
		Name:      "example_video_codec_configuration_720p",
		Bitrate:   2400000,
		FrameRate: 25.0,
		Width:     1280,
		Height:    720,
		Profile:   bitmovintypes.H264ProfileHigh,
	}
	video1080pResp, err := h264S.Create(video1080pConfig)
	errorHandler(video1080pResp.Status, err)
	video720Resp, err := h264S.Create(video720Config)
	errorHandler(video720Resp.Status, err)

	aacS := services.NewAACCodecConfigurationService(bitmovin)
	aacConfig := &models.AACCodecConfiguration{
		Name:         "example_audio_codec_configuration",
		Bitrate:      128000,
		SamplingRate: 48000.0,
	}
	aacResp, err := aacS.Create(aacConfig)
	errorHandler(aacResp.Status, err)

	videoInputStream := models.InputStream{
		InputID:       httpResp.Data.Result.ID,
		InputPath:     "YOUR INPUT FILE PATH AND LOCATION",
		SelectionMode: bitmovintypes.SelectionModeAuto,
	}
	audioInputStream := models.InputStream{
		InputID:       httpResp.Data.Result.ID,
		InputPath:     "YOUR INPUT FILE PATH AND LOCATION",
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

	videoStream1080pResp, err := encodingS.AddStream(encodingResp.Data.Result.ID, videoStream1080p)
	errorHandler(videoStream1080pResp.Status, err)
	videoStream720pResp, err := encodingS.AddStream(encodingResp.Data.Result.ID, videoStream720p)
	errorHandler(videoStream720pResp.Status, err)

	ais := []models.InputStream{audioInputStream}
	audioStream := &models.Stream{
		CodecConfigurationID: aacResp.Data.Result.ID,
		InputStreams:         ais,
	}
	aacStreamResp, err := encodingS.AddStream(encodingResp.Data.Result.ID, audioStream)
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
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: "golang_test/video/1080p",
		ACL:        acl,
	}
	videoMuxing720pOutput := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: "golang_test/video/720p",
		ACL:        acl,
	}
	audioMuxingOutput := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: "golang_test/audio",
		ACL:        acl,
	}

	videoMuxing1080p := &models.FMP4Muxing{
		SegmentLength:   4.0,
		SegmentNaming:   "seg_%number%.m4s",
		InitSegmentName: "init.mp4",
		Streams:         []models.StreamItem{videoMuxingStream1080p},
		Outputs:         []models.Output{videoMuxing1080pOutput},
	}
	videoMuxing1080pResp, err := encodingS.AddFMP4Muxing(encodingResp.Data.Result.ID, videoMuxing1080p)
	errorHandler(videoMuxing1080pResp.Status, err)

	videoMuxing720p := &models.FMP4Muxing{
		SegmentLength:   4.0,
		SegmentNaming:   "seg_%number%.m4s",
		InitSegmentName: "init.mp4",
		Streams:         []models.StreamItem{videoMuxingStream720p},
		Outputs:         []models.Output{videoMuxing720pOutput},
	}
	videoMuxing720pResp, err := encodingS.AddFMP4Muxing(encodingResp.Data.Result.ID, videoMuxing720p)
	errorHandler(videoMuxing720pResp.Status, err)

	audioMuxing := &models.FMP4Muxing{
		SegmentLength:   4.0,
		SegmentNaming:   "seg_%number%.m4s",
		InitSegmentName: "init.mp4",
		Streams:         []models.StreamItem{audioMuxingStream},
		Outputs:         []models.Output{audioMuxingOutput},
	}
	audioMuxingResp, err := encodingS.AddFMP4Muxing(encodingResp.Data.Result.ID, audioMuxing)
	errorHandler(audioMuxingResp.Status, err)

	startResp, err := encodingS.Start(encodingResp.Data.Result.ID)
	errorHandler(startResp.Status, err)

	var status string
	status = ""
	for status != "FINISHED" {
		time.Sleep(10 * time.Second)
		statusResp, err := encodingS.RetrieveStatus(encodingResp.Data.Result.ID)
		if err != nil {
			fmt.Println("error in Encoding Status")
			fmt.Println(err)
			return
		}
		// Polling and Printing out the response
		fmt.Printf("%+v\n", statusResp)
		status = statusResp.Data.Result.Status
		if status == "ERROR" {
			fmt.Println("error in Encoding Status")
			fmt.Printf("%+v\n", statusResp)
			return
		}
	}

	manifestOutput := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: "golang_test/manifest",
		ACL:        acl,
	}
	dashManifest := &models.DashManifest{
		ManifestName: "your_manifest_name.mpd",
		Outputs:      []models.Output{manifestOutput},
	}
	dashService := services.NewDashManifestService(bitmovin)
	dashManifestResp, err := dashService.Create(dashManifest)
	errorHandler(dashManifestResp.Status, err)

	period := &models.Period{}
	periodResp, err := dashService.AddPeriod(dashManifestResp.Data.Result.ID, period)
	errorHandler(periodResp.Status, err)

	vas := &models.VideoAdaptationSet{}
	vasResp, err := dashService.AddVideoAdaptationSet(dashManifestResp.Data.Result.ID, periodResp.Data.Result.ID, vas)
	errorHandler(vasResp.Status, err)

	aas := &models.AudioAdaptationSet{
		Language: "en",
	}
	aasResp, err := dashService.AddAudioAdaptationSet(dashManifestResp.Data.Result.ID, periodResp.Data.Result.ID, aas)
	errorHandler(aasResp.Status, err)

	fmp4Rep1080 := &models.FMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    videoMuxing1080pResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: "golang_test/video/1080p",
	}
	fmp4Rep1080Resp, err := dashService.AddFMP4Representation(dashManifestResp.Data.Result.ID, periodResp.Data.Result.ID, vasResp.Data.Result.ID, fmp4Rep1080)
	errorHandler(fmp4Rep1080Resp.Status, err)

	fmp4Rep720 := &models.FMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    videoMuxing720pResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: "golang_test/video/720p",
	}
	fmp4Rep720Resp, err := dashService.AddFMP4Representation(dashManifestResp.Data.Result.ID, periodResp.Data.Result.ID, vasResp.Data.Result.ID, fmp4Rep720)
	errorHandler(fmp4Rep720Resp.Status, err)

	fmp4RepAudio := &models.FMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    audioMuxingResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: "golang_test/audio",
	}
	fmp4RepAudioResp, err := dashService.AddFMP4Representation(dashManifestResp.Data.Result.ID, periodResp.Data.Result.ID, aasResp.Data.Result.ID, fmp4RepAudio)
	errorHandler(fmp4RepAudioResp.Status, err)

	startResp, err = dashService.Start(dashManifestResp.Data.Result.ID)
	errorHandler(startResp.Status, err)

	status = ""
	for status != "FINISHED" {
		time.Sleep(5 * time.Second)
		statusResp, err := dashService.RetrieveStatus(dashManifestResp.Data.Result.ID)
		if err != nil {
			fmt.Println("error in Manifest Status")
			fmt.Println(err)
			return
		}
		// Polling and Printing out the response
		fmt.Printf("%+v\n", statusResp)
		status = statusResp.Data.Result.Status
		if status == "ERROR" {
			fmt.Println("error in Manifest Status")
			fmt.Printf("%+v\n", statusResp)
			return
		}
	}

	// Delete Encoding
	deleteResp, err := encodingS.Delete(encodingResp.Data.Result.ID)
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

```
