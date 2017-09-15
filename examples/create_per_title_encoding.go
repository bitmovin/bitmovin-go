package main

import (
	"fmt"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/bitmovintypes"
	"github.com/bitmovin/bitmovin-go/models"
	"github.com/bitmovin/bitmovin-go/services"
	"time"
)

const MEDIAN_BITRATE = 500000
const MAX_COMPLEXITY_FACTOR = 1.5
const MIN_COMPLEXITY_FACTOR = 0.5

type H264CodecConfigDefinition struct {
	height  *int64
	bitrate *int64
	fps     *float64
}

type EncodingConfig struct {
	codecConfigDef  *H264CodecConfigDefinition
	codecConfigResp *models.H264CodecConfigurationResponse
	streamResp      *models.StreamResponse
	fmp4MuxingResp  *models.FMP4MuxingResponse
	tsMuxingResp    *models.TSMuxingResponse
}

func main() {
	bitmovin := bitmovin.NewBitmovinDefaultTimeout("YOUR_BITMOVIN_API_KEY", "https://api.bitmovin.com/v1/")

	videoEncodingProfiles := []*H264CodecConfigDefinition{
		{height: intToPtr(180), bitrate: intToPtr(200), fps: nil},
		{height: intToPtr(180), bitrate: intToPtr(250), fps: nil},
		{height: intToPtr(180), bitrate: intToPtr(300), fps: nil},
		{height: intToPtr(270), bitrate: intToPtr(500), fps: nil},
		{height: intToPtr(360), bitrate: intToPtr(800), fps: nil},
		{height: intToPtr(360), bitrate: intToPtr(1000), fps: nil},
		{height: intToPtr(480), bitrate: intToPtr(1500), fps: nil},
		{height: intToPtr(720), bitrate: intToPtr(3000), fps: nil},
		{height: intToPtr(720), bitrate: intToPtr(4000), fps: nil},
		{height: intToPtr(1080), bitrate: intToPtr(6000), fps: nil},
		{height: intToPtr(1080), bitrate: intToPtr(7000), fps: nil},
		{height: intToPtr(1080), bitrate: intToPtr(10000), fps: nil},
	}

	inputFilePath := "/path/to/your/input/file.mkv"
	outputBasePath := "/your/output/path/"

	fmt.Println("Creating GCS Input")
	gcsIS := services.NewGCSInputService(bitmovin)
	gcsInput := &models.GCSInput{
		AccessKey:  stringToPtr("YOUR_ACCESS_KEY"),
		SecretKey:  stringToPtr("YOUR_SECRET_KEY"),
		BucketName: stringToPtr("YOUR_BUCKET_NAME"),
	}
	gcsInputResp, err := gcsIS.Create(gcsInput)
	errorHandler(gcsInputResp.Status, err)
	fmt.Println("Created GCS Input!")

	fmt.Println("Creating GCS Output")
	gcsOS := services.NewGCSOutputService(bitmovin)
	gcsOutput := &models.GCSOutput{
		AccessKey:  stringToPtr("YOUR_ACCESS_KEY"),
		SecretKey:  stringToPtr("YOUR_SECRET_KEY"),
		BucketName: stringToPtr("YOUR_BUCKET_NAME"),
	}
	outputResponse, err := gcsOS.Create(gcsOutput)
	errorHandler(outputResponse.Status, err)
	fmt.Println("Created GCS Output!")

	fmt.Println("Creating Analysis Encoding")
	encodingS := services.NewEncodingService(bitmovin)
	analysisEncoding := models.Encoding{
		Name:        stringToPtr("Per title analysis encoding"),
		CloudRegion: bitmovintypes.CloudRegionGoogleEuropeWest1,
	}
	analysisEncodingResp, err := encodingS.Create(&analysisEncoding)
	errorHandler(analysisEncodingResp.Status, err)
	fmt.Println("Created Analysis Encoding!")

	fmt.Println("Creating Codec Configuration")
	h264S := services.NewH264CodecConfigurationService(bitmovin)
	analysisH264CodecConfig := models.H264CodecConfiguration{
		Name:    stringToPtr("H264 Per Title Analysis Configuration"),
		CRF:     floatToPtr(23.0),
		Height:  intToPtr(360),
		Profile: bitmovintypes.H264ProfileMain,
	}
	analysisH264CodecConfigResp, err := h264S.Create(&analysisH264CodecConfig)
	errorHandler(analysisH264CodecConfigResp.Status, err)
	fmt.Println("Created Codec Configuration!")

	fmt.Println("Creating Streams...")
	inputStream := models.InputStream{
		InputID:       gcsInputResp.Data.Result.ID,
		InputPath:     stringToPtr(inputFilePath),
		SelectionMode: bitmovintypes.SelectionModeAuto,
	}
	vis := []models.InputStream{inputStream}
	analysisVideoStream := &models.Stream{
		CodecConfigurationID: analysisH264CodecConfigResp.Data.Result.ID,
		InputStreams:         vis,
		Name:                 stringToPtr("Per Title Analysis Stream"),
	}
	fmt.Println("Created Streams!")

	analysisVideoStreamResp, err := encodingS.AddStream(*analysisEncodingResp.Data.Result.ID, analysisVideoStream)
	errorHandler(analysisVideoStreamResp.Status, err)
	fmt.Println("Created Stream!")

	fmt.Println("Creating Muxing...")
	analysisMuxingStream := models.StreamItem{
		StreamID: analysisVideoStreamResp.Data.Result.ID,
	}

	aclEntry := models.ACLItem{
		Permission: bitmovintypes.ACLPermissionPublicRead,
	}
	analysisMuxingOutput := models.Output{
		OutputID:   outputResponse.Data.Result.ID,
		OutputPath: stringToPtr(outputBasePath + "quality_analysis/"),
		ACL:        []models.ACLItem{aclEntry},
	}

	analysisMuxing := models.FMP4Muxing{
		SegmentLength:   floatToPtr(4.0),
		SegmentNaming:   stringToPtr("seg_%number%.m4s"),
		InitSegmentName: stringToPtr("init.mp4"),
		Streams:         []models.StreamItem{analysisMuxingStream},
		Outputs:         []models.Output{analysisMuxingOutput},
		Name:            stringToPtr("Per Title Analysis Muxing"),
	}
	analysisMuxingResp, err := encodingS.AddFMP4Muxing(*analysisEncodingResp.Data.Result.ID, &analysisMuxing)
	errorHandler(analysisMuxingResp.Status, err)
	fmt.Println("Created Muxing!")

	fmt.Println("Starting Analysis Encoding...")
	startResp, err := encodingS.Start(*analysisEncodingResp.Data.Result.ID)
	errorHandler(startResp.Status, err)
	fmt.Println("Started Analysis Encoding!")

	fmt.Println("Waiting for Analysis Encoding to be finished...")
	waitForEncodingToBeFinished(analysisEncodingResp, encodingS)
	fmt.Println("Analysis Encoding finished!")

	analysisMuxingResp, err = encodingS.RetrieveFMP4Muxing(*analysisEncodingResp.Data.Result.ID, *analysisMuxingResp.Data.Result.ID)
	analysisMuxing = analysisMuxingResp.Data.Result
	complexityFactor := float64(*analysisMuxing.AvgBitrate) / float64(MEDIAN_BITRATE)

	if complexityFactor > MAX_COMPLEXITY_FACTOR {
		complexityFactor = MAX_COMPLEXITY_FACTOR
	}
	if complexityFactor < MIN_COMPLEXITY_FACTOR {
		complexityFactor = MIN_COMPLEXITY_FACTOR
	}
	fmt.Printf("Used values for calculation -> avgBitrate %d, medianBitrate %d\n", *analysisMuxing.AvgBitrate, MEDIAN_BITRATE)
	fmt.Printf("Got complexity factor of %f\n", complexityFactor)

	fmt.Println("Creating the Encoding...")
	realEncoding := models.Encoding{
		Name:        stringToPtr("Golang - Per Title Encoding"),
		CloudRegion: bitmovintypes.CloudRegionGoogleEuropeWest1,
	}
	realEncodingResp, err := encodingS.Create(&realEncoding)
	errorHandler(realEncodingResp.Status, err)
	fmt.Println("Created Encoding!")

	updateVideoEncodingProfilesWithComplexityFactor(videoEncodingProfiles, complexityFactor)
	encodingConfigs := createEncodingConfigs(videoEncodingProfiles)
	createAndAddH264CodecConfigurationsToEncodingConfigs(encodingConfigs, h264S)
	createAndAddVideoStreamsToEncodingConfigs(encodingConfigs, inputStream, encodingS, *realEncodingResp.Data.Result.ID)
	createAndAddFmp4MuxingsToEncodingConfigs(encodingConfigs, *outputResponse, outputBasePath, encodingS, *realEncodingResp.Data.Result.ID)
	createAndAddTsMuxingsToEncodingConfigs(encodingConfigs, *outputResponse, outputBasePath, encodingS, *realEncodingResp.Data.Result.ID)

	fmt.Println("Creating AAC Configuration...")
	aacS := services.NewAACCodecConfigurationService(bitmovin)
	audioCodecConfiguration := models.AACCodecConfiguration{
		Name:         stringToPtr("AAC Codec Configuration"),
		Bitrate:      intToPtr(128000),
		SamplingRate: floatToPtr(48000),
	}
	audioCodecConfigurationResp, err := aacS.Create(&audioCodecConfiguration)
	errorHandler(audioCodecConfigurationResp.Status, err)
	fmt.Println("Created AAC Configuration!")

	fmt.Println("Creating Audio Stream...")
	ais := []models.InputStream{inputStream}
	audioConfigurationStream := &models.Stream{
		CodecConfigurationID: audioCodecConfigurationResp.Data.Result.ID,
		InputStreams:         ais,
		Name:                 stringToPtr("Audio Stream"),
	}
	audioStreamResp, err := encodingS.AddStream(*realEncodingResp.Data.Result.ID, audioConfigurationStream)
	errorHandler(audioStreamResp.Status, err)
	fmt.Println("Created Audio Stream!")

	fmt.Println("Creating Audio FMP4 Muxing...")
	audioMuxingStream := models.StreamItem{
		StreamID: audioStreamResp.Data.Result.ID,
	}

	audioMuxingOutput := models.Output{
		OutputID:   outputResponse.Data.Result.ID,
		OutputPath: stringToPtr(outputBasePath + "dash/audio/"),
		ACL:        []models.ACLItem{aclEntry},
	}

	audioFmp4Muxing := models.FMP4Muxing{
		SegmentLength:   floatToPtr(4.0),
		SegmentNaming:   stringToPtr("seg_%number%.m4s"),
		InitSegmentName: stringToPtr("init.mp4"),
		Streams:         []models.StreamItem{audioMuxingStream},
		Outputs:         []models.Output{audioMuxingOutput},
		Name:            stringToPtr("Audio FMP4 Muxing"),
	}
	audioFmp4MuxingResp, err := encodingS.AddFMP4Muxing(*realEncodingResp.Data.Result.ID, &audioFmp4Muxing)
	errorHandler(audioFmp4MuxingResp.Status, err)
	fmt.Println("Created Audio FMP4 Muxing!")

	fmt.Println("Creating Audio TS Muxing...")
	audioTsMuxingOutput := models.Output{
		OutputID:   outputResponse.Data.Result.ID,
		OutputPath: stringToPtr(outputBasePath + "hls/audio/"),
		ACL:        []models.ACLItem{aclEntry},
	}

	audioTsMuxing := models.TSMuxing{
		SegmentLength: floatToPtr(4.0),
		SegmentNaming: stringToPtr("seg_%number%.ts"),
		Streams:       []models.StreamItem{audioMuxingStream},
		Outputs:       []models.Output{audioTsMuxingOutput},
		Name:          stringToPtr("Audio Muxing"),
	}
	audioTsMuxingResp, err := encodingS.AddTSMuxing(*realEncodingResp.Data.Result.ID, &audioTsMuxing)
	errorHandler(audioTsMuxingResp.Status, err)
	fmt.Println("Created Audio TS Muxing!")

	fmt.Println("Starting Encoding...")
	encStartResp, err := encodingS.Start(*realEncodingResp.Data.Result.ID)
	errorHandler(encStartResp.Status, err)
	fmt.Println("Started Encoding!")

	fmt.Println("Waiting for Actual Encoding to be finished...")
	waitForEncodingToBeFinished(realEncodingResp, encodingS)
	fmt.Println("Encoding finished!")

	fmt.Println("Creating DASH Manifest...")
	createDashManifest(
		encodingConfigs,
		audioFmp4MuxingResp,
		*realEncodingResp.Data.Result.ID,
		*outputResponse.Data.Result.ID,
		outputBasePath,
		[]models.ACLItem{aclEntry},
		bitmovin,
	)
	fmt.Println("Created DASH Manifest!")

	fmt.Println("Creating HLS Manifests...")
	createHlsManifest(
		encodingConfigs,
		*audioTsMuxingResp.Data.Result.ID,
		*audioStreamResp.Data.Result.ID,
		*realEncodingResp.Data.Result.ID,
		*outputResponse.Data.Result.ID,
		outputBasePath,
		[]models.ACLItem{aclEntry},
		bitmovin,
	)
	fmt.Println("Created HLS Manifests!")
}

