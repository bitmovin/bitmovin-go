package services

import (
	"reflect"
	"testing"

	"github.com/bitmovin/bitmovin-go/bitmovin"
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
