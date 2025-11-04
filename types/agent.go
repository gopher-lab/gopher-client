package types

import (
	"github.com/masa-finance/tee-worker/v2/api/types"
)

type Sentiment uint
type Reasoning string
type Asset string

// TopicSummary represents a topic with sentiment and influencers discovered by the agent
type AssetSummary struct {
	Asset     Asset     `json:"asset"`
	Reasoning Reasoning `json:"reasoning"`
	Sentiment Sentiment `json:"sentiment"`
}

// Output is the agent's structured output
type Output struct {
	Assets    []AssetSummary   `json:"assets"`
	Documents []types.Document `json:"documents,omitempty"`
}
