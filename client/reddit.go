package client

import (
	"encoding/json"
	"fmt"

    "github.com/masa-finance/tee-worker/v2/api/args/reddit"
    "github.com/masa-finance/tee-worker/v2/api/params"
    "github.com/masa-finance/tee-worker/v2/api/types"
)

// ScrapeRedditURLAsync performs a Reddit URL scraping job and returns a job ID
func (c *Client) ScrapeRedditURLAsync(url string) (*types.ResultResponse, error) {
	args := reddit.NewScrapeUrlsArguments()
	args.URLs = []string{url}
	res, err := c.SearchRedditWithArgsAsync(args)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SearchRedditPostsAsync performs a Reddit posts search job and returns a job ID
func (c *Client) SearchRedditPostsAsync(query string) (*types.ResultResponse, error) {
	args := reddit.NewSearchPostsArguments()
	args.Queries = []string{query}
	res, err := c.SearchRedditWithArgsAsync(args)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SearchRedditUsersAsync performs a Reddit users search job and returns a job ID
func (c *Client) SearchRedditUsersAsync(query string) (*types.ResultResponse, error) {
	args := reddit.NewSearchUsersArguments()
	args.Queries = []string{query}
	res, err := c.SearchRedditWithArgsAsync(args)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SearchRedditCommunitiesAsync performs a Reddit communities search job and returns a job ID
func (c *Client) SearchRedditCommunitiesAsync(query string) (*types.ResultResponse, error) {
	args := reddit.NewSearchCommunitiesArguments()
	args.Queries = []string{query}
	res, err := c.SearchRedditWithArgsAsync(args)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ScrapeRedditURL performs a Reddit URL scraping and waits for completion, returning results directly
func (c *Client) ScrapeRedditURL(url string) ([]types.Document, error) {
	resp, err := c.ScrapeRedditURLAsync(url)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// SearchRedditPosts performs a Reddit posts search and waits for completion, returning results directly
func (c *Client) SearchRedditPosts(query string) ([]types.Document, error) {
	resp, err := c.SearchRedditPostsAsync(query)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// SearchRedditUsers performs a Reddit users search and waits for completion, returning results directly
func (c *Client) SearchRedditUsers(query string) ([]types.Document, error) {
	resp, err := c.SearchRedditUsersAsync(query)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// SearchRedditCommunities performs a Reddit communities search and waits for completion, returning results directly
func (c *Client) SearchRedditCommunities(query string) ([]types.Document, error) {
	resp, err := c.SearchRedditCommunitiesAsync(query)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// SearchRedditWithArgs searches Reddit with custom arguments and waits for completion, returning results directly
func (c *Client) SearchRedditWithArgs(args reddit.SearchArguments) ([]types.Document, error) {
	resp, err := c.SearchRedditWithArgsAsync(args)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("job submission failed: %s", resp.Error)
	}
	return c.WaitForJobCompletion(resp.UUID)
}

// SearchRedditWithArgsAsync searches Reddit with custom arguments and returns a job ID
func (c *Client) SearchRedditWithArgsAsync(args reddit.SearchArguments) (*types.ResultResponse, error) {
	params := params.Params[*reddit.SearchArguments]{}
	params.JobType = types.RedditJob
	params.Args = &args

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	return c.doRequest(c.BaseURL+jobEndpoint, body)
}
