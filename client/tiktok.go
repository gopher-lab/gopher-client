package client

import (
	"encoding/json"

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
