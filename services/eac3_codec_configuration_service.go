package services

import (
	"encoding/json"

	"github.com/streamco/bitmovin-go/bitmovin"
	"github.com/streamco/bitmovin-go/models"
)

type EAC3CodecConfigurationService struct {
	RestService *RestService
}

const EAC3CodecConfigurationEndpoint string = "encoding/configurations/audio/eac3"

func NewEAC3CodecConfigurationService(bitmovin *bitmovin.Bitmovin) *EAC3CodecConfigurationService {
	r := NewRestService(bitmovin)
	return &EAC3CodecConfigurationService{RestService: r}
}

func (s *EAC3CodecConfigurationService) Create(a *models.EAC3CodecConfiguration) (*models.EAC3CodecConfigurationResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	o, err := s.RestService.Create(EAC3CodecConfigurationEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.EAC3CodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
