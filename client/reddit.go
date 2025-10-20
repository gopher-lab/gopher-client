package client

import (
	"encoding/json"

	"github.com/masa-finance/tee-worker/api/args/reddit/search"
	"github.com/masa-finance/tee-worker/api/jobs"
	"github.com/masa-finance/tee-worker/api/types"
)

// PostRedditJob posts a Reddit job to the API
func (c *Client) PostRedditJob(args search.Arguments) (*types.ResultResponse, error) {
	body, err := json.Marshal(jobs.RedditParams{
		JobType: types.RedditJob,
		Args:    args,
	})
	if err != nil {
		return nil, err
	}
	return c.doRequest(c.BaseURL+jobEndpoint, body)
}

// PerformRedditScrapeURL performs a Reddit scrape job
func (c *Client) PerformRedditScrapeURL(url string, maxItems uint) (*types.ResultResponse, error) {
	args := search.NewScrapeUrlsArguments()
	args.URLs = []string{url}
	args.MaxItems = maxItems
	res, err := c.PostRedditJob(args)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// PerformRedditSearchPosts performs a Reddit search job
func (c *Client) PerformRedditSearchPosts(query string, maxItems uint) (*types.ResultResponse, error) {
	args := search.NewSearchPostsArguments()
	args.Queries = []string{query}
	args.MaxItems = maxItems
	res, err := c.PostRedditJob(args)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// PerformRedditSearchUsers performs a Reddit search users job
func (c *Client) PerformRedditSearchUsers(query string, maxItems uint) (*types.ResultResponse, error) {
	args := search.NewSearchUsersArguments()
	args.Queries = []string{query}
	args.MaxItems = maxItems
	res, err := c.PostRedditJob(args)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// PerformRedditSearchCommunities performs a Reddit search communities job
func (c *Client) PerformRedditSearchCommunities(query string, maxItems uint) (*types.ResultResponse, error) {
	args := search.NewSearchCommunitiesArguments()
	args.Queries = []string{query}
	args.MaxItems = maxItems
	res, err := c.PostRedditJob(args)
	if err != nil {
		return nil, err
	}
	return res, nil
}
