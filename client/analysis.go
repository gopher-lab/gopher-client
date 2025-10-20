package client

import (
	"encoding/json"

	"github.com/gopher-lab/gopher-client/types"
)

// AnalyzeData analyzes tweets and other data using various AI models via OpenRouter
//
// Args:
//   - tweets: Array of tweets to analyze
//   - prompt: Analysis prompt
//   - model: AI model to use for analysis (optional, defaults to "openai/gpt-4o-mini")
//   - app: Whether this is an app request (optional, defaults to false)
//   - chatHistory: Previous chat history for context (optional)
//   - currentQuery: Current query being analyzed (optional)
//
// Returns:
//   - A pointer to AnalysisResponse containing the analysis results and metadata, or an error if the operation fails
func (c *Client) AnalyzeData(data []string, prompt string, model string, app bool, chatHistory []types.ChatHistoryItem, currentQuery string) (*types.AnalysisResponse, error) {
	// Set default model if not provided
	if model == "" {
		model = "openai/gpt-4o-mini"
	}

	request := types.AnalysisRequest{
		Tweets:       data,
		Prompt:       prompt,
		Model:        model,
		App:          app,
		ChatHistory:  chatHistory,
		CurrentQuery: currentQuery,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	var response types.AnalysisResponse
	err = c.doImmediateRequest(c.BaseURL+"/v1/analysis", requestBody, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// AnalyzeDataSimple analyzes tweets with a simple prompt using default settings
//
// Args:
//   - tweets: Array of tweets to analyze
//   - prompt: Analysis prompt
//
// Returns:
//   - A pointer to AnalysisResponse containing the analysis results and metadata, or an error if the operation fails
func (c *Client) AnalyzeDataSimple(tweets []string, prompt string) (*types.AnalysisResponse, error) {
	return c.AnalyzeData(tweets, prompt, "", false, nil, "")
}

// GetAvailableModels retrieves the list of available AI models for analysis
//
// Returns:
//   - A slice of strings containing available model names, or an error if the operation fails
func (c *Client) GetAvailableModels() ([]string, error) {
	var models []string
	err := c.doResultRequest(c.BaseURL+"/v1/analysis", &models)
	if err != nil {
		return nil, err
	}
	return models, nil
}
