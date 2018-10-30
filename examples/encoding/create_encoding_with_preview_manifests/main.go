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

type h264VideoRepresentation struct {
	name         *string
	outputPath   *string
	width        *int64
	bitrate      *int64
	fps          *float64
	audioGroupId *string
}

type aacAudioRep struct {
	name         *string
	outputPath   *string
	bitrate      *int64
	sampleRate   *float64
	audioGroupId *string
	lang         *string
}

type encodingConfig struct {
	vRep       *h264VideoRepresentation
	aRep       *aacAudioRep
	vConfResp  *models.H264CodecConfigurationResponse
	aConfResp  *models.AACCodecConfigurationResponse
	streamResp *models.StreamResponse
	fmp4MxResp *models.FMP4MuxingResponse
	tsMxResp   *models.TSMuxingResponse
}

const bitmovinApiKey = "<YOUR BITMOVIN API KEY>"

const vodHlsManifestName = "myVodHlsManifest.m3u8"
const vodDashManifestName = "myVodDashManifest.mpd"
const previewHlsManifestName = "myPreviewHlsManifest.m3u8"
const previewDashManifestName = "myPreviewDashManifest.mpd"

const httpsInput = "your.host.com"
const fileInputPath = "/path/to/your/input/file.mkv"

const s3OutputAccessKey = "<YOUR S3 OUTPUT ACCESS KEY>"
const s3OutputSecretKey = "<YOUR S3 OUTPUT SECRET KEY>"
const s3OutputBucket = "<YOUR S3 OUTPUT BUCKET>"
const baseOutputPath = "/path/to/your/output"

