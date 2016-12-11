package services

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/gorilla/mux"
)

func TestCreatePasses(t *testing.T) {
	rawInput, err := ioutil.ReadFile("fixtures/h264_codec_configuration_create.json")
	if err != nil {
		t.Errorf("couldn't load fixture file")
	}
	svr := httptest.NewServer(handlers())
	b := bitmovin.NewBitmovinDefaultTimeout("apikey", svr.URL)
	r := NewRestService(b)
	_, err = r.Create("/create/path/passes", rawInput)
	if err != nil {
		t.Errorf("Create Fails when it should pass")
	}
}

func TestCreateFailsOnUnparsableURL(t *testing.T) {
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
}

func handlers() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/create/path/passes", respondWithFileHandler("fixtures/h264_codec_configuration_create_response.json", http.StatusCreated)).Methods("POST")
	// r.HandleFunc("/create/path/fails", respondWithFileHandler("fixtures/h264_codec_configuration_create.json", http.StatusCreated)).Methods("POST")

	// r.HandleFunc("/create/path/passes", createPassesHandler).Methods("POST")

	return r
}

func createPassesHandler(w http.ResponseWriter, r *http.Request) {
	rawResponse, err := ioutil.ReadFile("fixtures/h264_codec_configuration_create_response.json")
	if err != nil {
		panic("couldn't load fixture")
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(rawResponse)
}

func respondWithFileHandler(filename string, httpStatus int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rawResponse, err := ioutil.ReadFile("fixtures/h264_codec_configuration_create_response.json")
		if err != nil {
			panic("couldn't load fixture")
		}
		w.WriteHeader(httpStatus)
		w.Write(rawResponse)
	}
}

// func createFailsHandler(w http.ResponseWriter, r *http.Request) {
// }
