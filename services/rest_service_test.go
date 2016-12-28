package services

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/gorilla/mux"
)

func TestNewRestService(t *testing.T) {
	b := bitmovin.NewBitmovinDefaultTimeout("apikey", "someURL")
	r := NewRestService(b)
	if r == nil {
		t.Errorf("RestService should not be created nil")
	}
}

func TestCreatePasses(t *testing.T) {
	rawInput, err := ioutil.ReadFile("fixtures/h264_codec_configuration_create.json")
	if err != nil {
		t.Errorf("couldn't load fixture file")
	}
	svr := httptest.NewServer(restHandlers())
	b := bitmovin.NewBitmovinDefaultTimeout("apikey", svr.URL)
	r := NewRestService(b)
	_, err = r.Create("/create/path/passes", rawInput)
	if err != nil {
		t.Errorf("Create Fails when it should pass")
	}
}

func TestRetrievePasses(t *testing.T) {
	svr := httptest.NewServer(restHandlers())
	b := bitmovin.NewBitmovinDefaultTimeout("apikey", svr.URL)
	r := NewRestService(b)
	_, err := r.Retrieve("/retrieve/path/passes/thisismyid")
	if err != nil {
		t.Errorf("Retrieve Fails when it should pass")
	}
}

func TestDeletePasses(t *testing.T) {
	svr := httptest.NewServer(restHandlers())
	b := bitmovin.NewBitmovinDefaultTimeout("apikey", svr.URL)
	r := NewRestService(b)
	_, err := r.Delete("/retrieve/path/passes/thisismyid")
	if err != nil {
		t.Errorf("Retrieve Fails when it should pass")
	}
}

func TestCreateRetrieveDeleteFailsOnUnparsableURL(t *testing.T) {
	rawInput, err := ioutil.ReadFile("fixtures/h264_codec_configuration_create.json")
	if err != nil {
		t.Errorf("couldn't load fixture file")
	}
	b := bitmovin.NewBitmovinDefaultTimeout("apikey", "http://192.168.0.%31/")
	r := NewRestService(b)
	_, err = r.Create("/create/path/passes", rawInput)
	if err == nil {
		t.Errorf("Create should error on unparsable URL")
	}
	_, err = r.Retrieve("/retrieve/path/passes/sdfjf")
	if err == nil {
		t.Errorf("Retrieve should error on unparsable URL")
	}
	_, err = r.Delete("/delete/path/passes/sdfjf")
	if err == nil {
		t.Errorf("Delete should error on unparsable URL")
	}
}

func TestCreateRetrieveDeleteFailsOn404(t *testing.T) {
	rawInput, err := ioutil.ReadFile("fixtures/h264_codec_configuration_create.json")
	if err != nil {
		t.Errorf("couldn't load fixture file")
	}
	b := bitmovin.NewBitmovin("apikey", "http://192.168.123.54/", 1)
	r := NewRestService(b)
	_, err = r.Create("/create/path/passes", rawInput)
	if err == nil {
		t.Errorf("Create should error on unparsable URL")
	}
	_, err = r.Retrieve("/retrieve/path/passes/sdfjf")
	if err == nil {
		t.Errorf("Retrieve should error on unparsable URL")
	}
	_, err = r.Delete("/delete/path/passes/sdfjf")
	if err == nil {
		t.Errorf("Delete should error on unparsable URL")
	}
}

func restHandlers() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("create/path/passes", RespondWithFileHandler("fixtures/h264_codec_configuration_create_response.json", http.StatusCreated)).Methods("POST")
	r.HandleFunc("retrieve/path/passes/{id}", RespondWithFileHandler("fixtures/h264_codec_configuration_retrieve_response.json", http.StatusOK)).Methods("GET")
	r.HandleFunc("delete/path/passes/{id}", RespondWithFileHandler("fixtures/h264_codec_configuration_delete_response.json", http.StatusOK)).Methods("DELETE")

	return r
}

// func respondWithFileHandler(filename string, httpStatus int) func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		rawResponse, err := ioutil.ReadFile("fixtures/h264_codec_configuration_create_response.json")
// 		if err != nil {
// 			panic("couldn't load fixture")
// 		}
// 		w.WriteHeader(httpStatus)
// 		w.Write(rawResponse)
// 	}
// }
