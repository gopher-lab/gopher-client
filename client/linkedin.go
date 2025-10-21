package client

import (
	"encoding/json"
	"fmt"

	"github.com/masa-finance/tee-worker/api/args/linkedin/profile"
	"github.com/masa-finance/tee-worker/api/jobs"
	"github.com/masa-finance/tee-worker/api/types"
	ptypes "github.com/masa-finance/tee-worker/api/types/linkedin/profile"
)

// SearchLinkedInWithArgsAsync searches LinkedIn with custom arguments and returns a job ID
func (c *Client) SearchLinkedInWithArgsAsync(args profile.Arguments) (*types.ResultResponse, error) {
	body, err := json.Marshal(jobs.LinkedInParams{
		JobType: types.LinkedInJob,
		Args:    args,
	})
	if err != nil {
		return nil, err
	}
	return c.doRequest(c.BaseURL+jobEndpoint, body)
}

// SearchLinkedInAsync performs a LinkedIn search job and returns a job ID
func (c *Client) SearchLinkedInAsync(query string, mode ptypes.ScraperMode) (*types.ResultResponse, error) {
	args := profile.NewArguments()
	args.ScraperMode = mode
	args.Query = query

	res, err := c.SearchLinkedInWithArgsAsync(args)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SearchLinkedIn performs a LinkedIn search and waits for completion, returning results directly
func (c *Client) SearchLinkedIn(query string, mode ptypes.ScraperMode) ([]types.Document, error) {
	resp, err := c.SearchLinkedInAsync(query, mode)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// SearchLinkedInWithArgs searches LinkedIn with custom arguments and waits for completion, returning results directly
func (c *Client) SearchLinkedInWithArgs(args profile.Arguments) ([]types.Document, error) {
	resp, err := c.SearchLinkedInWithArgsAsync(args)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}
