package client

import (
	"encoding/json"
	"fmt"

	"github.com/masa-finance/tee-worker/api/args/linkedin/profile"
	"github.com/masa-finance/tee-worker/api/jobs"
	"github.com/masa-finance/tee-worker/api/types"
	ptypes "github.com/masa-finance/tee-worker/api/types/linkedin/profile"
)

// PostLinkedInJob posts a LinkedIn job to the API
func (c *Client) PostLinkedInJob(args profile.Arguments) (*types.ResultResponse, error) {
	body, err := json.Marshal(jobs.LinkedInParams{
		JobType: types.LinkedInJob,
		Args:    args,
	})
	if err != nil {
		return nil, err
	}
	return c.doRequest(c.BaseURL+jobEndpoint, body)
}

// PerformLinkedInSearch performs a LinkedIn search job
func (c *Client) PerformLinkedInSearch(query string, mode ptypes.ScraperMode) (*types.ResultResponse, error) {
	args := profile.NewArguments()
	args.ScraperMode = mode
	args.Query = query

	res, err := c.PostLinkedInJob(args)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// PerformLinkedInSearchAndWait performs a LinkedIn search and waits for completion
func (c *Client) PerformLinkedInSearchAndWait(query string, mode ptypes.ScraperMode) ([]types.Document, error) {
	resp, err := c.PerformLinkedInSearch(query, mode)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// PostLinkedInJobAndWait posts a LinkedIn job and waits for completion
func (c *Client) PostLinkedInJobAndWait(args profile.Arguments) ([]types.Document, error) {
	resp, err := c.PostLinkedInJob(args)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}
