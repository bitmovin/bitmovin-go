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
	bitmovin := bitmovin.NewBitmovin("YOUR_API_KEY", "https://api.bitmovin.com/v1/", 5)

	// Creating the RTMP Input
	rtmpIS := services.NewRTMPInputService(bitmovin)
	rtmpInputListResp, err := rtmpIS.List(0, 10)
	errorHandler(err)

	if len(rtmpInputListResp.Data.Result.Items) < 1 {
		fmt.Println("No RTMP inputs on account!")
		return
	}

	gcsOutputService := services.NewGCSOutputService(bitmovin)
	gcsOutput := &models.GCSOutput{
		AccessKey:   stringToPtr("YOUR_ACCESS_KEY"),
		SecretKey:   stringToPtr("YOUR_SECRET_KEY"),
		BucketName:  stringToPtr("YOUR_BUCKET_NAME"),
		CloudRegion: bitmovintypes.GoogleCloudRegionEuropeWest1,
	}

	t := time.Now()
	outputBasePath := "golang_live_drm_test_" + t.Format("2006-01-02-15-04-05")

	gcsOutputResponse, err := gcsOutputService.Create(gcsOutput)
	errorHandler(err)

	encodingS := services.NewEncodingService(bitmovin)
	encoding := &models.Encoding{
		Name:        stringToPtr("Example Golang Live Encoding with DRM"),
		CloudRegion: bitmovintypes.CloudRegionGoogleUSEast1,
	}
	encodingResp, err := encodingS.Create(encoding)
	errorHandler(err)

	encodingID := *encodingResp.Data.Result.ID

	/*
		H264 Codec Configurations
	*/
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

	/*
		Input Streams
	*/
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

	/*
		DASH FMP4 Muxings
	*/
	fmt.Println("Creating DASH FMP4 Muxings...")
	videoMuxing1080p := &models.FMP4Muxing{
		SegmentLength:   floatToPtr(4.0),
		SegmentNaming:   stringToPtr("seg_%number%.m4s"),
		InitSegmentName: stringToPtr("init.mp4"),
		Streams:         []models.StreamItem{videoMuxingStream1080p},
	}
	videoMuxing1080pResp, err := encodingS.AddFMP4Muxing(encodingID, videoMuxing1080p)
	errorHandler(err)

	videoMuxing720p := &models.FMP4Muxing{
		SegmentLength:   floatToPtr(4.0),
		SegmentNaming:   stringToPtr("seg_%number%.m4s"),
		InitSegmentName: stringToPtr("init.mp4"),
		Streams:         []models.StreamItem{videoMuxingStream720p},
	}
	videoMuxing720pResp, err := encodingS.AddFMP4Muxing(encodingID, videoMuxing720p)
	errorHandler(err)

	audioMuxing := &models.FMP4Muxing{
		SegmentLength:   floatToPtr(4.0),
		SegmentNaming:   stringToPtr("seg_%number%.m4s"),
		InitSegmentName: stringToPtr("init.mp4"),
		Streams:         []models.StreamItem{audioMuxingStream},
	}
	audioMuxingResp, err := encodingS.AddFMP4Muxing(encodingID, audioMuxing)
	errorHandler(err)
	fmt.Println("Successfully created DASH FMP4 Muxings!")

	/*
		HLS TS Muxings
	*/
	fmt.Println("Creating HLS TS Muxings...")
	videoTsMuxing1080p := &models.TSMuxing{
		SegmentLength: floatToPtr(4.0),
		Streams:       []models.StreamItem{videoMuxingStream1080p},
	}

	videoTsMuxing1080Resp, err := encodingS.AddTSMuxing(encodingID, videoTsMuxing1080p)
	errorHandler(err)

	videoTsMuxing720p := &models.TSMuxing{
		SegmentLength: floatToPtr(4.0),
		Streams:       []models.StreamItem{videoMuxingStream720p},
	}
	videoTsMuxing720Resp, err := encodingS.AddTSMuxing(encodingID, videoTsMuxing720p)
	errorHandler(err)

	audioTsMuxing := &models.TSMuxing{
		SegmentLength: floatToPtr(4.0),
		Streams:       []models.StreamItem{audioMuxingStream},
	}
	audioTsMuxingResp, err := encodingS.AddTSMuxing(encodingID, audioTsMuxing)
	errorHandler(err)

	/*
		Outputs for DRM FMP4 Muxings
	*/
	videoFMP4Muxing1080pOutput := models.Output{
		OutputID:   gcsOutputResponse.Data.Result.ID,
		OutputPath: stringToPtr(outputBasePath + "/video/dash/1080p"),
		ACL:        acl,
	}
	videoFMP4Muxing720pOutput := models.Output{
		OutputID:   gcsOutputResponse.Data.Result.ID,
		OutputPath: stringToPtr(outputBasePath + "/video/dash/720p"),
		ACL:        acl,
	}
	audioFMP4MuxingOutput := models.Output{
		OutputID:   gcsOutputResponse.Data.Result.ID,
		OutputPath: stringToPtr(outputBasePath + "/audio/dash"),
		ACL:        acl,
	}
	fmt.Println("Successfully created HLS TS Muxings!")

	/*
		Widevine and Playready CENC DRM
	*/
	fmt.Println("Creating Widevine and PlayReady CENC DRMs...")
	widevine := models.WidevineCencDrm{
		PSSH: stringToPtr("WIDEVINE_PSSH"),
	}

	playready := models.PlayReadyCencDrm{
		LaURL: stringToPtr("PLAYREADY_LA_URL"),
	}

	cencDrm := models.CencDrm{
		Key:       stringToPtr("YOUR_CENC_DRM_KEY"),
		KID:       stringToPtr("YOUR_CENC_DRM_KID"),
		Name:      stringToPtr("My CENC DRM"),
		Widevine:  widevine,
		PlayReady: playready,
	}

	drmService := services.NewDrmService(bitmovin)

	cencDrm.Outputs = []models.Output{videoFMP4Muxing720pOutput}
	cencResp720, err := drmService.CreateFmp4Drm(encodingID, *videoMuxing720pResp.Data.Result.ID, cencDrm)
	if err != nil {
		fmt.Println(err)
		return
	}
	cencDrm720Response := cencResp720.(*models.CencDrmResponse)
	errorHandler(err)

	cencDrm.Outputs = []models.Output{videoFMP4Muxing1080pOutput}
	cencResp1080, err := drmService.CreateFmp4Drm(encodingID, *videoMuxing1080pResp.Data.Result.ID, cencDrm)
	if err != nil {
		fmt.Println(err)
		return
	}
	cencDrm1080Response := cencResp1080.(*models.CencDrmResponse)
	errorHandler(err)

	cencDrm.Outputs = []models.Output{audioFMP4MuxingOutput}
	cencRespAudio, err := drmService.CreateFmp4Drm(encodingID, *audioMuxingResp.Data.Result.ID, cencDrm)
	if err != nil {
		fmt.Println(err)
		return
	}
	cencDrmAudioResponse := cencRespAudio.(*models.CencDrmResponse)
	errorHandler(err)

	fmt.Println("Successfully created Widevine and PlayReady CENC DRMs!")

	/*
		FairPlay DRM
	*/
	videoTSMuxing1080pOutput := models.Output{
		OutputID:   gcsOutputResponse.Data.Result.ID,
		OutputPath: stringToPtr(outputBasePath + "/video/hls/1080p"),
		ACL:        acl,
	}
	videoTSMuxing720pOutput := models.Output{
		OutputID:   gcsOutputResponse.Data.Result.ID,
		OutputPath: stringToPtr(outputBasePath + "/video/hls/720p"),
		ACL:        acl,
	}
	audioTSMuxingOutput := models.Output{
		OutputID:   gcsOutputResponse.Data.Result.ID,
		OutputPath: stringToPtr(outputBasePath + "/audio/hls"),
		ACL:        acl,
	}

	fmt.Println("Creating Fairplay DRMs for HLS...")
	fairPlayDrm := models.FairPlayDrm{
		IV:   stringToPtr("YOUR_FAIRPLAY_IV"),
		URI:  stringToPtr("YOUR_FAIRPLAY_URI"),
		Key:  stringToPtr("YOUR_FAIRPLAY_KEY"),
		Name: stringToPtr("My Fairplay DRM"),
	}

	fairPlayDrm.Outputs = []models.Output{videoTSMuxing720pOutput}
	tsDrm720pResp, err := drmService.CreateTsDrm(encodingID, *videoTsMuxing720Resp.Data.Result.ID, fairPlayDrm)
	if err != nil {
		fmt.Println(err)
		return
	}
	tsDrm720p := tsDrm720pResp.(*models.FairPlayDrmResponse)
	errorHandler(err)

	fairPlayDrm.Outputs = []models.Output{videoTSMuxing1080pOutput}
	tsDrm1080pResp, err := drmService.CreateTsDrm(encodingID, *videoTsMuxing1080Resp.Data.Result.ID, fairPlayDrm)
	if err != nil {
		fmt.Println(err)
		return
	}
	tsDrm1080p := tsDrm1080pResp.(*models.FairPlayDrmResponse)
	errorHandler(err)

	fairPlayDrm.Outputs = []models.Output{audioTSMuxingOutput}
	tsDrmAudioResp, err := drmService.CreateTsDrm(encodingID, *audioTsMuxingResp.Data.Result.ID, fairPlayDrm)
	if err != nil {
		fmt.Println(err)
		return
	}
	tsDrmAudio := tsDrmAudioResp.(*models.FairPlayDrmResponse)
	errorHandler(err)

	fmt.Println("Successfully created Fairplay DRMs for HLS!")

	/*
		DASH Manifest creation
	*/
	fmt.Println("Creating Live DASH Manifest...")

	manifestOutput := models.Output{
		OutputID:   gcsOutputResponse.Data.Result.ID,
		OutputPath: stringToPtr(outputBasePath + "/manifest/"),
		ACL:        acl,
	}

	dashManifest := &models.DashManifest{
		Name:         stringToPtr("your_manifest_name.mpd"),
		ManifestName: stringToPtr("your_manifest_name.mpd"),
		Outputs:      []models.Output{manifestOutput},
	}
	dashService := services.NewDashManifestService(bitmovin)
	dashManifestResp, err := dashService.Create(dashManifest)
	errorHandler(err)

	period := &models.Period{}
	periodResp, err := dashService.AddPeriod(*dashManifestResp.Data.Result.ID, period)
	errorHandler(err)

	/*
		Video AdaptationSet
	*/
	vas := &models.VideoAdaptationSet{}
	vasResp, err := dashService.AddVideoAdaptationSet(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, vas)
	errorHandler(err)

	fmt.Println("Adding content protection to video adaptation set...")
	vcp := &models.AdaptationSetContentProtection{
		DrmId:      cencDrm1080Response.Data.Result.ID,
		EncodingId: stringToPtr(encodingID),
		MuxingId:   videoMuxing1080pResp.Data.Result.ID,
	}
	_, err = dashService.AddContentProtectionToAdaptationSet(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *vasResp.Data.Result.ID, vcp)
	errorHandler(err)
	fmt.Println("Successfully added content protection to video adaptation set!")

	/*
		Audio AdaptationSet
	*/
	aas := &models.AudioAdaptationSet{
		Language: stringToPtr("en"),
	}
	aasResp, err := dashService.AddAudioAdaptationSet(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, aas)
	errorHandler(err)

	fmt.Println("Adding content protection to audio adaptation set...")
	acp := &models.AdaptationSetContentProtection{
		DrmId:      cencDrmAudioResponse.Data.Result.ID,
		EncodingId: stringToPtr(encodingID),
		MuxingId:   audioMuxingResp.Data.Result.ID,
	}
	_, err = dashService.AddContentProtectionToAdaptationSet(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *aasResp.Data.Result.ID, acp)
	errorHandler(err)
	fmt.Println("Successfully added content protection to audio adaptation set!")

	/*
		DRM Representations
	*/
	fmp4Rep1080 := &models.DrmFMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    videoMuxing1080pResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: stringToPtr("../video/dash/1080p"),
		DrmID:       cencDrm1080Response.Data.Result.ID,
	}
	_, err = dashService.AddDrmFMP4Representation(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *vasResp.Data.Result.ID, fmp4Rep1080)
	errorHandler(err)

	fmp4Rep720 := &models.DrmFMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    videoMuxing720pResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: stringToPtr("../video/dash/720p"),
		DrmID:       cencDrm720Response.Data.Result.ID,
	}
	_, err = dashService.AddDrmFMP4Representation(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *vasResp.Data.Result.ID, fmp4Rep720)
	errorHandler(err)

	fmp4RepAudio := &models.DrmFMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    audioMuxingResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: stringToPtr("../audio/dash"),
		DrmID:       cencDrmAudioResponse.Data.Result.ID,
	}
	_, err = dashService.AddDrmFMP4Representation(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *aasResp.Data.Result.ID, fmp4RepAudio)
	errorHandler(err)

	liveDashManifest := models.LiveDashManifest{
		ManifestID:     dashManifestResp.Data.Result.ID,
		LiveEdgeOffset: floatToPtr(45.0),
	}
	fmt.Println("Successfully created Live DASH Manifest!")

	/*
		Live HLS manifest
	*/
	fmt.Println("Creating Live HLS Manifest...")

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
		SegmentPath:     stringToPtr("../audio/hls"),
		EncodingID:      encodingResp.Data.Result.ID,
		StreamID:        aacStreamResp.Data.Result.ID,
		MuxingID:        audioTsMuxingResp.Data.Result.ID,
		DRMID:           tsDrmAudio.Data.Result.ID,
	}
	_, err = hlsService.AddMediaInfo(*hlsManifestResp.Data.Result.ID, audioMediaInfo)
	errorHandler(err)

	video1080pStreamInfo := &models.StreamInfo{
		Audio:       stringToPtr("audio_group"),
		SegmentPath: stringToPtr("../video/hls/1080p"),
		URI:         stringToPtr("video_hi.m3u8"),
		EncodingID:  encodingResp.Data.Result.ID,
		StreamID:    videoStream1080pResp.Data.Result.ID,
		MuxingID:    videoTsMuxing1080Resp.Data.Result.ID,
		DRMID:       tsDrm1080p.Data.Result.ID,
	}
	_, err = hlsService.AddStreamInfo(*hlsManifestResp.Data.Result.ID, video1080pStreamInfo)
	errorHandler(err)

	video720pStreamInfo := &models.StreamInfo{
		Audio:       stringToPtr("audio_group"),
		SegmentPath: stringToPtr("../video/hls/720p"),
		URI:         stringToPtr("video_lo.m3u8"),
		EncodingID:  encodingResp.Data.Result.ID,
		StreamID:    videoStream720pResp.Data.Result.ID,
		MuxingID:    videoTsMuxing720Resp.Data.Result.ID,
		DRMID:       tsDrm720p.Data.Result.ID,
	}
	_, err = hlsService.AddStreamInfo(*hlsManifestResp.Data.Result.ID, video720pStreamInfo)
	errorHandler(err)

	liveHlsManifest := models.LiveHLSManifest{
		ManifestID: hlsManifestResp.Data.Result.ID,
	}

	fmt.Println("Successfully created Live HLS Manifest!")

	/*
		Live Stream
	*/
	fmt.Println("Waiting for live encoding to be ready...")
	liveStreamConfig := &models.LiveStreamConfiguration{
		StreamKey:     stringToPtr("bitmovin"),
		DashManifests: []models.LiveDashManifest{liveDashManifest},
		HLSManifests:  []models.LiveHLSManifest{liveHlsManifest},
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
