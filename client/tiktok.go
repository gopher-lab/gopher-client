package client

import (
	"encoding/json"
	"fmt"

	"github.com/masa-finance/tee-worker/api/args/tiktok"
	"github.com/masa-finance/tee-worker/api/params"
	"github.com/masa-finance/tee-worker/api/types"
)

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

// SearchTikTokAsync performs a TikTok search job and returns a job ID
func (c *Client) SearchTikTokAsync(q string, maxItems uint) (*types.ResultResponse, error) {
	args := tiktok.NewQueryArguments()
	args.Search = []string{q}
	args.MaxItems = maxItems

	body, err := json.Marshal(params.TikTokSearch{
		JobType: types.TiktokJob,
		Args:    &args,
	})
	if err != nil {
		return nil, err
	}
	return c.doRequest(c.BaseURL+jobEndpoint, body)
}

// SearchTikTokTrendingAsync performs a TikTok trending search job and returns a job ID
func (c *Client) SearchTikTokTrendingAsync(sortBy string, maxItems int) (*types.ResultResponse, error) {
	args := tiktok.NewTrendingArguments()
	args.SortBy = sortBy
	args.MaxItems = maxItems

	body, err := json.Marshal(params.TikTokTrending{
		JobType: types.TiktokJob,
		Args:    &args,
	})
	if err != nil {
		return nil, err
	}
	return c.doRequest(c.BaseURL+jobEndpoint, body)
}

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

// SearchTikTok performs a TikTok search and waits for completion, returning results directly
func (c *Client) SearchTikTok(q string, maxItems uint) ([]types.Document, error) {
	resp, err := c.SearchTikTokAsync(q, maxItems)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// SearchTikTokTrending performs a TikTok trending search and waits for completion, returning results directly
func (c *Client) SearchTikTokTrending(sortBy string, maxItems int) ([]types.Document, error) {
	resp, err := c.SearchTikTokTrendingAsync(sortBy, maxItems)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
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
