package client

import (
	"encoding/json"
	"fmt"

	"github.com/masa-finance/tee-worker/api/args/tiktok/query"
	"github.com/masa-finance/tee-worker/api/args/tiktok/transcription"
	"github.com/masa-finance/tee-worker/api/args/tiktok/trending"
	"github.com/masa-finance/tee-worker/api/jobs"
	"github.com/masa-finance/tee-worker/api/types"
)

// PostTikTokJob posts a TikTok job to the API
func (c *Client) PostTikTokJob(args map[string]any) (*types.ResultResponse, error) {
	body, err := json.Marshal(jobs.TikTokParams{
		JobType: types.TiktokJob,
		Args:    args,
	})
	if err != nil {
		return nil, err
	}
	return c.doRequest(c.BaseURL+jobEndpoint, body)
}

func (c *Client) PerformTikTokTranscription(url string) (*types.ResultResponse, error) {
	args := transcription.NewArguments()
	args.VideoURL = url

	body, err := json.Marshal(jobs.TikTokTranscriptionParams{
		JobType: types.TiktokJob,
		Args:    args,
	})
	if err != nil {
		return nil, err
	}
	return c.doRequest(c.BaseURL+jobEndpoint, body)
}

func (c *Client) PerformTikTokSearch(q string, maxItems uint) (*types.ResultResponse, error) {
	args := query.NewArguments()
	args.Search = []string{q}
	args.MaxItems = maxItems

	body, err := json.Marshal(jobs.TikTokSearchParams{
		JobType: types.TiktokJob,
		Args:    args,
	})
	if err != nil {
		return nil, err
	}
	return c.doRequest(c.BaseURL+jobEndpoint, body)
}

func (c *Client) PerformTikTokSearchByTrending(sortBy string, maxItems int) (*types.ResultResponse, error) {
	args := trending.NewArguments()
	args.SortBy = sortBy
	args.MaxItems = maxItems

	body, err := json.Marshal(jobs.TikTokTrendingParams{
		JobType: types.TiktokJob,
		Args:    args,
	})
	if err != nil {
		return nil, err
	}
	return c.doRequest(c.BaseURL+jobEndpoint, body)
}

// PerformTikTokTranscriptionAndWait performs a TikTok transcription and waits for completion
func (c *Client) PerformTikTokTranscriptionAndWait(url string) ([]types.Document, error) {
	resp, err := c.PerformTikTokTranscription(url)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// PerformTikTokSearchAndWait performs a TikTok search and waits for completion
func (c *Client) PerformTikTokSearchAndWait(q string, maxItems uint) ([]types.Document, error) {
	resp, err := c.PerformTikTokSearch(q, maxItems)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// PerformTikTokSearchByTrendingAndWait performs a TikTok trending search and waits for completion
func (c *Client) PerformTikTokSearchByTrendingAndWait(sortBy string, maxItems int) ([]types.Document, error) {
	resp, err := c.PerformTikTokSearchByTrending(sortBy, maxItems)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// PostTikTokJobAndWait posts a TikTok job and waits for completion
func (c *Client) PostTikTokJobAndWait(args map[string]any) ([]types.Document, error) {
	resp, err := c.PostTikTokJob(args)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// PostTikTokSearchJobAndWait posts a TikTok search job with flexible arguments and waits for completion
func (c *Client) PostTikTokSearchJobAndWait(args query.Arguments) ([]types.Document, error) {
	body, err := json.Marshal(jobs.TikTokSearchParams{
		JobType: types.TiktokJob,
		Args:    args,
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

// PostTikTokTrendingJobAndWait posts a TikTok trending job with flexible arguments and waits for completion
func (c *Client) PostTikTokTrendingJobAndWait(args trending.Arguments) ([]types.Document, error) {
	body, err := json.Marshal(jobs.TikTokTrendingParams{
		JobType: types.TiktokJob,
		Args:    args,
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

// PostTikTokTranscriptionJobAndWait posts a TikTok transcription job with flexible arguments and waits for completion
func (c *Client) PostTikTokTranscriptionJobAndWait(args transcription.Arguments) ([]types.Document, error) {
	body, err := json.Marshal(jobs.TikTokTranscriptionParams{
		JobType: types.TiktokJob,
		Args:    args,
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