func createDashManifest(
	encodingConfigs []*EncodingConfig,
	audioMuxingResp *models.FMP4MuxingResponse,
	encodingId string,
	outputId string,
	outputBasePath string,
	acl []models.ACLItem,
	bitmovin *bitmovin.Bitmovin,
) {

	manifestOutput := models.Output{
		OutputID:   stringToPtr(outputId),
		OutputPath: stringToPtr(outputBasePath + "/manifest"),
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

	/*
		AUDIO
	*/
	aas := &models.AudioAdaptationSet{
		Language: stringToPtr("en"),
	}
	aasResp, err := dashService.AddAudioAdaptationSet(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, aas)
	errorHandler(aasResp.Status, err)

	fmp4RepAudio := &models.FMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    audioMuxingResp.Data.Result.ID,
		EncodingID:  stringToPtr(encodingId),
		SegmentPath: stringToPtr("../dash/audio"),
	}
	fmp4RepAudioResp, err := dashService.AddFMP4Representation(
		*dashManifestResp.Data.Result.ID,
		*periodResp.Data.Result.ID,
		*aasResp.Data.Result.ID,
		fmp4RepAudio)
	errorHandler(fmp4RepAudioResp.Status, err)

	/*
		VIDEO
	*/
	vas := &models.VideoAdaptationSet{}
	vasResp, err := dashService.AddVideoAdaptationSet(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, vas)
	errorHandler(vasResp.Status, err)

	for _, encodingConfig := range encodingConfigs {
		fmp4Represetation := &models.FMP4Representation{
			Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
			MuxingID:    encodingConfig.fmp4MuxingResp.Data.Result.ID,
			EncodingID:  stringToPtr(encodingId),
			SegmentPath: stringToPtr(fmt.Sprintf("../dash/video/%dp_%dk", *encodingConfig.codecConfigDef.height, *encodingConfig.codecConfigDef.bitrate)),
		}
		fmp4RepresetationResp, err := dashService.AddFMP4Representation(
			*dashManifestResp.Data.Result.ID,
			*periodResp.Data.Result.ID,
			*vasResp.Data.Result.ID,
			fmp4Represetation,
		)
		errorHandler(fmp4RepresetationResp.Status, err)
	}

	fmt.Printf("Starting DASH manifest generation with manifest id %s...\n", *dashManifestResp.Data.Result.ID)

	startResp, err := dashService.Start(*dashManifestResp.Data.Result.ID)
	errorHandler(startResp.Status, err)

	status := ""
	for status != "FINISHED" {
		time.Sleep(5 * time.Second)
		statusResp, err := dashService.RetrieveStatus(*dashManifestResp.Data.Result.ID)
		if err != nil {
			fmt.Println("error in Manifest Status")
			fmt.Println(err)
			return
		}
		status = *statusResp.Data.Result.Status
		fmt.Printf("STATUS: %s\n", status)
		if status == "ERROR" {
			fmt.Println("error in Manifest Status")
			fmt.Printf("STATUS: %s\n", status)
			return
		}
	}
	fmt.Println("DASH manifest created successfully!")
}

