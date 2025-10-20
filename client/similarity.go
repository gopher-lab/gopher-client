package client

import (
	"encoding/json"

	"github.com/masa-finance/gopher-client/log"
	"github.com/masa-finance/tee-worker/api/jobs"
	"github.com/masa-finance/tee-worker/api/types"
)

// PerformSimilaritySearch sends a POST request to the web similarity search endpoint
func (c *Client) PerformSimilaritySearch(
	query string,
	sources []types.Source,
	keywords []string,
	operator string,
	maxResults int,
	receiver any,
) error {
	requestBody, err := json.Marshal(jobs.SimilaritySearchParams{
		Query:           query,
		Keywords:        keywords,
		Sources:         sources,
		MaxResults:      maxResults,
		KeywordOperator: operator,
	})
	if err != nil {
		log.Error("Error while performing similarity search", "query", query, "keywords", keywords, "error", err.Error())
		return err
	}

	err = c.doImmediateRequest(c.BaseURL+"/v1/search/similarity", requestBody, receiver)
	if err != nil {
		log.Error("Error while performing similarity search", "query", query, "keywords", keywords, "error", err.Error())
	}
	return err
}
