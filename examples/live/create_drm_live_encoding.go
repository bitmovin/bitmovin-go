package main

import (
	"fmt"
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
	errorHandler(rtmpInputListResp.Status, err)

	if len(rtmpInputListResp.Data.Result.Items) < 1 {
		fmt.Println("No RTMP inputs on account!")
		return
	}

	gcsOutputService := services.NewGCSOutputService(bitmovin)
	gcsOutput := &models.GCSOutput{
		AccessKey: 		stringToPtr("YOUR_ACCESS_KEY"),
		SecretKey: 		stringToPtr("YOUR_SECRET_KEY"),
		BucketName: 	stringToPtr("YOUR_BUCKET_NAME"),
		CloudRegion: 	bitmovintypes.GoogleCloudRegionEuropeWest1,
	}

	t := time.Now()
	outputBasePath := "golang_live_drm_test_" + t.Format("20060102150405")

	gcsOutputResponse, err := gcsOutputService.Create(gcsOutput)
	errorHandler(gcsOutputResponse.Status, err)

	encodingS := services.NewEncodingService(bitmovin)
	encoding := &models.Encoding{
		Name:        stringToPtr("Example Golang Live Encoding with DRM"),
		CloudRegion: bitmovintypes.CloudRegionGoogleUSEast1,
	}
	encodingResp, err := encodingS.Create(encoding)
	errorHandler(encodingResp.Status, err)

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
	errorHandler(videoStream1080pResp.Status, err)
	videoStream720pResp, err := encodingS.AddStream(encodingID, videoStream720p)
	errorHandler(videoStream720pResp.Status, err)

	ais := []models.InputStream{audioInputStream}
	audioStream := &models.Stream{
		CodecConfigurationID: aacResp.Data.Result.ID,
		InputStreams:         ais,
	}
	aacStreamResp, err := encodingS.AddStream(encodingID, audioStream)
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
	errorHandler(videoMuxing1080pResp.Status, err)

	videoMuxing720p := &models.FMP4Muxing{
		SegmentLength:   floatToPtr(4.0),
		SegmentNaming:   stringToPtr("seg_%number%.m4s"),
		InitSegmentName: stringToPtr("init.mp4"),
		Streams:         []models.StreamItem{videoMuxingStream720p},
	}
	videoMuxing720pResp, err := encodingS.AddFMP4Muxing(encodingID, videoMuxing720p)
	errorHandler(videoMuxing720pResp.Status, err)

	audioMuxing := &models.FMP4Muxing{
		SegmentLength:   floatToPtr(4.0),
		SegmentNaming:   stringToPtr("seg_%number%.m4s"),
		InitSegmentName: stringToPtr("init.mp4"),
		Streams:         []models.StreamItem{audioMuxingStream},
	}
	audioMuxingResp, err := encodingS.AddFMP4Muxing(encodingID, audioMuxing)
	errorHandler(audioMuxingResp.Status, err)
	fmt.Println("Successfully created DASH FMP4 Muxings!")

	/*
	HLS TS Muxings
	 */
	fmt.Println("Creating HLS TS Muxings...")
	videoTsMuxing1080p := &models.TSMuxing{
		SegmentLength: floatToPtr(4.0),
		Streams: []models.StreamItem{videoMuxingStream1080p},
	}

	videoTsMuxing1080Resp, err := encodingS.AddTSMuxing(encodingID, videoTsMuxing1080p)
	errorHandler(videoTsMuxing1080Resp.Status, err)

	videoTsMuxing720p := &models.TSMuxing{
		SegmentLength: floatToPtr(4.0),
		Streams: []models.StreamItem{videoMuxingStream720p},
	}
	videoTsMuxing720Resp, err := encodingS.AddTSMuxing(encodingID, videoTsMuxing720p)
	errorHandler(videoTsMuxing720Resp.Status, err)

	audioTsMuxing := &models.TSMuxing{
		SegmentLength: floatToPtr(4.0),
		Streams: []models.StreamItem{audioMuxingStream},
	}
	audioTsMuxingResp, err := encodingS.AddTSMuxing(encodingID, audioTsMuxing)
	errorHandler(audioTsMuxingResp.Status, err)

	/*
	Outputs for DRM FMP4 Muxings
	 */
	videoMuxing1080pOutput := models.Output{
		OutputID:   gcsOutputResponse.Data.Result.ID,
		OutputPath: stringToPtr(outputBasePath + "/video/1080p"),
		ACL:        acl,
	}
	videoMuxing720pOutput := models.Output{
		OutputID:   gcsOutputResponse.Data.Result.ID,
		OutputPath: stringToPtr(outputBasePath + "/video/720p"),
		ACL:        acl,
	}
	audioMuxingOutput := models.Output{
		OutputID:   gcsOutputResponse.Data.Result.ID,
		OutputPath: stringToPtr(outputBasePath + "/audio"),
		ACL:        acl,
	}
	fmt.Println("Successfully created HLS TS Muxings!")

	/*
	Widevine and Playready CENC DRM
	 */
	fmt.Println("Creating Widevine and PlayReady CENC DRMs...")
	widevine := models.WidevineCencDrm{
		PSSH: stringToPtr("WIDEVINE_PSSSH"),
	}

	playready := models.PlayReadyCencDrm{
		LaURL: stringToPtr("PLAYREADY_LA_URL"),
	}

	cencDrm := models.CencDrm{
		Key: 	stringToPtr("YOUR_CENC_DRM_KEY"),
		KID: 	stringToPtr("YOUR_CENC_DRM_KID"),
		Name: 	stringToPtr("My CENC DRM"),
		Widevine: widevine,
		PlayReady: playready,
	}

	drmService := services.NewDrmService(bitmovin)

	cencDrm.Outputs = []models.Output{videoMuxing720pOutput}
	cencResp720, err := drmService.CreateFmp4Drm(encodingID, *videoMuxing720pResp.Data.Result.ID, cencDrm)
	if err != nil {
		fmt.Println(err)
		return
	}
	cencDrm720Response := cencResp720.(*models.CencDrmResponse)
	errorHandler(cencDrm720Response.Status, err)

	cencDrm.Outputs = []models.Output{videoMuxing1080pOutput}
	cencResp1080, err := drmService.CreateFmp4Drm(encodingID, *videoMuxing1080pResp.Data.Result.ID, cencDrm)
	if err != nil {
		fmt.Println(err)
		return
	}
	cencDrm1080Response := cencResp1080.(*models.CencDrmResponse)
	errorHandler(cencDrm1080Response.Status, err)

	cencDrm.Outputs = []models.Output{audioMuxingOutput}
	cencRespAudio, err := drmService.CreateFmp4Drm(encodingID, *audioMuxingResp.Data.Result.ID, cencDrm)
	if err != nil {
		fmt.Println(err)
		return
	}
	cencDrmAudioResponse := cencRespAudio.(*models.CencDrmResponse)
	errorHandler(cencDrmAudioResponse.Status, err)

	fmt.Println("Successfully created Widevine and PlayReady CENC DRMs!")

	/*
	FairPlay DRM
	 */
	fmt.Println("Creating Fairplay DRMs for HLS...")
	fairPlayDrm := models.FairPlayDrm{
		IV:		stringToPtr("YOUR_FAIRPLAY_IV"),
		URI: 	stringToPtr("YOUR_FAIRPLAY_URI"),
		Key: 	stringToPtr("YOUR_FAIRPLAY_KEY"),
		Name:	stringToPtr("My Fairplay DRM"),
	}

	fairPlayDrm.Outputs = []models.Output{videoMuxing720pOutput}
	drmService.CreateTsDrm(encodingID, *videoTsMuxing720Resp.Data.Result.ID, fairPlayDrm)
	fairPlayDrm.Outputs = []models.Output{videoMuxing1080pOutput}
	drmService.CreateFmp4Drm(encodingID, *videoTsMuxing1080Resp.Data.Result.ID, fairPlayDrm)
	fairPlayDrm.Outputs = []models.Output{audioMuxingOutput}
	drmService.CreateFmp4Drm(encodingID, *audioTsMuxingResp.Data.Result.ID, fairPlayDrm)

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
		Name: 		  stringToPtr("your_manifest_name.mpd"),
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
	Video AdaptationSet
	 */
	vas := &models.VideoAdaptationSet{}
	vasResp, err := dashService.AddVideoAdaptationSet(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, vas)
	errorHandler(vasResp.Status, err)

	fmt.Println("Adding content protection to video adaptation set...")
	vcp := &models.AdaptationSetContentProtection{
		DrmId: cencDrm1080Response.Data.Result.ID,
		EncodingId: stringToPtr(encodingID),
		MuxingId: videoMuxing1080pResp.Data.Result.ID,
	}
	vcpResp, err := dashService.AddContentProtectionToAdaptationSet(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *vasResp.Data.Result.ID, vcp)
	errorHandler(vcpResp.Status, err)
	fmt.Println("Successfully added content protection to vide adaptation set!")

	/*
	Audio AdaptationSet
	 */
	aas := &models.AudioAdaptationSet{
		Language: stringToPtr("en"),
	}
	aasResp, err := dashService.AddAudioAdaptationSet(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, aas)
	errorHandler(aasResp.Status, err)

	fmt.Println("Adding content protection to audio adaptation set...")
	acp := &models.AdaptationSetContentProtection{
		DrmId: cencDrmAudioResponse.Data.Result.ID,
		EncodingId: stringToPtr(encodingID),
		MuxingId: audioMuxingResp.Data.Result.ID,
	}
	acpResp, err := dashService.AddContentProtectionToAdaptationSet(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *aasResp.Data.Result.ID, acp)
	errorHandler(acpResp.Status, err)
	fmt.Println("Successfully added content protection to audio adaptation set!")

	/*
	DRM Representations
	 */
	fmp4Rep1080 := &models.DrmFMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    videoMuxing1080pResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: stringToPtr("../video/1080p"),
		DrmID:		 cencDrm1080Response.Data.Result.ID,
	}
	fmp4Rep1080Resp, err := dashService.AddDrmFMP4Representation(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *vasResp.Data.Result.ID, fmp4Rep1080)
	errorHandler(fmp4Rep1080Resp.Status, err)

	fmp4Rep720 := &models.DrmFMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    videoMuxing720pResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: stringToPtr("../video/720p"),
		DrmID:		 cencDrm720Response.Data.Result.ID,
	}
	fmp4Rep720Resp, err := dashService.AddDrmFMP4Representation(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *vasResp.Data.Result.ID, fmp4Rep720)
	errorHandler(fmp4Rep720Resp.Status, err)

	fmp4RepAudio := &models.DrmFMP4Representation{
		Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
		MuxingID:    audioMuxingResp.Data.Result.ID,
		EncodingID:  encodingResp.Data.Result.ID,
		SegmentPath: stringToPtr("../audio"),
		DrmID:		 cencDrmAudioResponse.Data.Result.ID,
	}
	fmp4RepAudioResp, err := dashService.AddDrmFMP4Representation(*dashManifestResp.Data.Result.ID, *periodResp.Data.Result.ID, *aasResp.Data.Result.ID, fmp4RepAudio)
	errorHandler(fmp4RepAudioResp.Status, err)

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
		MuxingID:        audioTsMuxingResp.Data.Result.ID,
	}
	audioMediaInfoResp, err := hlsService.AddMediaInfo(*hlsManifestResp.Data.Result.ID, audioMediaInfo)
	errorHandler(audioMediaInfoResp.Status, err)

	video1080pStreamInfo := &models.StreamInfo{
		Audio:       stringToPtr("audio_group"),
		SegmentPath: stringToPtr("../video/1080p"),
		URI:         stringToPtr("video_hi.m3u8"),
		EncodingID:  encodingResp.Data.Result.ID,
		StreamID:    videoStream1080pResp.Data.Result.ID,
		MuxingID:    videoTsMuxing1080Resp.Data.Result.ID,
	}
	video1080pStreamInfoResponse, err := hlsService.AddStreamInfo(*hlsManifestResp.Data.Result.ID, video1080pStreamInfo)
	errorHandler(video1080pStreamInfoResponse.Status, err)

	video720pStreamInfo := &models.StreamInfo{
		Audio:       stringToPtr("audio_group"),
		SegmentPath: stringToPtr("../video/720p"),
		URI:         stringToPtr("video_lo.m3u8"),
		EncodingID:  encodingResp.Data.Result.ID,
		StreamID:    videoStream720pResp.Data.Result.ID,
		MuxingID:    videoTsMuxing720Resp.Data.Result.ID,
	}
	video720pStreamInfoResponse, err := hlsService.AddStreamInfo(*hlsManifestResp.Data.Result.ID, video720pStreamInfo)
	errorHandler(video720pStreamInfoResponse.Status, err)

	liveHlsManifest := models.LiveHLSManifest{
		ManifestID:     hlsManifestResp.Data.Result.ID,
	}

	fmt.Println("Successfully created Live HLS Manifest!")

	/*
	Live Stream
	 */
	fmt.Println("Waiting for live encoding to be ready...")
	liveStreamConfig := &models.LiveStreamConfiguration{
		StreamKey:     stringToPtr("bitmovin"),
		DashManifests: []models.LiveDashManifest{liveDashManifest},
		HLSManifests: []models.LiveHLSManifest{liveHlsManifest},
	}

	startResp, err := encodingS.StartLive(encodingID, liveStreamConfig)
	errorHandler(startResp.Status, err)

	for numRetries := 0; numRetries < MaxRetries; numRetries++ {
		time.Sleep(10 * time.Second)
		statusResp, err := encodingS.RetrieveLiveStatus(encodingID)
		if err != nil {
			if err.Error() != "ERROR 2023: Live encoding details not available!" {
				fmt.Println("Error in starting live encoding")
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