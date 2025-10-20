package client

import (
	"encoding/json"

	"github.com/masa-finance/gopher-client/log"
	"github.com/masa-finance/tee-worker/api/jobs"
	"github.com/masa-finance/tee-worker/api/types"
)

// PerformHybridSearch sends a POST request to the web hybrid search endpoint
func (c *Client) PerformHybridSearch(
	query string,
	sources []types.Source,
	text string,
	queryWeight float64,
	textWeight float64,
	keywords []string,
	operator string,
	maxResults int,
	receiver any,
) error {
	requestBody, err := json.Marshal(jobs.HybridSearchParams{
		TextQuery:       jobs.HybridQuery{Query: query, Weight: queryWeight},
		SimilarityQuery: jobs.HybridQuery{Query: text, Weight: textWeight},
		Sources:         sources,
		Keywords:        keywords,
		Operator:        operator,
		MaxResults:      maxResults,
	})
	if err != nil {
		log.Error("Error while performing hybrid web search", "query", query, "text", text, "error", err.Error())
		return err
	}

	err = c.doImmediateRequest(c.BaseURL+"/v1/search/hybrid", requestBody, receiver)
	if err != nil {
		log.Error("Error while performing hybrid web search", "query", query, "text", text, "error", err.Error())
	}
	return err
}
