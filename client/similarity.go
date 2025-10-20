package client

import (
	"encoding/json"

	"github.com/gopher-lab/gopher-client/log"
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
) ([]types.Document, error) {
	requestBody, err := json.Marshal(jobs.SimilaritySearchParams{
		Query:           query,
		Keywords:        keywords,
		Sources:         sources,
		MaxResults:      maxResults,
		KeywordOperator: operator,
	})
	if err != nil {
		log.Error("Error while performing similarity search", "query", query, "keywords", keywords, "error", err.Error())
		return nil, err
	}

	var results []types.Document
	err = c.doImmediateRequest(c.BaseURL+"/v1/search/similarity", requestBody, &results)
	if err != nil {
		log.Error("Error while performing similarity search", "query", query, "keywords", keywords, "error", err.Error())
		return nil, err
	}
	return results, nil
}
