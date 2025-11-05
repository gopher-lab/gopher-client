package types

// ContextualizeRequest represents the request payload for the /v1/contextualize endpoint
type ContextualizeRequest struct {
	CurrentQuery    string            `json:"currentQuery"`              // Current search query to contextualize
	ChatHistory     []ChatHistoryItem `json:"chatHistory"`               // Previous chat history for context
	MaxHistoryItems int               `json:"maxHistoryItems,omitempty"` // Maximum number of history items to use (1-10, default 5)
}

// ContextualizeResponse represents the response from the /v1/contextualize endpoint
type ContextualizeResponse struct {
	ContextualizedQuery string `json:"contextualizedQuery"` // Enhanced query with context
	OriginalQuery       string `json:"originalQuery"`       // Original query before contextualization
	UsedContext         bool   `json:"usedContext"`         // Whether context was used to enhance the query
	Reasoning           string `json:"reasoning"`           // Explanation of how context was used
}
