package client

import (
	"encoding/json"

	"github.com/gopher-lab/gopher-client/log"
	"github.com/masa-finance/tee-worker/api/params"
	"github.com/masa-finance/tee-worker/api/types"
)

// SearchHybrid performs a hybrid search and returns results directly
func (c *Client) SearchHybrid(
	query string,
	sources []types.Source,
	text string,
	queryWeight float64,
	textWeight float64,
	keywords []string,
	operator string,
	maxResults int,
) ([]types.Document, error) {
	requestBody, err := json.Marshal(params.HybridSearch{
		TextQuery:       params.HybridQuery{Query: query, Weight: queryWeight},
		SimilarityQuery: params.HybridQuery{Query: text, Weight: textWeight},
		Sources:         sources,
		Keywords:        keywords,
		Operator:        operator,
		MaxResults:      maxResults,
	})
	if err != nil {
		log.Error("Error while performing hybrid web search", "query", query, "text", text, "error", err.Error())
		return nil, err
	}

	var results []types.Document
	err = c.doImmediateRequest(c.BaseURL+"/v1/search/hybrid", requestBody, &results)
	if err != nil {
		log.Error("Error while performing hybrid web search", "query", query, "text", text, "error", err.Error())
		return nil, err
	}
	return results, nil
}
