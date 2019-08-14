package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/streamco/bitmovin-go/bitmovin"
	"github.com/streamco/bitmovin-go/models"
)

type RestService struct {
	Bitmovin *bitmovin.Bitmovin
}

func NewRestService(bitmovin *bitmovin.Bitmovin) *RestService {
	return &RestService{
		Bitmovin: bitmovin,
	}
}

func (r *RestService) makeRequest(method, url string, input []byte, retries int) ([]byte, error) {
	giveUp := retries <= 1
	req, err := http.NewRequest(method, url, bytes.NewBuffer(input))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", *r.Bitmovin.APIKey)
	if r.Bitmovin.OrganizationID != nil {
		req.Header.Set("X-Tenant-Org-Id", *r.Bitmovin.OrganizationID)
	}
	req.Header.Set("X-Api-Client", ClientName)
	req.Header.Set("X-Api-Client-Version", Version)

	if os.Getenv("DUMP_TRAFFIC") != "" {
		b, _ := httputil.DumpRequest(req, true)
		println(string(b))
	}
	resp, err := r.Bitmovin.HTTPClient.Do(req)
	if err != nil {
		if giveUp {
			log.Printf("makeRequest error: %s: %s %s", err, method, url)
			log.Println("giving up", method, url)
			return nil, err
		}
		log.Printf("makeRequest error: %s: %s %s", err, method, url)
		log.Println("retrying", method, url)
		time.Sleep(1 * time.Second)
		return r.makeRequest(method, url, input, retries-1)
	}
	if os.Getenv("DUMP_TRAFFIC") != "" {
		b, _ := httputil.DumpResponse(resp, true)
		println(string(b))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 399 {
		if giveUp {
			log.Printf("makeRequest error: %s: %s %s: statusCode: %d", err, method, url, resp.StatusCode)
			log.Println("giving up", method, url)
			return nil, formatError(body)
		}
		log.Printf("makeRequest error: %s: %s %s", err, method, url)
		log.Println("retrying", method, url)
		time.Sleep(1 * time.Second)
		return r.makeRequest(method, url, input, retries-1)
	}
	return body, nil
}

func (r *RestService) Create(relativeURL string, input []byte) ([]byte, error) {
	fullURL := *r.Bitmovin.APIBaseURL + relativeURL
	_, err := url.Parse(fullURL)
	if err != nil {
		return nil, err
	}
	return r.makeRequest("POST", fullURL, input, 2)
}

func (r *RestService) Retrieve(relativeURL string) ([]byte, error) {
	fullURL := *r.Bitmovin.APIBaseURL + relativeURL
	_, err := url.Parse(fullURL)
	if err != nil {
		return nil, err
	}
	return r.makeRequest("GET", fullURL, nil, 2)
}

func (r *RestService) Delete(relativeURL string) ([]byte, error) {
	fullURL := *r.Bitmovin.APIBaseURL + relativeURL
	_, err := url.Parse(fullURL)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("DELETE", fullURL, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", *r.Bitmovin.APIKey)
	if r.Bitmovin.OrganizationID != nil {
		req.Header.Set("X-Tenant-Org-Id", *r.Bitmovin.OrganizationID)
	}
	req.Header.Set("X-Api-Client", ClientName)
	req.Header.Set("X-Api-Client-Version", Version)
	if os.Getenv("DUMP_TRAFFIC") != "" {
		b, _ := httputil.DumpRequest(req, true)
		println(string(b))
	}
	resp, err := r.Bitmovin.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if os.Getenv("DUMP_TRAFFIC") != "" {
		b, _ := httputil.DumpResponse(resp, true)
		println(string(b))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 399 {
		return nil, formatError(body)
	}

	return body, nil
}

//TODO default value version
func (r *RestService) List(relativeURL string, offset int64, limit int64) ([]byte, error) {
	queryParams := fmt.Sprintf("?offset=%v&limit=%v", offset, limit)
	fullURL := *r.Bitmovin.APIBaseURL + relativeURL + queryParams

	req, err := http.NewRequest("GET", fullURL, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", *r.Bitmovin.APIKey)
	if r.Bitmovin.OrganizationID != nil {
		req.Header.Set("X-Tenant-Org-Id", *r.Bitmovin.OrganizationID)
	}
	req.Header.Set("X-Api-Client", ClientName)
	req.Header.Set("X-Api-Client-Version", Version)
	if os.Getenv("DUMP_TRAFFIC") != "" {
		b, _ := httputil.DumpRequest(req, true)
		println(string(b))
	}
	resp, err := r.Bitmovin.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if os.Getenv("DUMP_TRAFFIC") != "" {
		b, _ := httputil.DumpResponse(resp, true)
		println(string(b))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (r *RestService) RetrieveCustomData(relativeURL string) ([]byte, error) {
	fullURL := *r.Bitmovin.APIBaseURL + relativeURL + "/customData"
	_, err := url.Parse(fullURL)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", fullURL, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", *r.Bitmovin.APIKey)
	if r.Bitmovin.OrganizationID != nil {
		req.Header.Set("X-Tenant-Org-Id", *r.Bitmovin.OrganizationID)
	}
	req.Header.Set("X-Api-Client", ClientName)
	req.Header.Set("X-Api-Client-Version", Version)
	if os.Getenv("DUMP_TRAFFIC") != "" {
		b, _ := httputil.DumpRequest(req, true)
		println(string(b))
	}
	resp, err := r.Bitmovin.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if os.Getenv("DUMP_TRAFFIC") != "" {
		b, _ := httputil.DumpResponse(resp, true)
		println(string(b))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (r *RestService) Update(relativeURL string, input []byte) ([]byte, error) {
	fullURL := *r.Bitmovin.APIBaseURL + relativeURL
	_, err := url.Parse(fullURL)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", fullURL, bytes.NewBuffer(input))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", *r.Bitmovin.APIKey)
	if r.Bitmovin.OrganizationID != nil {
		req.Header.Set("X-Tenant-Org-Id", *r.Bitmovin.OrganizationID)
	}
	req.Header.Set("X-Api-Client", ClientName)
	req.Header.Set("X-Api-Client-Version", Version)
	if os.Getenv("DUMP_TRAFFIC") != "" {
		b, _ := httputil.DumpRequest(req, true)
		println(string(b))
	}
	resp, err := r.Bitmovin.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if os.Getenv("DUMP_TRAFFIC") != "" {
		b, _ := httputil.DumpResponse(resp, true)
		println(string(b))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 399 {
		return nil, formatError(body)
	}

	return body, nil
}

func unmarshalError(body []byte) (*models.DataEnvelope, error) {
	var d models.DataEnvelope
	err := json.Unmarshal(body, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func formatError(body []byte) error {
	dataEnvelope, err := unmarshalError(body)
	if err != nil {
		return err
	}
	be := models.BitmovinError{
		DataEnvelope: *dataEnvelope,
	}
	return be
}
