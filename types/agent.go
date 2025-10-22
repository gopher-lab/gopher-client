package types

import (
	"github.com/masa-finance/tee-worker/api/types"
)

type Sentiment string
type TopInfluencers []string
type Topic string

const (
	SentimentBullish Sentiment = "bullish"
	SentimentBearish Sentiment = "bearish"
	SentimentNeutral Sentiment = "neutral"
)

// TopicSummary represents a topic with sentiment and influencers discovered by the agent
type TopicSummary struct {
	Topic          Topic          `json:"topic"`
	Sentiment      Sentiment      `json:"sentiment"`
	TopInfluencers TopInfluencers `json:"top_influencers"`
}

// Output is the agent's structured output
type Output struct {
	Topics    []TopicSummary   `json:"topics"`
	Documents []types.Document `json:"documents,omitempty"`
}
