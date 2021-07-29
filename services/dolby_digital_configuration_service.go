package services

import (
	"encoding/json"

	"github.com/streamco/bitmovin-go/bitmovin"
	"github.com/streamco/bitmovin-go/models"
)

type DolbyDigitalCodecConfigurationService struct {
	RestService *RestService
}

const DolbyDigitalCodecConfigurationServiceEndpoint string = "encoding/configurations/audio/dolby-digital"

func NewDolbyDigitalCodecConfigurationService(bitmovin *bitmovin.Bitmovin) *DolbyDigitalCodecConfigurationService {
	r := NewRestService(bitmovin)
	return &DolbyDigitalCodecConfigurationService{RestService: r}
}

func (s *DolbyDigitalCodecConfigurationService) Create(a *models.DolbyDigitalAudioConfiguration) (*models.DolbyDigitalAudioConfigurationResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	o, err := s.RestService.Create(DolbyDigitalCodecConfigurationServiceEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.DolbyDigitalAudioConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