func main() {
	// Creating Bitmovin object
	bitmovin := bitmovin.NewBitmovin(bitmovinApiKey, "https://api.bitmovin.com/v1/", 5)

	videoRepresentations := []*h264VideoRepresentation{
		{width: intToPtr(426), bitrate: intToPtr(300000), fps: nil, name: stringToPtr("300 kbps profile"), outputPath: stringToPtr("video/300kbps"), audioGroupId: stringToPtr("audio1")},
		{width: intToPtr(640), bitrate: intToPtr(500000), fps: nil, name: stringToPtr("500 kbps profile"), outputPath: stringToPtr("video/500kbps"), audioGroupId: stringToPtr("audio1")},
		{width: intToPtr(852), bitrate: intToPtr(800000), fps: nil, name: stringToPtr("800 kbps profile"), outputPath: stringToPtr("video/800kbps"), audioGroupId: stringToPtr("audio1")},
		{width: intToPtr(1280), bitrate: intToPtr(1500000), fps: nil, name: stringToPtr("1500 kbps profile"), outputPath: stringToPtr("video/1500kbps"), audioGroupId: stringToPtr("audio1")},
		{width: intToPtr(1920), bitrate: intToPtr(8000000), fps: nil, name: stringToPtr("8000 kbps profile"), outputPath: stringToPtr("video/8000kbps"), audioGroupId: stringToPtr("audio1")},
	}

	audioRepresentations := []*aacAudioRep{
		{
			sampleRate:   floatToPtr(44100.0),
			bitrate:      intToPtr(128000),
			name:         stringToPtr("128 kbit audio profile"),
			outputPath:   stringToPtr("audio/128kbit"),
			audioGroupId: stringToPtr("audio1"),
			lang:         stringToPtr("en"),
		},
	}

	// Add video and audio representations to encoding config
	var encodingConfigs []*encodingConfig
	for _, vr := range videoRepresentations {
		encodingConfigs = append(encodingConfigs, &encodingConfig{
			vRep: vr,
		})
	}
	for _, ar := range audioRepresentations {
		encodingConfigs = append(encodingConfigs, &encodingConfig{
			aRep: ar,
		})
	}

	// Creating the HTTP Input
	httpIS := services.NewHTTPInputService(bitmovin)
	httpInput := &models.HTTPInput{
		Host: stringToPtr(httpsInput),
	}
	httpResp, err := httpIS.Create(httpInput)
	errorHandler(err)

	//Creating S3 Output
	s3OS := services.NewS3OutputService(bitmovin)
	s3Output := &models.S3Output{
		AccessKey:   stringToPtr(s3OutputAccessKey),
		SecretKey:   stringToPtr(s3OutputSecretKey),
		BucketName:  stringToPtr(s3OutputBucket),
		CloudRegion: bitmovintypes.AWSCloudRegionEUWest1,
	}
	s3OutputResp, err := s3OS.Create(s3Output)
	errorHandler(err)

	// Create encoding
	encodingS := services.NewEncodingService(bitmovin)
	encoding := &models.Encoding{
		Name:        stringToPtr("My Golang Example Encoding"),
		CloudRegion: bitmovintypes.CloudRegionGoogleEuropeWest1,
	}
	encodingResp, err := encodingS.Create(encoding)
	errorHandler(err)

	// Add audio and video codec configs to encoding
	h264S := services.NewH264CodecConfigurationService(bitmovin)
	aacS := services.NewAACCodecConfigurationService(bitmovin)

	for _, ec := range encodingConfigs {
		if ec.vRep != nil {
			vc := &models.H264CodecConfiguration{
				Name:    ec.vRep.name,
				Bitrate: ec.vRep.bitrate,
				Width:   ec.vRep.width,
				Profile: bitmovintypes.H264ProfileHigh,
			}
			vcResp, err := h264S.Create(vc)
			errorHandler(err)
			ec.vConfResp = vcResp
		}
		if ec.aRep != nil {
			ac := &models.AACCodecConfiguration{
				Name:         ec.aRep.name,
				Bitrate:      ec.aRep.bitrate,
				SamplingRate: ec.aRep.sampleRate,
			}
			acResp, err := aacS.Create(ac)
			errorHandler(err)
			ec.aConfResp = acResp
		}
	}

	// Define video and audio input streams
	vis := models.InputStream{
		InputID:       httpResp.Data.Result.ID,
		InputPath:     stringToPtr(fileInputPath),
		SelectionMode: bitmovintypes.SelectionModeAuto,
	}
	ais := models.InputStream{
		InputID:       httpResp.Data.Result.ID,
		InputPath:     stringToPtr(fileInputPath),
		SelectionMode: bitmovintypes.SelectionModeAuto,
	}

	visList := []models.InputStream{vis}
	aisList := []models.InputStream{ais}

	//Add streams to encoding
	for _, ec := range encodingConfigs {
		if ec.vConfResp != nil {
			s := &models.Stream{
				CodecConfigurationID: ec.vConfResp.Data.Result.ID,
				InputStreams:         visList,
			}
			sResp, err := encodingS.AddStream(*encodingResp.Data.Result.ID, s)
			errorHandler(err)
			ec.streamResp = sResp
		}
		if ec.aConfResp != nil {
			s := &models.Stream{
				CodecConfigurationID: ec.aConfResp.Data.Result.ID,
				InputStreams:         aisList,
			}
			sResp, err := encodingS.AddStream(*encodingResp.Data.Result.ID, s)
			errorHandler(err)
			ec.streamResp = sResp
		}
	}

	aclEntry := models.ACLItem{
		Permission: bitmovintypes.ACLPermissionPublicRead,
	}
	acl := []models.ACLItem{aclEntry}

	//Add FMP4 muxings to encoding
	for _, ec := range encodingConfigs {
		if ec.streamResp != nil {
			ms := models.StreamItem{
				StreamID: ec.streamResp.Data.Result.ID,
			}

			if ec.vRep != nil {
				mop := models.Output{
					OutputID:   s3OutputResp.Data.Result.ID,
					OutputPath: stringToPtr(fmt.Sprintf("%s/%s_dash", baseOutputPath, *ec.vRep.outputPath)),
					ACL:        acl,
				}
				vmx := &models.FMP4Muxing{
					SegmentLength:   floatToPtr(4.0),
					SegmentNaming:   stringToPtr("seg_%number%.m4s"),
					InitSegmentName: stringToPtr("init.mp4"),
					Streams:         []models.StreamItem{ms},
					Outputs:         []models.Output{mop},
				}
				vmxResp, err := encodingS.AddFMP4Muxing(*encodingResp.Data.Result.ID, vmx)
				errorHandler(err)
				ec.fmp4MxResp = vmxResp
			}
			if ec.aRep != nil {
				mop := models.Output{
					OutputID:   s3OutputResp.Data.Result.ID,
					OutputPath: stringToPtr(fmt.Sprintf("%s/%s_dash", baseOutputPath, *ec.aRep.outputPath)),
					ACL:        acl,
				}
				amx := &models.FMP4Muxing{
					SegmentLength:   floatToPtr(4.0),
					SegmentNaming:   stringToPtr("seg_%number%.m4s"),
					InitSegmentName: stringToPtr("init.mp4"),
					Streams:         []models.StreamItem{ms},
					Outputs:         []models.Output{mop},
				}
				amxResp, err := encodingS.AddFMP4Muxing(*encodingResp.Data.Result.ID, amx)
				errorHandler(err)
				ec.fmp4MxResp = amxResp
			}

		}
	}

	//Add TS muxings to encoding
	for _, ec := range encodingConfigs {
		if ec.streamResp != nil {
			ms := models.StreamItem{
				StreamID: ec.streamResp.Data.Result.ID,
			}

			if ec.vRep != nil {
				mop := models.Output{
					OutputID:   s3OutputResp.Data.Result.ID,
					OutputPath: stringToPtr(fmt.Sprintf("%s/%s_hls", baseOutputPath, *ec.vRep.outputPath)),
					ACL:        acl,
				}
				vmx := &models.TSMuxing{
					SegmentLength: floatToPtr(4.0),
					SegmentNaming: stringToPtr("seg_%number%.m4s"),
					Streams:       []models.StreamItem{ms},
					Outputs:       []models.Output{mop},
				}
				vmxResp, err := encodingS.AddTSMuxing(*encodingResp.Data.Result.ID, vmx)
				errorHandler(err)
				ec.tsMxResp = vmxResp
			}
			if ec.aRep != nil {
				mop := models.Output{
					OutputID:   s3OutputResp.Data.Result.ID,
					OutputPath: stringToPtr(fmt.Sprintf("%s/%s_hls", baseOutputPath, *ec.aRep.outputPath)),
					ACL:        acl,
				}
				amx := &models.TSMuxing{
					SegmentLength: floatToPtr(4.0),
					SegmentNaming: stringToPtr("seg_%number%.m4s"),
					Streams:       []models.StreamItem{ms},
					Outputs:       []models.Output{mop},
				}
				amxResp, err := encodingS.AddTSMuxing(*encodingResp.Data.Result.ID, amx)
				errorHandler(err)
				ec.tsMxResp = amxResp
			}
		}
	}

	//Create manifest objects
	hlsService := services.NewHLSManifestService(bitmovin)
	dashService := services.NewDashManifestService(bitmovin)

	manifestOutput := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: stringToPtr(baseOutputPath),
		ACL:        acl,
	}

	// Dash Manifests
	vodDashM := &models.DashManifest{
		ManifestName: stringToPtr(vodDashManifestName),
		Outputs:      []models.Output{manifestOutput},
	}
	vodDashMResp, err := dashService.Create(vodDashM)
	errorHandler(err)

	previewDashM := &models.DashManifest{
		ManifestName: stringToPtr(previewDashManifestName),
		Outputs:      []models.Output{manifestOutput},
	}
	previewDashMResp, err := dashService.Create(previewDashM)
	errorHandler(err)

	period := &models.Period{}
	//Add to vod manifest
	vodPeriodResp, err := dashService.AddPeriod(*vodDashMResp.Data.Result.ID, period)
	errorHandler(err)
	// Add to preview manifest
	previewPeriodResp, err := dashService.AddPeriod(*previewDashMResp.Data.Result.ID, period)
	errorHandler(err)

	vas := &models.VideoAdaptationSet{}
	// Add to vod manifest
	vodVasResp, err := dashService.AddVideoAdaptationSet(*vodDashMResp.Data.Result.ID, *vodPeriodResp.Data.Result.ID, vas)
	errorHandler(err)

	// Add to preview manifest
	previewVasResp, err := dashService.AddVideoAdaptationSet(*previewDashMResp.Data.Result.ID, *previewPeriodResp.Data.Result.ID, vas)
	errorHandler(err)

	vodAasMap := make(map[string]*models.AudioAdaptationSetResponse)
	previewAasMap := make(map[string]*models.AudioAdaptationSetResponse)
	for _, ec := range encodingConfigs {
		if ec.aRep != nil {
			aas := &models.AudioAdaptationSet{
				Language: ec.aRep.lang,
			}
			// Add to vod manifest
			vodAasResp, err := dashService.AddAudioAdaptationSet(*vodDashMResp.Data.Result.ID, *vodPeriodResp.Data.Result.ID, aas)
			errorHandler(err)

			// Add to preview manifest
			previewAasResp, err := dashService.AddAudioAdaptationSet(*previewDashMResp.Data.Result.ID, *previewPeriodResp.Data.Result.ID, aas)
			errorHandler(err)

			vodAasMap[*ec.aRep.lang] = vodAasResp
			previewAasMap[*ec.aRep.lang] = previewAasResp
		}
	}

	// Hls Manifests
	vodHlsM := &models.HLSManifest{
		ManifestName: stringToPtr(vodHlsManifestName),
		Outputs:      []models.Output{manifestOutput},
	}
	vodHlsMResp, err := hlsService.Create(vodHlsM)
	errorHandler(err)

	previewHlsM := &models.HLSManifest{
		ManifestName: stringToPtr(previewHlsManifestName),
		Outputs:      []models.Output{manifestOutput},
	}
	previewHlsMResp, err := hlsService.Create(previewHlsM)
	errorHandler(err)

	for _, ec := range encodingConfigs {
		if ec.vRep != nil {
			if ec.fmp4MxResp != nil {
				fmp4Rep := &models.FMP4Representation{
					Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
					MuxingID:    ec.fmp4MxResp.Data.Result.ID,
					EncodingID:  encodingResp.Data.Result.ID,
					SegmentPath: stringToPtr(fmt.Sprintf("%s_dash", *ec.vRep.outputPath)),
				}
				// Add to vod manifest
				_, err := dashService.AddFMP4Representation(*vodDashMResp.Data.Result.ID, *vodPeriodResp.Data.Result.ID, *vodVasResp.Data.Result.ID, fmp4Rep)
				errorHandler(err)

				// Add to preview manifest
				_, err = dashService.AddFMP4Representation(*previewDashMResp.Data.Result.ID, *previewPeriodResp.Data.Result.ID, *previewVasResp.Data.Result.ID, fmp4Rep)
				errorHandler(err)
			}
			if ec.tsMxResp != nil {
				vsi := &models.StreamInfo{
					Audio:       ec.vRep.audioGroupId,
					URI:         stringToPtr(fmt.Sprintf("video_%d.m3u8", *ec.vRep.bitrate)),
					EncodingID:  encodingResp.Data.Result.ID,
					SegmentPath: stringToPtr(fmt.Sprintf("%s_hls", *ec.vRep.outputPath)),
					StreamID:    ec.streamResp.Data.Result.ID,
					MuxingID:    ec.tsMxResp.Data.Result.ID,
				}
				// Add to vod manifest
				_, err := hlsService.AddStreamInfo(*vodHlsMResp.Data.Result.ID, vsi)
				errorHandler(err)

				// Add to preview manifest
				_, err = hlsService.AddStreamInfo(*previewHlsMResp.Data.Result.ID, vsi)
				errorHandler(err)
			}
		}
		if ec.aRep != nil {
			if ec.fmp4MxResp != nil {
				fmp4Rep := &models.FMP4Representation{
					Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
					MuxingID:    ec.fmp4MxResp.Data.Result.ID,
					EncodingID:  encodingResp.Data.Result.ID,
					SegmentPath: stringToPtr(fmt.Sprintf("%s_dash", *ec.aRep.outputPath)),
				}
				vodAasResp := vodAasMap[*ec.aRep.lang]
				previewAasResp := previewAasMap[*ec.aRep.lang]

				// Add to vod manifest
				_, err := dashService.AddFMP4Representation(*vodDashMResp.Data.Result.ID, *vodPeriodResp.Data.Result.ID, *vodAasResp.Data.Result.ID, fmp4Rep)
				errorHandler(err)

				//Add to preview manifest
				_, err = dashService.AddFMP4Representation(*previewDashMResp.Data.Result.ID, *previewPeriodResp.Data.Result.ID, *previewAasResp.Data.Result.ID, fmp4Rep)
				errorHandler(err)
			}
			if ec.tsMxResp != nil {
				audioMediaInfo := &models.MediaInfo{
					Type:            bitmovintypes.MediaTypeAudio,
					URI:             stringToPtr(fmt.Sprintf("audio_%d.m3u8", *ec.aRep.bitrate)),
					Name:            stringToPtr(fmt.Sprintf("Audio Media Info %s %dkbit", *ec.aRep.audioGroupId, *ec.aRep.bitrate)),
					IsDefault:       boolToPtr(false),
					Autoselect:      boolToPtr(false),
					Forced:          boolToPtr(false),
					Characteristics: []string{"public.accessibility.describes-video"},
					GroupID:         ec.aRep.audioGroupId,
					Language:        ec.aRep.lang,
					SegmentPath:     stringToPtr(fmt.Sprintf("%s_hls", *ec.aRep.outputPath)),
					EncodingID:      encodingResp.Data.Result.ID,
					StreamID:        ec.streamResp.Data.Result.ID,
					MuxingID:        ec.tsMxResp.Data.Result.ID,
				}
				// Add to vod manifest
				_, err := hlsService.AddMediaInfo(*vodHlsMResp.Data.Result.ID, audioMediaInfo)
				errorHandler(err)

				// Add to preview manifest
				_, err = hlsService.AddMediaInfo(*previewHlsMResp.Data.Result.ID, audioMediaInfo)
				errorHandler(err)
			}
		}

	}

	// Start encoding with manifests for preview and vod
	vodHls := models.VodHlsManifest{
		ManifestID: *vodHlsMResp.Data.Result.ID,
	}
	previewHls := models.PreviewHlsManifest{
		ManifestID: *previewHlsMResp.Data.Result.ID,
	}
	vodDash := models.VodDashManifest{
		ManifestID: *vodDashMResp.Data.Result.ID,
	}
	previewDash := models.PreviewDashManifest{
		ManifestID: *previewDashMResp.Data.Result.ID,
	}

	options := &models.StartOptions{
		VodHlsManifests:      []models.VodHlsManifest{vodHls},
		PreviewHlsManifests:  []models.PreviewHlsManifest{previewHls},
		VodDashManifests:     []models.VodDashManifest{vodDash},
		PreviewDashManifests: []models.PreviewDashManifest{previewDash},
	}

	_, err = encodingS.StartWithOptions(*encodingResp.Data.Result.ID, options)
	errorHandler(err)

	waitForEncodingToBeFinished(encodingResp, encodingS)
	fmt.Println("Encoding finished successfully!")
}

func waitForEncodingToBeFinished(encodingResp *models.EncodingResponse, encodingS *services.EncodingService) {
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
		fmt.Printf("ENCODING STATUS: %s\n", status)
		if status == "ERROR" {
			fmt.Println("error in Encoding Status")
			fmt.Printf("STATUS: %s\n", status)
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
