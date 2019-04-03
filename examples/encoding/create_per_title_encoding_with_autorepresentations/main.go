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
		Name:        stringToPtr("Per-Title Encoding with Auto-Representations"),
		CloudRegion: bitmovintypes.CloudRegionGoogleEuropeWest1,
	}
	encodingResp, err := encodingS.Create(encoding)
	errorHandler(err)

	// Add video codec config to encoding
	h264S := services.NewH264CodecConfigurationService(bitmovin)

	vc := &models.H264CodecConfiguration{
		Profile: bitmovintypes.H264ProfileHigh,
	}

	vcResp, err := h264S.Create(vc)
	errorHandler(err)

	// Define per title template video stream
	vis := models.InputStream{
		InputID:       httpResp.Data.Result.ID,
		InputPath:     stringToPtr(fileInputPath),
		SelectionMode: bitmovintypes.SelectionModeAuto,
	}

	visList := []models.InputStream{vis}

	s := &models.Stream{
		CodecConfigurationID: vcResp.Data.Result.ID,
		InputStreams:         visList,
		Mode:                 bitmovintypes.StreamModePerTitleTemplate,
	}

	sResp, err := encodingS.AddStream(*encodingResp.Data.Result.ID, s)
	errorHandler(err)

	aclEntry := models.ACLItem{
		Permission: bitmovintypes.ACLPermissionPublicRead,
	}
	acl := []models.ACLItem{aclEntry}

	ms := models.StreamItem{
		StreamID: sResp.Data.Result.ID,
	}

	mop := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: stringToPtr(fmt.Sprintf("%s/video/{bitrate}", baseOutputPath)),
		ACL:        acl,
	}
	vmx := &models.FMP4Muxing{
		SegmentLength:   floatToPtr(4.0),
		SegmentNaming:   stringToPtr("seg_%number%.m4s"),
		InitSegmentName: stringToPtr("init.mp4"),
		Streams:         []models.StreamItem{ms},
		Outputs:         []models.Output{mop},
	}
	_, err = encodingS.AddFMP4Muxing(*encodingResp.Data.Result.ID, vmx)
	errorHandler(err)

	// Configure per-title options
	perTitle := &models.PerTitle{
		H264Configuration: &models.H264PerTitleConfiguration{
			AutoRepresentations: &models.AutoRepresentations{},
		},
	}

	options := &models.StartOptions{
		PerTitle:     perTitle,
		EncodingMode: bitmovintypes.EncodingModeThreePass,
	}

	// Start the encoding
	_, err = encodingS.StartWithOptions(*encodingResp.Data.Result.ID, options)
	errorHandler(err)

	waitForEncodingToBeFinished(encodingResp, encodingS)
	fmt.Println("Per-Title Encoding finished successfully!")

	/// Generate DASH manifest

	// Create manifest objects
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

	streamsResp, err := encodingS.ListStream(*encodingResp.Data.Result.ID, 0, 20)
	muxingsResp, err := encodingS.ListFMP4Muxing(*encodingResp.Data.Result.ID, 0, 20)
	errorHandler(err)

	for _, stream := range streamsResp.Data.Result.Items {
		if stream.Mode == bitmovintypes.StreamModePerTitleTemplate {
			continue
		}

		muxing := getMuxingOfStream(*stream.ID, muxingsResp)

		muxingOutputPath := *muxing.Outputs[0].OutputPath
		segmentPath := strings.Replace(`/`+muxingOutputPath, baseOutputPath+`/`, "", -1)

		fmp4Rep := &models.FMP4Representation{
			Type:        bitmovintypes.FMP4RepresentationTypeTemplate,
			MuxingID:    muxing.ID,
			EncodingID:  encodingResp.Data.Result.ID,
			SegmentPath: stringToPtr(segmentPath),
		}
		// Add to vod manifest
		_, err := dashService.AddFMP4Representation(*dashMResp.Data.Result.ID, *periodResp.Data.Result.ID, *videoAdaptationSetResp.Data.Result.ID, fmp4Rep)
		errorHandler(err)
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

func floatToPtr(f float64) *float64 {
	return &f
}
