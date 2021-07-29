package services

import (
	"encoding/json"

	"github.com/streamco/bitmovin-go/bitmovin"
	"github.com/streamco/bitmovin-go/models"
)

type DolbyDigitalPlusCodecConfigurationService struct {
	RestService *RestService
}

const DolbyDigitalPlusCodecConfigurationServiceEndpoint string = "encoding/configurations/audio/dolby-digital"

func NewDolbyDigitalPlusCodecConfigurationService(bitmovin *bitmovin.Bitmovin) *DolbyDigitalPlusCodecConfigurationService {
	r := NewRestService(bitmovin)
	return &DolbyDigitalPlusCodecConfigurationService{RestService: r}
}

func (s *DolbyDigitalPlusCodecConfigurationService) Create(a *models.DolbyDigitalPlusAudioConfiguration) (*models.DolbyDigitalPlusAudioConfigurationResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	o, err := s.RestService.Create(DolbyDigitalPlusCodecConfigurationServiceEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.DolbyDigitalPlusAudioConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
