package types

// ExtractionRequest represents the request payload for the /v1/extraction endpoint
type ExtractionRequest struct {
	UserInput string `json:"userInput"`           // User input to extract search terms from
	MaxTerms  int    `json:"maxTerms,omitempty"`  // Maximum number of search terms to extract (1-6, default 4)
}

// ExtractionResponse represents the response from the /v1/extraction endpoint
type ExtractionResponse struct {
	SearchTerm string `json:"searchTerm"` // Extracted and optimized search term
	Thinking   string `json:"thinking"`   // AI reasoning process for the extraction
	UUID       string `json:"uuid"`       // Unique identifier for this extraction request
}
