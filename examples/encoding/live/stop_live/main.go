package main

import (
	"fmt"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/bitmovintypes"
	"github.com/bitmovin/bitmovin-go/services"
)

func main() {
	// Creating Bitmovin object
	bitmovin := bitmovin.NewBitmovin("YOUR API KEY", "https://api.bitmovin.com/v1/", 5)
	encodingS := services.NewEncodingService(bitmovin)
	a, err := encodingS.StopLive("YOUR ENCODING ID")
	errorHandler(a.Status, err)
}

func errorHandler(responseStatus bitmovintypes.ResponseStatus, err error) {
	if err != nil {
		fmt.Println("go error")
		fmt.Println(err)
	} else if responseStatus == "ERROR" {
		fmt.Println("api error")
	}
}
