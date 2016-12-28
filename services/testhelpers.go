package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func RespondWithFileHandler(filename string, httpStatus int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rawResponse, err := ioutil.ReadFile(filename)
		if err != nil {
			panic("couldn't load fixture")
		}
		w.WriteHeader(httpStatus)
		w.Write(rawResponse)
	}
}

func LoadJSONFileIntoStruct(filename string, i interface{}) error {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw, i)
	if err != nil {
		return err
	}
	return nil
}
