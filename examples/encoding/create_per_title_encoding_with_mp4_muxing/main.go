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

const bitmovinApiKey = "<YOUR BITMOVIN API KEY>"

const s3InputAccessKey = "<YOUR S3 INPUT ACCESS KEY>"
const s3InputSecretKey = "<YOUR S3 INPUT SECRET KEY>"
const s3InputBucket = "<YOUR S3 INPUT BUCKET>"

const s3OutputAccessKey = "<YOUR S3 OUTPUT ACCESS KEY>"
const s3OutputSecretKey = "<YOUR S3 OUTPUT SECRET KEY>"
const s3OutputBucket = "<YOUR S3 OUTPUT BUCKET>"

const fileInputPath = "/path/to/your/input/file.mkv"

const baseOutputPath = "/path/to/your/output"

func main() {
	// Creating Bitmovin object
	bitmovin := bitmovin.NewBitmovinDefault(bitmovinApiKey)

	// Create the input resource to access the input file
	s3IS := services.NewS3InputService(bitmovin)
	s3Input := &models.S3Input{
		AccessKey:   stringToPtr(s3InputAccessKey),
		SecretKey:   stringToPtr(s3InputSecretKey),
		BucketName:  stringToPtr(s3InputBucket),
		CloudRegion: bitmovintypes.AWSCloudRegionEUWest1,
	}
	s3InputResp, err := s3IS.Create(s3Input)
	errorHandler(err)

	// Create the output resource to write the output files
	s3OS := services.NewS3OutputService(bitmovin)
	s3Output := &models.S3Output{
		AccessKey:   stringToPtr(s3OutputAccessKey),
		SecretKey:   stringToPtr(s3OutputSecretKey),
		BucketName:  stringToPtr(s3OutputBucket),
		CloudRegion: bitmovintypes.AWSCloudRegionEUWest1,
	}
	s3OutputResp, err := s3OS.Create(s3Output)
	errorHandler(err)

	// The encoding is created. The cloud region is set to AUTO to use the best cloud region depending on the input
	encodingS := services.NewEncodingService(bitmovin)
	encoding := &models.Encoding{
		Name:        stringToPtr("Go Example - Per-Title"),
		CloudRegion: bitmovintypes.CloudRegionGoogleEuropeWest1,
	}
	encodingResp, err := encodingS.Create(encoding)
	errorHandler(err)

	// Select the video and audio input stream that should be encoded
	vis := models.InputStream{
		InputID:       s3InputResp.Data.Result.ID,
		InputPath:     stringToPtr(fileInputPath),
		SelectionMode: bitmovintypes.SelectionModeAuto,
	}

	ais := models.InputStream{
		InputID:       s3InputResp.Data.Result.ID,
		InputPath:     stringToPtr(fileInputPath),
		SelectionMode: bitmovintypes.SelectionModeAuto,
	}

	h264S := services.NewH264CodecConfigurationService(bitmovin)
	aacS := services.NewAACCodecConfigurationService(bitmovin)

	vsResp := createPerTitleVideoStream(encodingResp, vis, h264S, encodingS)
	asResp := createAudioStream(encodingResp, ais, aacS, encodingS)

	createMp4Muxing(encodingResp, s3OutputResp, vsResp, asResp, encodingS)
	startEncoding(encodingResp, encodingS)
}

// This will create the Per-Title template video stream. This stream will be used as a template for the Per-Title
// encoding. The Codec Configuration, Muxings, DRMs and Filters applied to the generated Per-Title profile will be
// based on the same, or closest matching resolutions defined in the template.
// Please note, that template streams are not necessarily used for the encoding -
// they are just used as template.
// encodingResp: The reference of the encoding
// vis: The input stream that should be encoded
// h264S: The H264 configuration service
// encodingS: The encoding service
// Return: The created Per-Title template video stream. This will be used later for the MP4 muxing
func createPerTitleVideoStream(encodingResp *models.EncodingResponse, vis models.InputStream, h264S *services.H264CodecConfigurationService, encodingS *services.EncodingService) *models.StreamResponse {
	vc := &models.H264CodecConfiguration{
		Profile: bitmovintypes.H264ProfileHigh,
	}

	vcResp, err := h264S.Create(vc)
	errorHandler(err)

	vs := &models.Stream{
		CodecConfigurationID: vcResp.Data.Result.ID,
		InputStreams:         []models.InputStream{vis},
		Mode:                 bitmovintypes.StreamModePerTitleTemplate,
	}

	vsResp, err := encodingS.AddStream(*encodingResp.Data.Result.ID, vs)
	errorHandler(err)

	return vsResp
}

