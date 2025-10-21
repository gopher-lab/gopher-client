package client

import (
	"encoding/json"
	"fmt"

	"github.com/masa-finance/tee-worker/api/args/web/page"
	"github.com/masa-finance/tee-worker/api/jobs"
	"github.com/masa-finance/tee-worker/api/types"
)

// ScrapeWebWithArgsAsync scrapes a web page with custom arguments and returns a job ID
func (c *Client) ScrapeWebWithArgsAsync(args page.Arguments) (*types.ResultResponse, error) {
	body, err := json.Marshal(jobs.WebParams{
		JobType: types.WebJob,
		Args:    args,
	})
	if err != nil {
		return nil, err
	}
	return c.doRequest(c.BaseURL+jobEndpoint, body)
}

// ScrapeWebAsync performs a web scraping job using the provided URL and returns a job ID
func (c *Client) ScrapeWebAsync(url string) (*types.ResultResponse, error) {
	args := page.NewArguments()
	args.URL = url
	res, err := c.ScrapeWebWithArgsAsync(args)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ScrapeWeb performs a web scraping job and waits for completion, returning results directly
func (c *Client) ScrapeWeb(url string) ([]types.Document, error) {
	resp, err := c.ScrapeWebAsync(url)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// ScrapeWebWithArgs scrapes a web page with custom arguments and waits for completion, returning results directly
func (c *Client) ScrapeWebWithArgs(args page.Arguments) ([]types.Document, error) {
	resp, err := c.ScrapeWebWithArgsAsync(args)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}
