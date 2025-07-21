package file_upload_service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// FileUploadServiceInfo struct to store meta-data related to the file upload service
type FileUploadServiceInfo struct {
	ServiceID string
	SecretKey string
}

// need to make utility function out of this
func RegisterToFileUploadService(ctx context.Context, scheme string, domain string, port string, endpoint string, headers map[string]string, payload interface{}) (*FileUploadServiceInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// need to build the url
	// right now im assuming the endpoint starts with /
	url := fmt.Sprintf("%s://%s:%s%s", scheme, domain, port, endpoint)
	// need to create a http.Request using the payload
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// set the headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	// create the transport layer
	transport := &http.Transport{
		TLSHandshakeTimeout: 5 * time.Second,
		MaxIdleConns:        1,
		IdleConnTimeout:     0 * time.Second,
	}
	client := http.Client{
		Transport: transport,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// good go-lang practise
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	var responseBody struct {
		ServiceID string `json:"service_id"`
		SecretKey string `json:"secret_key"`
	}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	return &FileUploadServiceInfo{
		ServiceID: responseBody.ServiceID,
		SecretKey: responseBody.SecretKey,
	}, nil

}
