package services

import (
	"encoding/json"

	"github.com/streamco/bitmovin-go/bitmovin"
	"github.com/streamco/bitmovin-go/models"
)

type AC3CodecConfigurationService struct {
	RestService *RestService
}

const AC3CodecConfigurationEndpoint string = "encoding/configurations/audio/ac3"

func NewAC3CodecConfigurationService(bitmovin *bitmovin.Bitmovin) *AC3CodecConfigurationService {
	r := NewRestService(bitmovin)
	return &AC3CodecConfigurationService{RestService: r}
}

func (s *AC3CodecConfigurationService) Create(a *models.AC3CodecConfiguration) (*models.AC3CodecConfigurationResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	o, err := s.RestService.Create(AC3CodecConfigurationEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.AC3CodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
