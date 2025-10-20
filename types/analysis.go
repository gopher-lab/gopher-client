package types

// ChatHistoryItem represents a single item in chat history for analysis
type ChatHistoryItem struct {
	Query     string `json:"query"`     // Previous query
	Timestamp string `json:"timestamp"` // Timestamp of the query
}

// AnalysisRequest represents the request payload for the /v1/analysis endpoint
type AnalysisRequest struct {
	Tweets       []string          `json:"tweets"`                 // Array of tweets to analyze, TODO this can be renamed to "data" once Data API is updated
	Prompt       string            `json:"prompt"`                 // Analysis prompt
	Model        string            `json:"model,omitempty"`        // AI model to use for analysis (default: "openai/gpt-4o-mini")
	App          bool              `json:"app,omitempty"`          // Whether this is an app request (default: false)
	ChatHistory  []ChatHistoryItem `json:"chatHistory,omitempty"`  // Previous chat history for context
	CurrentQuery string            `json:"currentQuery,omitempty"` // Current query being analyzed
}

// AnalysisResponse represents the response from the /v1/analysis endpoint
type AnalysisResponse struct {
	Analysis   string `json:"analysis"`    // The final analysis result
	Reasoning  string `json:"reasoning"`   // AI reasoning process
	ModelUsed  string `json:"model_used"`  // AI model used for analysis
	TokensUsed int    `json:"tokens_used"` // Number of tokens consumed
	JobUUID    string `json:"job_uuid"`    // Unique identifier for this analysis request
}
