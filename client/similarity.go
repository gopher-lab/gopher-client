package client

import (
	"encoding/json"

	"github.com/gopher-lab/gopher-client/log"
	"github.com/masa-finance/tee-worker/api/jobs"
	"github.com/masa-finance/tee-worker/api/types"
)

// SearchSimilarity performs a similarity search and returns results directly
func (c *Client) SearchSimilarity(
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
