package services

import (
	"encoding/json"
	"fmt"

	"github.com/streamco/bitmovin-go/bitmovin"
	"github.com/streamco/bitmovin-go/bitmovintypes"
	"github.com/streamco/bitmovin-go/models"
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

// not part of the original bitmovin API
func (s *EncodingService) AddIngestStream(encodingID string, name string, inputID string, inputPath string,
	selectionMode bitmovintypes.SelectionMode, position int) (*models.StreamResponse, error) {

	b := []byte(fmt.Sprintf(`{"name" : %q, "inputId" : %q, "inputPath": %q, "selectionMode": %q, "position": %d}`,
		name, inputID, inputPath, selectionMode, position))

	path := EncodingEndpoint + "/" + encodingID + "/" + "input-streams" + "/" + "ingest"
	o, err := s.RestService.Create(path, b)
	if err != nil {
		return nil, err
	}
	var r models.StreamResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

// not part of the original bitmovin API
func (s *EncodingService) AddDolbyVisionSidecarStream(encodingID string,
	inputID string, inputPath string) (*models.StreamResponse, error) {

	b := []byte(fmt.Sprintf(`{"type" : "SIDECAR_DOLBY_VISION_METADATA", "inputId" : %q, "inputPath": %q}`, inputID, inputPath))

	path := EncodingEndpoint + "/" + encodingID + "/" + "input-streams" + "/" + "sidecar" + "/" + "dolby-vision-metadata-ingest"
	o, err := s.RestService.Create(path, b)
	if err != nil {
		return nil, err
	}
	var r models.StreamResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

// not part of the original bitmovin API
func (s *EncodingService) MarkStreamDolbyVision(encodingID string, streamID string, dvMetadataStreamID string) error {

	var b []byte

	if dvMetadataStreamID == "" {
		b = []byte(`{"profile" : "DVHE_05", "metadataSource" : "EMBEDDED"}`)
	} else {
		b = []byte(fmt.Sprintf(`{"profile" : "DVHE_05", "metadataSource" : "INPUT_STREAM", "metadataInputStreamId": %q}`,
			dvMetadataStreamID))
	}

	path := EncodingEndpoint + "/" + encodingID + "/" + "streams" + "/" + streamID + "/" + "hdr" + "/" + "dolby-vision"
	_, err := s.RestService.Create(path, b)
	return err
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) AddKeyframe(encodingID string, a *models.Keyframe) (*models.KeyframeResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	path := EncodingEndpoint + "/" + encodingID + "/" + "keyframes"
	o, err := s.RestService.Create(path, b)
	if err != nil {
		return nil, err
	}
	var r models.KeyframeResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) RetrieveKeyframe(encodingID string, keyframeID string) (*models.KeyframeResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "keyframes" + "/" + keyframeID
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.KeyframeResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) DeleteKeyframe(encodingID string, keyframeID string) (*models.KeyframeResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "keyframes" + "/" + keyframeID
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.KeyframeResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) ListKeyframes(encodingID string, offset int64, limit int64) (*models.KeyframeListResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "keyframes"
	o, err := s.RestService.List(path, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.KeyframeListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) RetrieveStreamInputData(encodingID string, streamID string) (*models.StreamInputResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "streams" + "/" + streamID + "/input"
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.StreamInputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) RetrieveMP4MuxingInformation(encodingID string, mp4MuxingID string) (*models.MP4MuxingInformationResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/mp4" + "/" + mp4MuxingID + "/information"
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.MP4MuxingInformationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) AddProgressiveMOVMuxing(encodingID string, a *models.ProgressiveMOVMuxing) (*models.ProgressiveMOVMuxingResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/progressive-mov"
	o, err := s.RestService.Create(path, b)
	if err != nil {
		return nil, err
	}
	var r models.ProgressiveMOVMuxingResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) RetrieveProgressiveMOVMuxing(encodingID string, progressiveMOVID string) (*models.ProgressiveMOVMuxingResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/progressive-mov" + "/" + progressiveMOVID
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.ProgressiveMOVMuxingResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) DeleteProgressiveMOVMuxing(encodingID string, progressiveMOVID string) (*models.ProgressiveMOVMuxingResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/progressive-mov" + "/" + progressiveMOVID
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.ProgressiveMOVMuxingResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) ListProgressiveMOVMuxing(encodingID string, offset int64, limit int64) (*models.ProgressiveMOVMuxingListResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/progressive-mov"
	o, err := s.RestService.List(path, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.ProgressiveMOVMuxingListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) RetrieveProgressiveMOVMuxingCustomData(encodingID string, progressiveWebMID string, offset int64, limit int64) (*models.CustomDataResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/progressive-mov" + "/" + progressiveWebMID
	o, err := s.RestService.RetrieveCustomData(path)
	if err != nil {
		return nil, err
	}
	var r models.CustomDataResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) RetrieveProgressiveMOVMuxingInformation(encodingID string, progressiveWebMID string) (*models.ProgressiveMOVMuxingInformationResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/progressive-mov" + "/" + progressiveWebMID + "/information"
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.ProgressiveMOVMuxingInformationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) AddProgressiveTSMuxing(encodingID string, a *models.ProgressiveTSMuxing) (*models.ProgressiveTSMuxingResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/progressive-ts"
	o, err := s.RestService.Create(path, b)
	if err != nil {
		return nil, err
	}
	var r models.ProgressiveTSMuxingResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) RetrieveProgressiveTSMuxing(encodingID string, progressiveTSID string) (*models.ProgressiveTSMuxingResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/progressive-ts" + "/" + progressiveTSID
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.ProgressiveTSMuxingResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) DeleteProgressiveTSMuxing(encodingID string, progressiveTSID string) (*models.ProgressiveTSMuxingResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/progressive-ts" + "/" + progressiveTSID
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.ProgressiveTSMuxingResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) ListProgressiveTSMuxing(encodingID string, offset int64, limit int64) (*models.ProgressiveTSMuxingListResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/progressive-ts"
	o, err := s.RestService.List(path, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.ProgressiveTSMuxingListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) RetrieveProgressiveTSMuxingCustomData(encodingID string, progressiveTSID string, offset int64, limit int64) (*models.CustomDataResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/progressive-ts" + "/" + progressiveTSID
	o, err := s.RestService.RetrieveCustomData(path)
	if err != nil {
		return nil, err
	}
	var r models.CustomDataResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) RetrieveProgressiveTSMuxingInformation(encodingID string, progressiveTSID string) (*models.ProgressiveTSMuxingInformationResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/progressive-ts" + "/" + progressiveTSID + "/information"
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.ProgressiveTSMuxingInformationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) AddProgressiveWebMMuxing(encodingID string, a *models.ProgressiveWebMMuxing) (*models.ProgressiveWebMMuxingResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/progressive-webm"
	o, err := s.RestService.Create(path, b)
	if err != nil {
		return nil, err
	}
	var r models.ProgressiveWebMMuxingResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) RetrieveProgressiveWebMMuxing(encodingID string, progressiveWebMID string) (*models.ProgressiveWebMMuxingResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/progressive-webm" + "/" + progressiveWebMID
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.ProgressiveWebMMuxingResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) DeleteProgressiveWebMMuxing(encodingID string, progressiveWebMID string) (*models.ProgressiveWebMMuxingResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/progressive-webm" + "/" + progressiveWebMID
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.ProgressiveWebMMuxingResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) ListProgressiveWebMMuxing(encodingID string, offset int64, limit int64) (*models.ProgressiveWebMMuxingListResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/progressive-webm"
	o, err := s.RestService.List(path, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.ProgressiveWebMMuxingListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) RetrieveProgressiveWebMMuxingCustomData(encodingID string, progressiveWebMID string, offset int64, limit int64) (*models.CustomDataResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/progressive-webm" + "/" + progressiveWebMID
	o, err := s.RestService.RetrieveCustomData(path)
	if err != nil {
		return nil, err
	}
	var r models.CustomDataResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) RetrieveProgressiveWebMMuxingInformation(encodingID string, progressiveWebMID string) (*models.ProgressiveWebMMuxingInformationResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + "muxings/progressive-webm" + "/" + progressiveWebMID + "/information"
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.ProgressiveWebMMuxingInformationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) StartWithOptions(encodingID string, startOptions *models.StartOptions) (*models.StartStopResponse, error) {
	b, err := json.Marshal(*startOptions)
	if err != nil {
		return nil, err
	}
	path := EncodingEndpoint + "/" + encodingID + "/start"
	o, err := s.RestService.Create(path, b)
	if err != nil {
		return nil, err
	}
	var r models.StartStopResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) Reschedule(encodingID string, rescheduleEncoding *models.RescheduleEncoding) (*models.StartStopResponse, error) {
	b, err := json.Marshal(*rescheduleEncoding)
	if err != nil {
		return nil, err
	}
	path := EncodingEndpoint + "/" + encodingID + "/reschedule"
	o, err := s.RestService.Create(path, b)
	if err != nil {
		return nil, err
	}
	var r models.StartStopResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
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
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) AddThumbnail(encodingID, streamID string, thumb *models.Thumbnail) (*models.ThumbnailResponse, error) {
	payload, err := json.Marshal(thumb)
	if err != nil {
		return nil, err
	}

	path := EncodingEndpoint + "/" + encodingID + "/streams/" + streamID + "/thumbnails"
	resp, err := s.RestService.Create(path, payload)
	if err != nil {
		return nil, err
	}
	var r models.ThumbnailResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) ListThumbnails(encodingID, streamID string, offset, limit int64) (*models.ThumbnailListResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + streamID + "/thumbnails"
	o, err := s.RestService.List(path, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.ThumbnailListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) AddSprites(encodingID, streamID string, sprite *models.Sprite) (*models.SpriteResponse, error) {
	payload, err := json.Marshal(sprite)
	if err != nil {
		return nil, err
	}

	path := EncodingEndpoint + "/" + encodingID + "/streams/" + streamID + "/sprites"
	resp, err := s.RestService.Create(path, payload)
	if err != nil {
		return nil, err
	}
	var r models.SpriteResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *EncodingService) ListSprites(encodingID, streamID string, offset, limit int64) (*models.SpriteListResponse, error) {
	path := EncodingEndpoint + "/" + encodingID + "/" + streamID + "/sprites"
	o, err := s.RestService.List(path, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.SpriteListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response due to error %q. Original text was: %s", err, string(o))
	}
	return &r, nil
}

func (s *EncodingService) AddFilter(encodingID, streamID, filterID string, position int64) (*models.AddFilterResponse, error) {
	path := fmt.Sprintf("%s/%s/streams/%s/filters", EncodingEndpoint, encodingID, streamID)
	f := &models.AddFilter{ID: filterID, Position: &position}
	payload, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}

	resp, err := s.RestService.Create(path, payload)
	if err != nil {
		return nil, err
	}
	var r models.AddFilterResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
