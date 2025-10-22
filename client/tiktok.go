package client

import (
	"encoding/json"
	"fmt"

	"github.com/masa-finance/tee-worker/api/args/tiktok"
	"github.com/masa-finance/tee-worker/api/params"
	"github.com/masa-finance/tee-worker/api/types"
)

// TranscribeTikTok performs a TikTok transcription and waits for completion, returning results directly
func (c *Client) TranscribeTikTok(url string) ([]types.Document, error) {
	resp, err := c.TranscribeTikTokAsync(url)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// TranscribeTikTokAsync performs a TikTok transcription job and returns a job ID
func (c *Client) TranscribeTikTokAsync(url string) (*types.ResultResponse, error) {
	args := tiktok.NewTranscriptionArguments()
	args.VideoURL = url

	body, err := json.Marshal(params.TikTokTranscription{
		JobType: types.TiktokJob,
		Args:    &args,
	})
	if err != nil {
		return nil, err
	}
	return c.doRequest(c.BaseURL+jobEndpoint, body)
}

// TranscribeTikTokWithArgs transcribes TikTok with custom arguments and waits for completion, returning results directly
func (c *Client) TranscribeTikTokWithArgs(args tiktok.TranscriptionArguments) ([]types.Document, error) {
	body, err := json.Marshal(params.TikTokTranscription{
		JobType: types.TiktokJob,
		Args:    &args,
	})
	if err != nil {
		return nil, err
	}
	resp, err := c.doRequest(c.BaseURL+jobEndpoint, body)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// TranscribeTikTokWithArgsAsync transcribes TikTok with custom arguments and returns a job ID
func (c *Client) TranscribeTikTokWithArgsAsync(args tiktok.TranscriptionArguments) (*types.ResultResponse, error) {
	body, err := json.Marshal(params.TikTokTranscription{
		JobType: types.TiktokJob,
		Args:    &args,
	})
	if err != nil {
		return nil, err
	}
	return c.doRequest(c.BaseURL+jobEndpoint, body)
}

// SearchTikTok performs a TikTok search and waits for completion, returning results directly
func (c *Client) SearchTikTok(query string) ([]types.Document, error) {
	resp, err := c.SearchTikTokAsync(query)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// SearchTikTokAsync performs a TikTok search job and returns a job ID
func (c *Client) SearchTikTokAsync(query string) (*types.ResultResponse, error) {
	args := tiktok.NewQueryArguments()
	args.Search = []string{query}

	body, err := json.Marshal(params.TikTokSearch{
		JobType: types.TiktokJob,
		Args:    &args,
	})
	if err != nil {
		return nil, err
	}
	return c.doRequest(c.BaseURL+jobEndpoint, body)
}

// SearchTikTokWithArgs searches TikTok with query arguments and waits for completion, returning results directly
func (c *Client) SearchTikTokWithArgs(args tiktok.QueryArguments) ([]types.Document, error) {
	body, err := json.Marshal(params.TikTokSearch{
		JobType: types.TiktokJob,
		Args:    &args,
	})
	if err != nil {
		return nil, err
	}
	resp, err := c.doRequest(c.BaseURL+jobEndpoint, body)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// SearchTikTokWithArgsAsync searches TikTok with query arguments and returns a job ID
func (c *Client) SearchTikTokWithArgsAsync(args tiktok.QueryArguments) (*types.ResultResponse, error) {
	body, err := json.Marshal(params.TikTokSearch{
		JobType: types.TiktokJob,
		Args:    &args,
	})
	if err != nil {
		return nil, err
	}
	return c.doRequest(c.BaseURL+jobEndpoint, body)
}

// SearchTikTokTrending performs a TikTok trending search and waits for completion, returning results directly
func (c *Client) SearchTikTokTrending(sortBy string) ([]types.Document, error) {
	resp, err := c.SearchTikTokTrendingAsync(sortBy)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// SearchTikTokTrendingAsync performs a TikTok trending search job and returns a job ID
func (c *Client) SearchTikTokTrendingAsync(sortBy string) (*types.ResultResponse, error) {
	args := tiktok.NewTrendingArguments()
	args.SortBy = sortBy

	body, err := json.Marshal(params.TikTokTrending{
		JobType: types.TiktokJob,
		Args:    &args,
	})
	if err != nil {
		return nil, err
	}
	return c.doRequest(c.BaseURL+jobEndpoint, body)
}

// SearchTikTokTrendingWithArgs searches TikTok trending with custom arguments and waits for completion, returning results directly
func (c *Client) SearchTikTokTrendingWithArgs(args tiktok.TrendingArguments) ([]types.Document, error) {
	body, err := json.Marshal(params.TikTokTrending{
		JobType: types.TiktokJob,
		Args:    &args,
	})
	if err != nil {
		return nil, err
	}
	resp, err := c.doRequest(c.BaseURL+jobEndpoint, body)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// SearchTikTokTrendingWithArgsAsync searches TikTok trending with custom arguments and returns a job ID
func (c *Client) SearchTikTokTrendingWithArgsAsync(args tiktok.TrendingArguments) (*types.ResultResponse, error) {
	body, err := json.Marshal(params.TikTokTrending{
		JobType: types.TiktokJob,
		Args:    &args,
	})
	if err != nil {
		return nil, err
	}
	return c.doRequest(c.BaseURL+jobEndpoint, body)
}
