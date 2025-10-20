package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/masa-finance/tee-worker/api/types"
)

// note, single endpoint for all jobs, supported by indexer and data app for acceptance tests
const jobEndpoint = "/v1/search/live"

// Client represents the API client
type Client struct {
	BaseURL string
	Token   string
}

// NewClient creates a new API client
func NewClient(baseURL string, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
	}
}

func getErrorFromResponse(body []byte) error {
	result := struct {
		Error string `json:"error"`
	}{}

	_ = json.Unmarshal(body, &result)

	if result.Error != "" {
		return fmt.Errorf("job errored: %s", result.Error)
	}

	return nil
}

func (c *Client) doRequest(url string, requestBody []byte) (*types.ResultResponse, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request to %s: %w", url, err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to do POST request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body from POST request to %s: %w", url, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("job errored: Status code %d during call to %s with request body %s. Response body: %s",
			resp.StatusCode, url, requestBody, body)
	}

	var searchResponse types.ResultResponse
	if err := json.Unmarshal(body, &searchResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal GET %s response %s: %w", url, body, err)
	}

	return &searchResponse, getErrorFromResponse(body)
}

func (c *Client) doStatusRequest(url string) (*types.IndexerJobResult, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request to %s: %w", url, err)
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do GET request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body from POST request to %s: %w", url, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("job errored: Status code %d during call to %s. Response body: %s", resp.StatusCode, url, body)
	}

	var jobStatusResponse types.IndexerJobResult
	if err := json.Unmarshal(body, &jobStatusResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal GET %s response %s: %w", url, body, err)
	}

	return &jobStatusResponse, getErrorFromResponse(body)
}

func (c *Client) doResultRequest(url string, receiver any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create GET request to %s: %w", url, err)
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do GET request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body from POST request to %s: %w", url, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("job errored: Status code %d during call to %s. Response body: %s", resp.StatusCode, url, body)
	}

	if err := json.Unmarshal(body, receiver); err != nil {
		return fmt.Errorf("failed to unmarshal GET %s response %s: %w", url, body, err)
	}

	return getErrorFromResponse(body)
}

func (c *Client) doImmediateRequest(url string, requestBody []byte, receiver any) error {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create POST request to %s: %w", url, err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do POST request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body from POST request to %s: %w", url, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("job errored: Status code %d during call to %s. Body: %s", resp.StatusCode, url, body)
	}

	if err := json.Unmarshal(body, receiver); err != nil {
		return fmt.Errorf("Error during unmarshal: %#w. URL: %s. Request: '%s'. Response: '%s'", err, url, requestBody, body)
	}

	return getErrorFromResponse(body)
}

// GetJobStatus sends a GET request to the job status endpoint
func (c *Client) GetJobStatus(jobID string) (*types.IndexerJobResult, error) {
	url := c.BaseURL + jobEndpoint + "/status/" + jobID
	return c.doStatusRequest(url)
}

// GetResult sends a GET request to the job result endpoint
func (c *Client) GetResult(jobID string, receiver any) error {
	url := c.BaseURL + jobEndpoint + "/result/" + jobID
	return c.doResultRequest(url, receiver)
}
