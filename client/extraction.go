package client

import (
	"encoding/json"

	"github.com/gopher-lab/gopher-client/types"
)

// ExtractSearchTerms extracts optimized search terms from user input using AI
//
// Args:
//   - userInput: The user input to extract search terms from
//   - maxTerms: Maximum number of search terms to extract (1-6, default 4). Use 0 for default.
//
// Returns:
//   - A pointer to ExtractionResponse containing the extracted search terms and metadata, or an error if the operation fails
func (c *Client) ExtractSearchTerms(userInput string, maxTerms int) (*types.ExtractionResponse, error) {
	// Set default maxTerms if not provided or invalid
	if maxTerms <= 0 || maxTerms > 6 {
		maxTerms = 4
	}

	request := types.ExtractionRequest{
		UserInput: userInput,
		MaxTerms:  maxTerms,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	var response types.ExtractionResponse
	err = c.doImmediateRequest(c.BaseURL+"/v1/extraction", requestBody, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
