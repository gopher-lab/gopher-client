package client

import (
	"encoding/json"

	"github.com/gopher-lab/gopher-client/types"
)

// ContextualizeQuery enhances search queries using chat history context via OpenRouter
//
// Args:
//   - currentQuery: Current search query to contextualize
//   - chatHistory: Previous chat history for context
//   - maxHistoryItems: Maximum number of history items to use (1-10, default 5). Use 0 for default.
//
// Returns:
//   - A pointer to ContextualizeResponse containing the contextualized query and metadata, or an error if the operation fails
func (c *Client) ContextualizeQuery(currentQuery string, chatHistory []types.ChatHistoryItem, maxHistoryItems int) (*types.ContextualizeResponse, error) {
	// Set default maxHistoryItems if not provided or invalid
	if maxHistoryItems <= 0 || maxHistoryItems > 10 {
		maxHistoryItems = 5
	}

	request := types.ContextualizeRequest{
		CurrentQuery:    currentQuery,
		ChatHistory:     chatHistory,
		MaxHistoryItems: maxHistoryItems,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	var response types.ContextualizeResponse
	err = c.doImmediateRequest(c.BaseURL+"/v1/contextualize", requestBody, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