func createHlsManifest(
	encodingConfigs []*EncodingConfig,
	audioMuxingId string,
	audioStreamId string,
	encodingId string,
	outputId string,
	outputBasePath string,
	acl []models.ACLItem,
	bitmovin *bitmovin.Bitmovin,
) {

	manifestOutput := models.Output{
		OutputID:   stringToPtr(outputId),
		OutputPath: stringToPtr(outputBasePath + "/manifest"),
		ACL:        acl,
	}

	hlsService := services.NewHLSManifestService(bitmovin)
	hlsManifest := &models.HLSManifest{
		ManifestName: stringToPtr("your_manifest_name.m3u8"),
		Outputs:      []models.Output{manifestOutput},
	}
	hlsManifestResp, err := hlsService.Create(hlsManifest)
	errorHandler(hlsManifestResp.Status, err)

	/*
		AUDIO
	*/
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
		SegmentPath:     stringToPtr("../hls/audio"),
		EncodingID:      stringToPtr(encodingId),
		StreamID:        stringToPtr(audioStreamId),
		MuxingID:        stringToPtr(audioMuxingId),
	}
	audioMediaInfoResp, err := hlsService.AddMediaInfo(*hlsManifestResp.Data.Result.ID, audioMediaInfo)
	errorHandler(audioMediaInfoResp.Status, err)

	/*
		VIDEO
	*/
	for _, encodingConfig := range encodingConfigs {
		videoStreamInfo := &models.StreamInfo{
			Audio:       stringToPtr("audio_group"),
			SegmentPath: stringToPtr(fmt.Sprintf("../hls/video/%dp_%dk", *encodingConfig.codecConfigDef.height, *encodingConfig.codecConfigDef.bitrate)),
			URI:         stringToPtr(fmt.Sprintf("video_%dp_%dk.m3u8", *encodingConfig.codecConfigDef.height, *encodingConfig.codecConfigDef.bitrate)),
			EncodingID:  stringToPtr(encodingId),
			StreamID:    encodingConfig.streamResp.Data.Result.ID,
			MuxingID:    encodingConfig.tsMuxingResp.Data.Result.ID,
		}
		videoStreamInfoResp, err := hlsService.AddStreamInfo(*hlsManifestResp.Data.Result.ID, videoStreamInfo)
		errorHandler(videoStreamInfoResp.Status, err)
	}

	fmt.Printf("Starting HLS manifest generation with manifest id %s...\n", *hlsManifestResp.Data.Result.ID)

	startResp, err := hlsService.Start(*hlsManifestResp.Data.Result.ID)
	errorHandler(startResp.Status, err)

	status := ""
	for status != "FINISHED" {
		time.Sleep(5 * time.Second)
		statusResp, err := hlsService.RetrieveStatus(*hlsManifestResp.Data.Result.ID)
		if err != nil {
			fmt.Println("error in Manifest Status")
			fmt.Println(err)
			return
		}
		status = *statusResp.Data.Result.Status
		fmt.Printf("STATUS: %s\n", status)
		if status == "ERROR" {
			fmt.Println("error in Manifest Status")
			fmt.Printf("STATUS: %s\n", status)
			return
		}
	}
	fmt.Println("HLS manifest created successfully!")
}

