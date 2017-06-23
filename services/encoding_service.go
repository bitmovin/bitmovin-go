package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type EncodingService struct {
	RestService *RestService
}

const (
	EncodingEndpoint string = "encoding/encodings"
)

func NewEncodingService(bitmovin *bitmovin.Bitmovin) *EncodingService {
	r := NewRestService(bitmovin)

	return &EncodingService{RestService: r}
}

func (s *EncodingService) Create(a *models.Encoding) (*models.EncodingResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	o, err := s.RestService.Create(EncodingEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.EncodingResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) Retrieve(id string) (*models.EncodingResponse, error) {
	path := EncodingEndpoint + "/" + id
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.EncodingResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) Delete(id string) (*models.EncodingResponse, error) {
	path := EncodingEndpoint + "/" + id
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.EncodingResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) List(offset int64, limit int64) (*models.EncodingListResponse, error) {
	o, err := s.RestService.List(EncodingEndpoint, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.EncodingListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) RetrieveCustomData(id string) (*models.CustomDataResponse, error) {
	path := EncodingEndpoint + "/" + id
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

func (s *EncodingService) AddStream(encodingID string, a *models.Stream) (*models.StreamResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	path := EncodingEndpoint + "/" + encodingID + "/" + "streams"
	o, err := s.RestService.Create(path, b)
	if err != nil {
		return nil, err
	}
	var r models.StreamResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) RetrieveStream(encodingID string, streamID string) (*models.StreamResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "streams" + "/" + streamID
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.StreamResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) DeleteStream(encodingID string, streamID string) (*models.StreamResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "streams" + "/" + streamID
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.StreamResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) ListStream(encodingID string, offset int64, limit int64) (*models.StreamListResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "streams"
	o, err := s.RestService.List(path, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.StreamListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) RetrieveStreamCustomData(encodingID string, streamID string, offset int64, limit int64) (*models.CustomDataResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "streams" + "/" + streamID
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

func (s *EncodingService) AddFMP4Muxing(encodingID string, a *models.FMP4Muxing) (*models.FMP4MuxingResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/fmp4"
	o, err := s.RestService.Create(path, b)
	if err != nil {
		return nil, err
	}
	var r models.FMP4MuxingResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) RetrieveFMP4Muxing(encodingID string, fmp4ID string) (*models.FMP4MuxingResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/fmp4" + "/" + fmp4ID
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.FMP4MuxingResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) DeleteFMP4Muxing(encodingID string, fmp4ID string) (*models.FMP4MuxingResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/fmp4" + "/" + fmp4ID
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.FMP4MuxingResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) ListFMP4Muxing(encodingID string, offset int64, limit int64) (*models.FMP4MuxingListResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/fmp4"
	o, err := s.RestService.List(path, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.FMP4MuxingListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) RetrieveFMP4MuxingCustomData(encodingID string, fmp4ID string, offset int64, limit int64) (*models.CustomDataResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/fmp4" + "/" + fmp4ID
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

func (s *EncodingService) AddTSMuxing(encodingID string, a *models.TSMuxing) (*models.TSMuxingResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/ts"
	o, err := s.RestService.Create(path, b)
	if err != nil {
		return nil, err
	}
	var r models.TSMuxingResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) RetrieveTSMuxing(encodingID string, tsID string) (*models.TSMuxingResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/ts" + "/" + tsID
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.TSMuxingResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) DeleteTSMuxing(encodingID string, tsID string) (*models.TSMuxingResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/ts" + "/" + tsID
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.TSMuxingResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) ListTSMuxing(encodingID string, offset int64, limit int64) (*models.TSMuxingListResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/ts"
	o, err := s.RestService.List(path, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.TSMuxingListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) RetrieveTSMuxingCustomData(encodingID string, tsID string, offset int64, limit int64) (*models.CustomDataResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/ts" + "/" + tsID
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

func (s *EncodingService) AddMP4Muxing(encodingID string, a *models.MP4Muxing) (*models.MP4MuxingResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/mp4"
	o, err := s.RestService.Create(path, b)
	if err != nil {
		return nil, err
	}
	var r models.MP4MuxingResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) RetrieveMP4Muxing(encodingID string, mp4ID string) (*models.MP4MuxingResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/mp4" + "/" + mp4ID
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.MP4MuxingResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) DeleteMP4Muxing(encodingID string, mp4ID string) (*models.MP4MuxingResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/mp4" + "/" + mp4ID
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.MP4MuxingResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) ListMP4Muxing(encodingID string, offset int64, limit int64) (*models.MP4MuxingListResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/mp4"
	o, err := s.RestService.List(path, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.MP4MuxingListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) RetrieveMP4MuxingCustomData(encodingID string, mp4ID string, offset int64, limit int64) (*models.CustomDataResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/mp4" + "/" + mp4ID
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

func (s *EncodingService) Start(encodingID string) (*models.StartStopResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/start"
	o, err := s.RestService.Create(path, nil)
	if err != nil {
		return nil, err
	}
	var r models.StartStopResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// Stop and Start use the same model
func (s *EncodingService) Stop(encodingID string) (*models.StartStopResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/stop"
	o, err := s.RestService.Create(path, nil)
	if err != nil {
		return nil, err
	}
	var r models.StartStopResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) RetrieveStatus(encodingID string) (*models.StatusResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/status"
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.StatusResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) StartLive(encodingID string, a *models.LiveStreamConfiguration) (*models.StartStopResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	path := EncodingEndpoint + "/" + encodingID + "/live/start"
	o, err := s.RestService.Create(path, b)
	if err != nil {
		return nil, err
	}
	var r models.StartStopResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) StopLive(encodingID string) (*models.StartStopResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/live/stop"
	o, err := s.RestService.Create(path, nil)
	if err != nil {
		return nil, err
	}
	var r models.StartStopResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) RetrieveLiveStatus(encodingID string) (*models.LiveStatusResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/live"
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.LiveStatusResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
