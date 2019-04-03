package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/streamco/bitmovin-go/bitmovin"
	"github.com/streamco/bitmovin-go/bitmovintypes"
	"github.com/streamco/bitmovin-go/models"
	"github.com/streamco/bitmovin-go/services"
)

const (
	minBitrate int64 = 300 * 1000
	maxBitrate int64 = 4000 * 1000
)

type h264VideoRepresentation struct {
	name         *string
	outputPath   *string
	width        *int64
	height       *int64
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
}

const bitmovinApiKey = "<YOUR BITMOVIN API KEY>"

const s3OutputAccessKey = "<YOUR S3 OUTPUT ACCESS KEY>"
const s3OutputSecretKey = "<YOUR S3 OUTPUT SECRET KEY>"
const s3OutputBucket = "<YOUR S3 OUTPUT BUCKET>"

const httpsInput = "your.host.com"
const fileInputPath = "/path/to/your/input/file.mkv"

const baseOutputPath = "/path/to/your/output"
const dashManifestName = "stream.mpd"

func main() {
	// Creating Bitmovin object
	bitmovin := bitmovin.NewBitmovinDefault(bitmovinApiKey)

	videoRepresentations := []*h264VideoRepresentation{
		{height: intToPtr(360), fps: nil, name: stringToPtr("360p profile"), outputPath: stringToPtr("video/360p/{bitrate}"), audioGroupId: stringToPtr("audio1")},
		{height: intToPtr(432), fps: nil, name: stringToPtr("432p profile"), outputPath: stringToPtr("video/432p/{bitrate}"), audioGroupId: stringToPtr("audio1")},
		{height: intToPtr(576), fps: nil, name: stringToPtr("576p profile"), outputPath: stringToPtr("video/576p/{bitrate}"), audioGroupId: stringToPtr("audio1")},
		{height: intToPtr(720), fps: nil, name: stringToPtr("720p profile"), outputPath: stringToPtr("video/720p/{bitrate}"), audioGroupId: stringToPtr("audio1")},
		{height: intToPtr(1080), fps: nil, name: stringToPtr("1080p profile"), outputPath: stringToPtr("video/1080p/{bitrate}"), audioGroupId: stringToPtr("audio1")},
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
		Name:        stringToPtr("Example Per-Title Encoding"),
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
				Width:   ec.vRep.width,
				Height:  ec.vRep.height,
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
				Mode:                 bitmovintypes.StreamModePerTitleTemplate,
				CustomData: map[string]interface{}{
					"type": "video",
				},
			}
			sResp, err := encodingS.AddStream(*encodingResp.Data.Result.ID, s)
			errorHandler(err)
			ec.streamResp = sResp
		}
		if ec.aConfResp != nil {
			s := &models.Stream{
				CodecConfigurationID: ec.aConfResp.Data.Result.ID,
				InputStreams:         aisList,
				CustomData: map[string]interface{}{
					"type": "audio",
					"lang": ec.aRep.lang,
				},
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
				_, err := encodingS.AddFMP4Muxing(*encodingResp.Data.Result.ID, vmx)
				errorHandler(err)
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
				_, err := encodingS.AddFMP4Muxing(*encodingResp.Data.Result.ID, amx)
				errorHandler(err)
			}

		}
	}

	perTitle := &models.PerTitle{
		H264Configuration: &models.H264PerTitleConfiguration{
			MinBitrate: intToPtr(minBitrate),
			MaxBitrate: intToPtr(maxBitrate),
		},
	}

	options := &models.StartOptions{
		PerTitle: perTitle,
	}

	_, err = encodingS.StartWithOptions(*encodingResp.Data.Result.ID, options)
	errorHandler(err)

	waitForEncodingToBeFinished(encodingResp, encodingS)
	fmt.Println("Per-Title Encoding finished successfully!")

	///// Generate DASH manifest

	//Create manifest objects
	dashService := services.NewDashManifestService(bitmovin)

	manifestOutput := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: stringToPtr(baseOutputPath),
		ACL:        acl,
	}

	// Dash Manifests
	dashM := &models.DashManifest{
		ManifestName: stringToPtr(dashManifestName),
		Outputs:      []models.Output{manifestOutput},
	}
	dashMResp, err := dashService.Create(dashM)
	errorHandler(err)

	period := &models.Period{}
	//Add to vod manifest
	periodResp, err := dashService.AddPeriod(*dashMResp.Data.Result.ID, period)
	errorHandler(err)

	videoAdaptationSet := &models.VideoAdaptationSet{}
	// Add to vod manifest
	videoAdaptationSetResp, err := dashService.AddVideoAdaptationSet(*dashMResp.Data.Result.ID, *periodResp.Data.Result.ID, videoAdaptationSet)
	errorHandler(err)

	vodAasMap := make(map[string]*models.AudioAdaptationSetResponse)
	for _, ec := range encodingConfigs {
		if ec.aRep != nil {
			aas := &models.AudioAdaptationSet{
				Language: ec.aRep.lang,
			}
			// Add to vod manifest
			vodAasResp, err := dashService.AddAudioAdaptationSet(*dashMResp.Data.Result.ID, *periodResp.Data.Result.ID, aas)
			errorHandler(err)

			vodAasMap[*ec.aRep.lang] = vodAasResp
		}
	}

	streamsResp, err := encodingS.ListStream(*encodingResp.Data.Result.ID, 0, 20)
	muxingsResp, err := encodingS.ListFMP4Muxing(*encodingResp.Data.Result.ID, 0, 20)
	errorHandler(err)

	for _, stream := range streamsResp.Data.Result.Items {
		streamCustomDataResp, err := encodingS.RetrieveStreamCustomData(*encodingResp.Data.Result.ID, *stream.ID, 0, 20)
		errorHandler(err)

		muxing := getMuxingOfStream(*stream.ID, muxingsResp)

		muxingOutputPath := *muxing.Outputs[0].OutputPath
		segmentPath := strings.Replace(`/`+muxingOutputPath, baseOutputPath+`/`, "", -1)

		if stream.Mode == bitmovintypes.StreamModePerTitleResult {
			_, err := h264S.Retrieve(*stream.CodecConfigurationID)
			errorHandler(err)

			fmp4Rep := &models.FMP4Representation{
				Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
				MuxingID:    muxing.ID,
				EncodingID:  encodingResp.Data.Result.ID,
				SegmentPath: stringToPtr(segmentPath),
			}
			// Add to vod manifest
			_, err = dashService.AddFMP4Representation(*dashMResp.Data.Result.ID, *periodResp.Data.Result.ID, *videoAdaptationSetResp.Data.Result.ID, fmp4Rep)
			errorHandler(err)
		} else if streamCustomDataResp.Data.Result.CustomData["type"].(string) == "audio" {
			fmp4Rep := &models.FMP4Representation{
				Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
				MuxingID:    muxing.ID,
				EncodingID:  encodingResp.Data.Result.ID,
				SegmentPath: stringToPtr(segmentPath),
			}
			vodAasResp := vodAasMap[streamCustomDataResp.Data.Result.CustomData["lang"].(string)]

			// Add to vod manifest
			_, err := dashService.AddFMP4Representation(*dashMResp.Data.Result.ID, *periodResp.Data.Result.ID, *vodAasResp.Data.Result.ID, fmp4Rep)
			errorHandler(err)
		}
	}

	fmt.Println("Starting dash manifest creation")
	_, err = dashService.Start(*dashMResp.Data.Result.ID)
	errorHandler(err)
	fmt.Println("Dash manifest creation finished successfully!")

	waitForDashManifestCreationToBeFinished(*dashMResp.Data.Result.ID, dashService)
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

func waitForDashManifestCreationToBeFinished(dashManifestId string, dashService *services.DashManifestService) {
	status := ""
	for status != "FINISHED" {
		time.Sleep(5 * time.Second)
		statusResp, err := dashService.RetrieveStatus(dashManifestId)
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
}

func getMuxingOfStream(streamId string, muxingsResp *models.FMP4MuxingListResponse) *models.FMP4Muxing {
	for _, m := range muxingsResp.Data.Result.Items {
		if *(m.Streams[0].StreamID) == streamId {
			return &m
		}
	}

	return nil
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
