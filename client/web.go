package client

import (
	"encoding/json"

	"github.com/masa-finance/tee-worker/api/args/web/page"
	"github.com/masa-finance/tee-worker/api/jobs"
	"github.com/masa-finance/tee-worker/api/types"
)

// PostWebJob posts a Web job to the API
func (c *Client) PostWebJob(args page.Arguments) (*types.ResultResponse, error) {
	body, err := json.Marshal(jobs.WebParams{
		JobType: types.WebJob,
		Args:    args,
	})
	if err != nil {
		return nil, err
	}
	return c.doRequest(c.BaseURL+jobEndpoint, body)
}

// PerformWebSearch performs a web search job using the provided URL
func (c *Client) PerformWebSearch(url string) (*types.ResultResponse, error) {
	args := page.NewArguments()
	args.URL = url
	res, err := c.PostWebJob(args)
	if err != nil {
		return nil, err
	}
	return res, nil
}
