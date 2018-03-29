package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/bitmovintypes"
	"github.com/bitmovin/bitmovin-go/models"
	"github.com/bitmovin/bitmovin-go/services"
)

func main() {

	// Creating Bitmovin object
	bitmovin := bitmovin.NewBitmovin("YOUR_API_KEY", "https://api.bitmovin.com/v1/", 5)

	// Creating the HTTP Input
	httpIS := services.NewHTTPInputService(bitmovin)
	httpInput := &models.HTTPInput{
		Host: stringToPtr("YOUR_HTTP_INPUT_HOST"),
	}
	httpResp, err := httpIS.Create(httpInput)
	errorHandler(httpResp.Status, err)

	s3OS := services.NewS3OutputService(bitmovin)
	s3Output := &models.S3Output{
		AccessKey:   stringToPtr("YOUR_ACCESS_KEY"),
		SecretKey:   stringToPtr("YOUR_SECRET_KEY"),
		BucketName:  stringToPtr("YOUR_BUCKET_NAME"),
		CloudRegion: bitmovintypes.AWSCloudRegionEUWest1,
	}
	s3OutputResp, err := s3OS.Create(s3Output)
	errorHandler(s3OutputResp.Status, err)

	infrastructureS := services.NewInfrastructureService(bitmovin)
	infrastructure, err := infrastructureS.Retrieve("YOUR_INFRASTRUCTURE_ID")

	if err != nil {
		panic(err)
	}

	kubernetesConfigR := &models.KubernetesClusterConfigurationRequest{
		ParallelEncodings:  intToPtr(2),
		WorkersPerEncoding: intToPtr(2),
	}

	kubernetesConfigS := services.NewKubernetesClusterConfigurationService(bitmovin)
	kubernetesConfig, err := kubernetesConfigS.Upsert(infrastructure.ID, kubernetesConfigR)
	errorHandler(kubernetesConfig.Status, err)

	kubernetesConfigResp, err := kubernetesConfigS.Retrieve(infrastructure.ID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("KubernetesClusterConfig Response: %s\n", getAsJsonString(*kubernetesConfigResp))

	encodingS := services.NewEncodingService(bitmovin)
	encoding := &models.Encoding{
		Name:             stringToPtr("Kubernetes Cluster Config Encoding"),
		Description:      stringToPtr("Example Golang Client Encoding with Kubernetes Cluster Configuration set"),
		CloudRegion:      bitmovintypes.CloudRegionExternal,
		InfrastructureID: stringToPtr(infrastructure.ID),
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
		InputPath:     stringToPtr("YOUR_INPUT_FILE.mp4"),
		SelectionMode: bitmovintypes.SelectionModeAuto,
	}
	audioInputStream := models.InputStream{
		InputID:       httpResp.Data.Result.ID,
		InputPath:     stringToPtr("YOUR_INPUT_FILE.mp4"),
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
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: stringToPtr("golang_test/video/1080p"),
		ACL:        acl,
	}
	videoMuxing720pOutput := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: stringToPtr("golang_test/video/720p"),
		ACL:        acl,
	}
	audioMuxingOutput := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: stringToPtr("golang_test/audio"),
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
	errorHandler(videoFMP4Muxing1080pResp.Status, err)

	videoFMP4Muxing720p := &models.FMP4Muxing{
		SegmentLength:   floatToPtr(4.0),
		SegmentNaming:   stringToPtr("seg_%number%.m4s"),
		InitSegmentName: stringToPtr("init.mp4"),
		Streams:         []models.StreamItem{videoMuxingStream720p},
		Outputs:         []models.Output{videoMuxing720pOutput},
	}
	videoFMP4Muxing720pResp, err := encodingS.AddFMP4Muxing(*encodingResp.Data.Result.ID, videoFMP4Muxing720p)
	errorHandler(videoFMP4Muxing720pResp.Status, err)

	audioFMP4Muxing := &models.FMP4Muxing{
		SegmentLength:   floatToPtr(4.0),
		SegmentNaming:   stringToPtr("seg_%number%.m4s"),
		InitSegmentName: stringToPtr("init.mp4"),
		Streams:         []models.StreamItem{audioMuxingStream},
		Outputs:         []models.Output{audioMuxingOutput},
	}
	audioFMP4MuxingResp, err := encodingS.AddFMP4Muxing(*encodingResp.Data.Result.ID, audioFMP4Muxing)
	errorHandler(audioFMP4MuxingResp.Status, err)

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

	startReq := &models.StartOptions{
		Scheduling: &models.EncodingScheduling{
			Priority: intToPtr(50),
		},
	}
	startResp, err := encodingS.StartWithOptions(*encodingResp.Data.Result.ID, startReq)
	errorHandler(startResp.Status, err)

	fmt.Printf("Start Response: %s\n", getAsJsonString(*startResp))

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
		fmt.Printf("Status: %s\n", getAsJsonString(*statusResp))

		status = *statusResp.Data.Result.Status
		if status == "ERROR" {
			fmt.Println("error in Encoding Status")
			fmt.Printf("Status %s\n", getAsJsonString(*statusResp))
			return
		}
	}

	manifestOutput := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: stringToPtr("golang_test/manifest"),
		ACL:        acl,
	}
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
		SegmentPath: stringToPtr("../video/1080p"),
	}
	fmp4Rep1080Resp, err := dashService.AddFMP4Representation(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *vasResp.Data.Result.ID, fmp4Rep1080)
	errorHandler(fmp4Rep1080Resp.Status, err)

	fmp4Rep720 := &models.FMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    videoFMP4Muxing720pResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: stringToPtr("../video/720p"),
	}
	fmp4Rep720Resp, err := dashService.AddFMP4Representation(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *vasResp.Data.Result.ID, fmp4Rep720)
	errorHandler(fmp4Rep720Resp.Status, err)

	fmp4RepAudio := &models.FMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    audioFMP4MuxingResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: stringToPtr("../audio"),
	}
	fmp4RepAudioResp, err := dashService.AddFMP4Representation(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *aasResp.Data.Result.ID, fmp4RepAudio)
	errorHandler(fmp4RepAudioResp.Status, err)

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
		fmt.Printf("DASH manifest creation status: %s\n", getAsJsonString(statusResp))
		status = *statusResp.Data.Result.Status
		if status == "ERROR" {
			fmt.Println("error in Manifest Status")
			fmt.Printf("%s\n", getAsJsonString(statusResp))
			return
		}
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
		fmt.Printf("HLS manifest creation status: %s\n", getAsJsonString(statusResp))
		status = *statusResp.Data.Result.Status
		if status == "ERROR" {
			fmt.Println("error in Manifest Status")
			fmt.Printf("%s\n", getAsJsonString(statusResp))
			return
		}
	}

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

func getAsJsonString(v interface{}) string {
	j, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(j)
}
