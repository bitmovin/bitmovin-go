package services

import (
	"encoding/json"
	"testing"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

const apiKey = "INSERT_API_KEY"

func TestCreateError(t *testing.T) {
	bitmovin := bitmovin.NewBitmovinDefaultTimeout(apiKey, "https://api.bitmovin.com/v1/")
	svc := NewRestService(bitmovin)
	gcsInput := &models.GCSInput{
		AccessKey:  stringToPtr(""),
		SecretKey:  stringToPtr(""),
		BucketName: stringToPtr(""),
	}
	json, _ := json.Marshal(*gcsInput)
	_, err := svc.Create(`encoding/inputs/gcs`, json)

	if err == nil {
		t.Fatal("Expected to receive error")
	}
	if err.Error() != "ERROR 1000: One or more fields are not present or invalid" {
		t.Fatalf("Expected error message - got %s", err.Error())
	}
}

func stringToPtr(s string) *string {
	return &s
}
