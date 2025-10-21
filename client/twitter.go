package client

import (
	"encoding/json"
	"fmt"

	"github.com/masa-finance/tee-worker/api/args/twitter/search"
	"github.com/masa-finance/tee-worker/api/jobs"
	"github.com/masa-finance/tee-worker/api/types"
)

// SearchTwitterWithArgsAsync searches Twitter with custom arguments and returns a job ID
func (c *Client) SearchTwitterWithArgsAsync(args search.Arguments) (*types.ResultResponse, error) {
	body, err := json.Marshal(jobs.TwitterParams{
		JobType: types.TwitterJob,
		Args:    args,
	})
	if err != nil {
		return nil, err
	}
	return c.doRequest(c.BaseURL+jobEndpoint, body)
}

// SearchTwitterAsync performs a Twitter search job and returns a job ID
func (c *Client) SearchTwitterAsync(query string) (*types.ResultResponse, error) {
	args := search.NewArguments()
	args.Query = query
	res, err := c.SearchTwitterWithArgsAsync(args)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SearchTwitter performs a Twitter search and waits for completion, returning results directly
func (c *Client) SearchTwitter(query string) ([]types.Document, error) {
	resp, err := c.SearchTwitterAsync(query)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// SearchTwitterWithArgs searches Twitter with custom arguments and waits for completion, returning results directly
func (c *Client) SearchTwitterWithArgs(args search.Arguments) ([]types.Document, error) {
	resp, err := c.SearchTwitterWithArgsAsync(args)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}
