package bitmovin

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestNewBitmovinDefaultTimeout(t *testing.T) {
	apiKey := "apiKey"
	baseURL := "baseURL"
	a := &Bitmovin{
		HTTPClient: &http.Client{
			Timeout: time.Second * 5,
		},
		APIKey:     &apiKey,
		APIBaseURL: &baseURL,
	}
	b := NewBitmovinDefaultTimeout(apiKey, baseURL)
	if !reflect.DeepEqual(a, b) {
		t.Errorf("Structs should be equivalent")
	}
}