func createEncodingConfigs(
	codecConfigDefinitions []*H264CodecConfigDefinition,
) []*EncodingConfig {
	var encodingConfigs []*EncodingConfig
	for _, codecConfigDefinition := range codecConfigDefinitions {
		encodingConfigs = append(encodingConfigs, &EncodingConfig{codecConfigDef: codecConfigDefinition})
	}
	return encodingConfigs
}

func updateVideoEncodingProfilesWithComplexityFactor(
	codecConfigDefinitions []*H264CodecConfigDefinition,
	complexityFactor float64,
) {
	for _, codecConfigDefinition := range codecConfigDefinitions {
		fmt.Printf("Bitrate before: %d\n", *codecConfigDefinition.bitrate)
		codecConfigDefinition.bitrate = intToPtr(int64(float64(*codecConfigDefinition.bitrate) * complexityFactor))
		fmt.Printf("Bitrate after: %d\n", *codecConfigDefinition.bitrate)
	}
}

func createAndAddTsMuxingsToEncodingConfigs(
	encodingConfigs []*EncodingConfig,
	outputResp models.GCSOutputResponse,
	outputBasePath string,
	encodingS *services.EncodingService,
	encodingId string,
) {
	for _, encodingProfile := range encodingConfigs {
		fmt.Println("Creating Muxing...")
		tsMuxingStream := models.StreamItem{
			StreamID: encodingProfile.streamResp.Data.Result.ID,
		}

		aclEntry := models.ACLItem{
			Permission: bitmovintypes.ACLPermissionPublicRead,
		}
		muxingOutput := models.Output{
			OutputID:   outputResp.Data.Result.ID,
			OutputPath: stringToPtr(fmt.Sprintf(outputBasePath+"hls/video/%dp_%dk", *encodingProfile.codecConfigDef.height, *encodingProfile.codecConfigDef.bitrate)),
			ACL:        []models.ACLItem{aclEntry},
		}

		muxing := models.TSMuxing{
			SegmentLength: floatToPtr(4.0),
			SegmentNaming: stringToPtr("seg_%number%.ts"),
			Streams:       []models.StreamItem{tsMuxingStream},
			Outputs:       []models.Output{muxingOutput},
			Name:          stringToPtr(fmt.Sprintf("TS Muxing %dp_%dk", *encodingProfile.codecConfigDef.height, *encodingProfile.codecConfigDef.bitrate)),
		}
		muxingResp, err := encodingS.AddTSMuxing(encodingId, &muxing)
		errorHandler(muxingResp.Status, err)
		encodingProfile.tsMuxingResp = muxingResp
		fmt.Println("Created Muxing!")
	}
}

