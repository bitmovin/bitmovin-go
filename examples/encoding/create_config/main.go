package main

import (
	"log"

	"github.com/streamco/bitmovin-go/bitmovin"
	"github.com/streamco/bitmovin-go/bitmovintypes"
	"github.com/streamco/bitmovin-go/models"
	"github.com/streamco/bitmovin-go/services"
)

func main() {
	bitmovin := bitmovin.NewBitmovinDefault("YOUR API KEY")

	config := models.NewH264CodecConfigBuilder(`H264 Default Config`).
		Width(1920).Height(1080).Bitrate(4500000).
		Framerate(30).RcLookahead(50).
		Profile(bitmovintypes.H264ProfileHigh).
		BFrames(3).RefFrames(5).
		MVPredictionMode(bitmovintypes.MVPredictionModeAuto).MVSearchRangeMax(16).
		CABAC(true).Trellis(bitmovintypes.TrellisEnabledAll).
		Partitions([]bitmovintypes.Partition{bitmovintypes.PartitionI4X4, bitmovintypes.PartitionI8X8}).
		Build()

	svc := services.NewH264CodecConfigurationService(bitmovin)
	response, _ := svc.Create(config)
	log.Printf("Created h264 Code Configuration with ID: %s With Name: %s", *response.Data.Result.ID, *response.Data.Result.Name)
}
