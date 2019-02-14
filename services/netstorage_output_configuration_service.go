package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type AkamaiNetstorageOutputService struct {
	RestService *RestService
}

const (
	AkamaiNetstorageOutputEndpoint string = "encoding/outputs/akamai-netstorage"
)

func NewAkamaiNetstorageOutputService(bitmovin *bitmovin.Bitmovin) *AkamaiNetstorageOutputService {
	r := NewRestService(bitmovin)

	return &AkamaiNetstorageOutputService{RestService: r}
}

func (s *AkamaiNetstorageOutputService) Create(a *models.AkamaiNetstorageOutput) (*models.AkamaiNetstorageOutputResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	o, err := s.RestService.Create(AkamaiNetstorageOutputEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.AkamaiNetstorageOutputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *AkamaiNetstorageOutputService) Retrieve(id string) (*models.AkamaiNetstorageOutputResponse, error) {
	path := AkamaiNetstorageOutputEndpoint + "/" + id
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.AkamaiNetstorageOutputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *AkamaiNetstorageOutputService) Delete(id string) (*models.AkamaiNetstorageOutputResponse, error) {
	path := AkamaiNetstorageOutputEndpoint + "/" + id
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.AkamaiNetstorageOutputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *AkamaiNetstorageOutputService) List(offset int64, limit int64) (*models.AkamaiNetstorageOutputListResponse, error) {
	o, err := s.RestService.List(AkamaiNetstorageOutputEndpoint, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.AkamaiNetstorageOutputListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *AkamaiNetstorageOutputService) RetrieveCustomData(id string) (*models.CustomDataResponse, error) {
	path := NetstorageOutputEndpoint + "/" + id
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
