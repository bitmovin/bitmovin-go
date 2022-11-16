package services

import (
	"encoding/json"

	"github.com/streamco/bitmovin-go/bitmovin"
	"github.com/streamco/bitmovin-go/models"
)

type AV1CodecConfigurationService struct {
	RestService *RestService
}

const (
	AV1CodecConfigurationEndpoint string = "encoding/configurations/video/av1"
)

func NewAV1CodecConfigurationService(bitmovin *bitmovin.Bitmovin) *AV1CodecConfigurationService {
	r := NewRestService(bitmovin)

	return &AV1CodecConfigurationService{RestService: r}
}

func (s *AV1CodecConfigurationService) Create(a *models.AV1CodecConfiguration) (*models.AV1CodecConfigurationResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	o, err := s.RestService.Create(AV1CodecConfigurationEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.AV1CodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *AV1CodecConfigurationService) Retrieve(id string) (*models.AV1CodecConfigurationResponse, error) {
	path := AV1CodecConfigurationEndpoint + "/" + id
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.AV1CodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *AV1CodecConfigurationService) Delete(id string) (*models.AV1CodecConfigurationResponse, error) {
	path := AV1CodecConfigurationEndpoint + "/" + id
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.AV1CodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *AV1CodecConfigurationService) List(offset int64, limit int64) (*models.AV1CodecConfigurationListResponse, error) {
	o, err := s.RestService.List(AV1CodecConfigurationEndpoint, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.AV1CodecConfigurationListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *AV1CodecConfigurationService) RetrieveCustomData(id string) (*models.CustomDataResponse, error) {
	path := AV1CodecConfigurationEndpoint + "/" + id
	o, err := s.RestService.RetrieveCustomData(path)
	if err != nil {
		return nil, err
	}
	var r models.CustomDataResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