func createAndAddFmp4MuxingsToEncodingConfigs(
	encodingConfigs []*EncodingConfig,
	outputResp models.GCSOutputResponse,
	outputBasePath string,
	encodingS *services.EncodingService,
	encodingId string,
) {
	for _, encodingProfile := range encodingConfigs {
		fmt.Printf("Creating FMP4 Muxing %dp_%dk\n", *encodingProfile.codecConfigDef.height, *encodingProfile.codecConfigDef.bitrate)
		fmp4MuxingStream := models.StreamItem{
			StreamID: encodingProfile.streamResp.Data.Result.ID,
		}

		aclEntry := models.ACLItem{
			Permission: bitmovintypes.ACLPermissionPublicRead,
		}
		muxingOutput := models.Output{
			OutputID:   outputResp.Data.Result.ID,
			OutputPath: stringToPtr(fmt.Sprintf(outputBasePath+"dash/video/%dp_%dk", *encodingProfile.codecConfigDef.height, *encodingProfile.codecConfigDef.bitrate)),
			ACL:        []models.ACLItem{aclEntry},
		}

		muxing := models.FMP4Muxing{
			SegmentLength:   floatToPtr(4.0),
			SegmentNaming:   stringToPtr("seg_%number%.m4s"),
			InitSegmentName: stringToPtr("init.mp4"),
			Streams:         []models.StreamItem{fmp4MuxingStream},
			Outputs:         []models.Output{muxingOutput},
			Name:            stringToPtr(fmt.Sprintf("FMP4 Muxing %dp_%dk", *encodingProfile.codecConfigDef.height, *encodingProfile.codecConfigDef.bitrate)),
		}
		muxingResp, err := encodingS.AddFMP4Muxing(encodingId, &muxing)
		errorHandler(muxingResp.Status, err)
		encodingProfile.fmp4MuxingResp = muxingResp
		fmt.Println("Created Muxing!")
	}
}

