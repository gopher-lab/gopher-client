package client

import (
	"encoding/json"

	"github.com/masa-finance/tee-worker/api/args/twitter/search"
	"github.com/masa-finance/tee-worker/api/jobs"
	"github.com/masa-finance/tee-worker/api/types"
)

// PostTwitterJob posts a Twitter job to the API
func (c *Client) PostTwitterJob(args search.Arguments) (*types.ResultResponse, error) {
	body, err := json.Marshal(jobs.TwitterParams{
		JobType: types.TwitterJob,
		Args:    args,
	})
	if err != nil {
		return nil, err
	}
	return c.doRequest(c.BaseURL+jobEndpoint, body)
}

// PerformTwitterSearch performs a Twitter search job
func (c *Client) PerformTwitterSearch(query string) (*types.ResultResponse, error) {
	args := search.NewArguments()
	args.Query = query
	res, err := c.PostTwitterJob(args)
	if err != nil {
		return nil, err
	}
	return res, nil
}
