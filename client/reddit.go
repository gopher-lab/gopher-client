package client

import (
	"encoding/json"
	"fmt"

	"github.com/masa-finance/tee-worker/api/args/reddit/search"
	"github.com/masa-finance/tee-worker/api/jobs"
	"github.com/masa-finance/tee-worker/api/types"
)

// SearchRedditWithArgsAsync searches Reddit with custom arguments and returns a job ID
func (c *Client) SearchRedditWithArgsAsync(args search.Arguments) (*types.ResultResponse, error) {
	body, err := json.Marshal(jobs.RedditParams{
		JobType: types.RedditJob,
		Args:    args,
	})
	if err != nil {
		return nil, err
	}
	return c.doRequest(c.BaseURL+jobEndpoint, body)
}

// ScrapeRedditURLAsync performs a Reddit URL scraping job and returns a job ID
func (c *Client) ScrapeRedditURLAsync(url string, maxItems uint) (*types.ResultResponse, error) {
	args := search.NewScrapeUrlsArguments()
	args.URLs = []string{url}
	args.MaxItems = maxItems
	res, err := c.SearchRedditWithArgsAsync(args)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SearchRedditPostsAsync performs a Reddit posts search job and returns a job ID
func (c *Client) SearchRedditPostsAsync(query string, maxItems uint) (*types.ResultResponse, error) {
	args := search.NewSearchPostsArguments()
	args.Queries = []string{query}
	args.MaxItems = maxItems
	res, err := c.SearchRedditWithArgsAsync(args)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SearchRedditUsersAsync performs a Reddit users search job and returns a job ID
func (c *Client) SearchRedditUsersAsync(query string, maxItems uint) (*types.ResultResponse, error) {
	args := search.NewSearchUsersArguments()
	args.Queries = []string{query}
	args.MaxItems = maxItems
	res, err := c.SearchRedditWithArgsAsync(args)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SearchRedditCommunitiesAsync performs a Reddit communities search job and returns a job ID
func (c *Client) SearchRedditCommunitiesAsync(query string, maxItems uint) (*types.ResultResponse, error) {
	args := search.NewSearchCommunitiesArguments()
	args.Queries = []string{query}
	args.MaxItems = maxItems
	res, err := c.SearchRedditWithArgsAsync(args)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ScrapeRedditURL performs a Reddit URL scraping and waits for completion, returning results directly
func (c *Client) ScrapeRedditURL(url string, maxItems uint) ([]types.Document, error) {
	resp, err := c.ScrapeRedditURLAsync(url, maxItems)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// SearchRedditPosts performs a Reddit posts search and waits for completion, returning results directly
func (c *Client) SearchRedditPosts(query string, maxItems uint) ([]types.Document, error) {
	resp, err := c.SearchRedditPostsAsync(query, maxItems)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// SearchRedditUsers performs a Reddit users search and waits for completion, returning results directly
func (c *Client) SearchRedditUsers(query string, maxItems uint) ([]types.Document, error) {
	resp, err := c.SearchRedditUsersAsync(query, maxItems)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// SearchRedditCommunities performs a Reddit communities search and waits for completion, returning results directly
func (c *Client) SearchRedditCommunities(query string, maxItems uint) ([]types.Document, error) {
	resp, err := c.SearchRedditCommunitiesAsync(query, maxItems)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// SearchRedditWithArgs searches Reddit with custom arguments and waits for completion, returning results directly
func (c *Client) SearchRedditWithArgs(args search.Arguments) ([]types.Document, error) {
	resp, err := c.SearchRedditWithArgsAsync(args)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}