func createAndAddVideoStreamsToEncodingConfigs(
	encodingConfigs []*EncodingConfig,
	videoInputStream models.InputStream,
	encodingS *services.EncodingService,
	encodingId string,
) {
	for _, encodingConfig := range encodingConfigs {
		fmt.Printf("Creating Stream %dp_%dk\n", *encodingConfig.codecConfigDef.height, *encodingConfig.codecConfigDef.bitrate)
		vis := []models.InputStream{videoInputStream}
		videoStream := &models.Stream{
			CodecConfigurationID: encodingConfig.codecConfigResp.Data.Result.ID,
			InputStreams:         vis,
			Name:                 stringToPtr(fmt.Sprintf("Stream %dp_%dk", *encodingConfig.codecConfigDef.height, *encodingConfig.codecConfigDef.bitrate)),
		}
		videoStreamResp, err := encodingS.AddStream(encodingId, videoStream)
		errorHandler(videoStreamResp.Status, err)
		encodingConfig.streamResp = videoStreamResp
		fmt.Println("Created Stream!")
	}
}

func createAndAddH264CodecConfigurationsToEncodingConfigs(
	encodingConfigs []*EncodingConfig,
	h264S *services.H264CodecConfigurationService,
) {
	for _, encodingConfig := range encodingConfigs {
		fmt.Printf("Creating codec config (bitrate = %d, height = %d)\n", *encodingConfig.codecConfigDef.bitrate, *encodingConfig.codecConfigDef.height)
		h264CodecConfig := models.H264CodecConfiguration{
			Name:      stringToPtr(fmt.Sprintf("H264 Configuration %dp %dk", *encodingConfig.codecConfigDef.height, *encodingConfig.codecConfigDef.bitrate)),
			Bitrate:   intToPtr(*encodingConfig.codecConfigDef.bitrate * 1000),
			Height:    encodingConfig.codecConfigDef.height,
			FrameRate: encodingConfig.codecConfigDef.fps,
			Profile:   bitmovintypes.H264ProfileHigh,
		}
		codecConfigResp, err := h264S.Create(&h264CodecConfig)
		errorHandler(codecConfigResp.Status, err)
		fmt.Println("Created codec config!")
		encodingConfig.codecConfigResp = codecConfigResp
	}
}

func waitForEncodingToBeFinished(
	encodingResp *models.EncodingResponse,
	encodingS *services.EncodingService,
) {
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
