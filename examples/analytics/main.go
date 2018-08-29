package main

import (
	"log"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
	"github.com/bitmovin/bitmovin-go/services"
)

func main() {
	bitmovin := bitmovin.NewBitmovinDefaultTimeout("YOUR API KEY", "https://api.bitmovin.com/v1/")
	as := services.NewAnalyticsService(bitmovin)
	query := models.Query{
		Dimension: "PLAYER_STARTUPTIME",
		Start:     "2018-06-25T20:09:23.69Z",
		End:       "2018-08-25T20:09:23.69Z",
	}
	res, err := as.Avg(&query)

	if err != nil {
		log.Printf("could not query %+v: %v", query, err)
	}

	log.Printf("%+v", res)

}