// This will create the audio stream that will be encoded with the given codec configuration.
// encodingResp: The reference of the encoding
// ais: The input stream that should be encoded
// aacS: The AAC configuration service
// encodingS: The encoding service
// Return: The created audio stream. This will be used later for the MP4 muxing
func createAudioStream(encodingResp *models.EncodingResponse, ais models.InputStream, aacS *services.AACCodecConfigurationService, encodingS *services.EncodingService) *models.StreamResponse {

	// AUDIO
	// Add audio codec config
	ac := &models.AACCodecConfiguration{
		Name:         stringToPtr("AAC audio configuration"),
		Bitrate:      intToPtr(128000),
		SamplingRate: floatToPtr(48000),
	}

	acResp, err := aacS.Create(ac)
	errorHandler(err)

	// Add audio stream to encoding
	as := &models.Stream{
		Name:                 stringToPtr("Audio Stream"),
		CodecConfigurationID: acResp.Data.Result.ID,
		InputStreams:         []models.InputStream{ais},
	}

	asResp, err := encodingS.AddStream(*encodingResp.Data.Result.ID, as)
	errorHandler(err)

	return asResp
}

// An MP4 muxing will be created for with the Per-Title video stream template and the audio stream.
// This muxing must define either {uuid} or {bitrate} in the output path. These placeholders will be replaced during
// the generation of the Per-Title.
// encodingResp: The reference of the encoding
// s3OutputResp: The output the files should be written to
// vsResp: The Per-Title template video stream
// asResp: The audio stream
// encodingS: The encoding service
func createMp4Muxing(encodingResp *models.EncodingResponse, s3OutputResp *models.S3OutputResponse, vsResp *models.StreamResponse, asResp *models.StreamResponse, encodingS *services.EncodingService) {

	// Output acl settings
	aclEntry := models.ACLItem{
		Permission: bitmovintypes.ACLPermissionPublicRead,
	}
	acl := []models.ACLItem{aclEntry}

	mop := models.Output{
		OutputID:   s3OutputResp.Data.Result.ID,
		OutputPath: stringToPtr(fmt.Sprintf("%s/{width}_{bitrate}_{uuid}/", baseOutputPath)),
		ACL:        acl,
	}

	vms := models.StreamItem{
		StreamID: vsResp.Data.Result.ID,
	}

	ams := models.StreamItem{
		StreamID: asResp.Data.Result.ID,
	}

	mx := &models.MP4Muxing{
		Streams:  []models.StreamItem{vms, ams},
		Filename: stringToPtr(`per_title_mp4.mp4`),
		Outputs:  []models.Output{mop},
		Name:     stringToPtr(`MP4 Muxing`),
	}
	_, err := encodingS.AddMP4Muxing(*encodingResp.Data.Result.ID, mx)
	errorHandler(err)
}

// The encoding will be started with the per title object and the auto representations set. If the auto
// representation is set, stream configurations will be automatically added to the Per-Title profile. In that case
// at least one PER_TITLE_TEMPLATE stream configuration must be available. All other configurations will be
// automatically chosen by the Per-Title algorithm. All relevant settings for streams and muxings will be taken from
// the closest PER_TITLE_TEMPLATE stream defined. The closest stream will be chosen based on the resolution
// specified in the codec configuration.
// encodingResp: The reference of the encoding
// encodingS: The encoding service
func startEncoding(encodingResp *models.EncodingResponse, encodingS *services.EncodingService) {
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
	_, err := encodingS.StartWithOptions(*encodingResp.Data.Result.ID, options)
	errorHandler(err)

	waitForEncodingToBeFinished(encodingResp, encodingS)
	fmt.Println("Per-Title Encoding finished successfully!")
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

func floatToPtr(f float64) *float64 {
	return &f
}
