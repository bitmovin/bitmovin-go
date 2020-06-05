package services

import (
	"encoding/json"

	"github.com/streamco/bitmovin-go/bitmovin"
	"github.com/streamco/bitmovin-go/models"
)

type DolbyAtmosCodecConfigurationService struct {
	RestService *RestService
}

const DolbyAtmosCodecConfigurationEndpoint string = "encoding/configurations/audio/dolby-atmos"

func NewDolbyAtmosCodecConfigurationService(bitmovin *bitmovin.Bitmovin) *DolbyAtmosCodecConfigurationService {
	r := NewRestService(bitmovin)
	return &DolbyAtmosCodecConfigurationService{RestService: r}
}

func (s *DolbyAtmosCodecConfigurationService) Create(a *models.DolbyAtmosCodecConfiguration) (*models.DolbyAtmosCodecConfigurationResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	o, err := s.RestService.Create(DolbyAtmosCodecConfigurationEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.DolbyAtmosCodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
