package services

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
	"github.com/gorilla/mux"
)

func TestNewH264CodecConfigurationService(t *testing.T) {
	bitmovin := bitmovin.NewBitmovinDefaultTimeout("apikey", "someURL")
	r := NewRestService(bitmovin)
	a := NewH264CodecConfigurationService(bitmovin)
	b := &H264CodecConfigurationService{
		RestService: r,
	}
	if !reflect.DeepEqual(a, b) {
		t.Errorf("Structs should be equivalent")
	}
}

func TestH264ConfigurationServiceCreatePasses(t *testing.T) {
	var input models.H264CodecConfiguration
	err := LoadJSONFileIntoStruct("fixtures/h264_codec_configuration_create.json", &input)
	if err != nil {
		t.Errorf(err.Error())
	}
	svr := httptest.NewServer(h264PassingHandlers())
	b := bitmovin.NewBitmovinDefaultTimeout("apikey", svr.URL)
	h := NewH264CodecConfigurationService(b)
	c, err := h.Create(&input)
	if err != nil {
		t.Errorf(err.Error())
	}
	if c.Data.Result.ID == "" {
		t.Errorf("Result should have an ID")
	}
}

func TestH264ConfigurationServiceRetrievePasses(t *testing.T) {
	svr := httptest.NewServer(h264PassingHandlers())
	b := bitmovin.NewBitmovinDefaultTimeout("apikey", svr.URL)
	h := NewH264CodecConfigurationService(b)
	c, err := h.Retrieve("thisisanid")
	if err != nil {
		t.Errorf(err.Error())
	}
	if c.Data.Result.ID == "" {
		t.Errorf("Result should have an ID")
	}
}

func TestH264ConfigurationServiceDeletePasses(t *testing.T) {
	svr := httptest.NewServer(h264PassingHandlers())
	b := bitmovin.NewBitmovinDefaultTimeout("apikey", svr.URL)
	h := NewH264CodecConfigurationService(b)
	c, err := h.Delete("thisisanid")
	if err != nil {
		t.Errorf(err.Error())
	}
	if c.Data.Result.ID == "" {
		t.Errorf("Result should have an ID")
	}
}

func TestH264ConfigurationServiceCreateRetrieveDeleteFailsIfRestServiceFails(t *testing.T) {
	var input models.H264CodecConfiguration
	err := LoadJSONFileIntoStruct("fixtures/h264_codec_configuration_create.json", &input)
	if err != nil {
		t.Errorf(err.Error())
	}
	b := bitmovin.NewBitmovin("apikey", "http://192.168.123.54/", 1)
	h := NewH264CodecConfigurationService(b)
	_, err = h.Create(&input)
	if err == nil {
		t.Errorf("Rest Service should have 404 and this Create should have errored.")
	}
	_, err = h.Retrieve("id")
	if err == nil {
		t.Errorf("Rest Service should have 404 and this Retrieve should have errored.")
	}
	_, err = h.Delete("id")
	if err == nil {
		t.Errorf("Rest Service should have 404 and this Delete should have errored.")
	}
}

func TestH264ConfigurationServiceCreateRetrieveDeleteFailsIfProperJSONNotReturned(t *testing.T) {
	var input models.H264CodecConfiguration
	err := LoadJSONFileIntoStruct("fixtures/h264_codec_configuration_create.json", &input)
	if err != nil {
		t.Errorf(err.Error())
	}
	svr := httptest.NewServer(bogusH264JSONHandlers())
	b := bitmovin.NewBitmovinDefaultTimeout("apikey", svr.URL)
	h := NewH264CodecConfigurationService(b)
	_, err = h.Create(&input)
	if err == nil {
		t.Errorf("JSON was not returned, and unmarshalling should have failed and Create should have failed.")
	}
	_, err = h.Retrieve("id")
	if err == nil {
		t.Errorf("JSON was not returned, and unmarshalling should have failed and Retrieve should have failed.")
	}
	_, err = h.Delete("id")
	if err == nil {
		t.Errorf("JSON was not returned, and unmarshalling should have failed and Delete should have failed.")
	}
}

func h264PassingHandlers() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc(H264CodecConfigurationEndpoint, RespondWithFileHandler("fixtures/h264_codec_configuration_create_response.json", http.StatusCreated)).Methods("POST")
	r.HandleFunc(H264CodecConfigurationEndpoint+"/{id}", RespondWithFileHandler("fixtures/h264_codec_configuration_retrieve_response.json", http.StatusOK)).Methods("GET")
	r.HandleFunc(H264CodecConfigurationEndpoint+"/{id}", RespondWithFileHandler("fixtures/h264_codec_configuration_delete_response.json", http.StatusOK)).Methods("DELETE")

	return r
}

func bogusH264JSONHandlers() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc(H264CodecConfigurationEndpoint, RespondWithFileHandler("fixtures/notjson.notjson", http.StatusCreated)).Methods("POST")
	r.HandleFunc(H264CodecConfigurationEndpoint+"/{id}", RespondWithFileHandler("fixtures/notjson.notjson", http.StatusOK)).Methods("GET")
	r.HandleFunc(H264CodecConfigurationEndpoint+"/{id}", RespondWithFileHandler("fixtures/notjson.notjson", http.StatusOK)).Methods("DELETE")

	return r
}
