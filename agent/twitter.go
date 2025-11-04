package agent

import (
	"encoding/json"

	"github.com/gopher-lab/gopher-client/client"
	"github.com/masa-finance/tee-worker/v2/api/args/twitter"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

// TwitterQueryInstructions contains long-form guidance and examples for constructing Twitter queries.
// If no date range is specified by the user, searches should consider the last 1 day.
const TwitterQueryInstructions = `
watching now 	containing both "watching" and "now". This is the default operator.
"happy hour" 	containing the exact phrase "happy hour".
love OR hate 	containing either "love" or "hate" (or both).
beer -root 	containing "beer" but not "root".
#haiku 	containing the hashtag "haiku".
from:interior 	sent from Twitter account "interior". CRITICAL: NO SPACE after 'from:' (e.g., 'from:username', NOT 'from: username').
to:NASA 	a Tweet authored in reply to Twitter account "NASA".
@NASA 	mentioning Twitter account "NASA".
superhero since:2015-12-21 	containing "superhero" and sent since date "2015-12-21" (year-month-day).

If no date range is specified, default to the last 1 day (use since:YYYY-MM-DD format, typically one day before today).

CORRECT EXAMPLES:
- from:JamesWynnReal (BTC OR Bitcoin OR ETH OR Ethereum) since:2025-11-03
- from:CryptoWendyO #BTC OR #ETH since:2025-11-03
- from:VitalikButerin (ethereum OR ETH) since:2025-11-03

INCORRECT (DO NOT USE SPACES AFTER from:):
- from: JamesWynnReal (WRONG - has space after from:)
- from: CryptoWendyO (WRONG - has space after from:)
`

// TwitterSearch is a Cogito tool that bridges to the client's SearchTwitterWithArgs
type TwitterSearch struct {
	Client *client.Client
}

func (t *TwitterSearch) Name() string {
	return "search_twitter"
}

func (t *TwitterSearch) Description() string {
	return "Search Twitter using the provided query. Include operators like 'since:YYYY-MM-DD' (typically one day before today). Defaults to last 1 day if none provided. CRITICAL: Query ONLY 1 account per search using format 'from:username' with NO SPACE after 'from:' (e.g., 'from:JamesWynnReal', NOT 'from: JamesWynnReal'). Randomly sample accounts - do not exhaustively query all accounts. Use hashtags and keywords like '#BTC OR #ETH OR bitcoin OR ethereum' to find relevant tweets."
}

// Tool describes the tool for the underlying LLM provider (OpenAI-compatible)
func (t *TwitterSearch) Tool() openai.Tool {
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        t.Name(),
			Description: t.Description(),
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"query": {Type: jsonschema.String, Description: "Twitter advanced search query. CRITICAL: Use 'from:username' format with NO SPACE after 'from:' (e.g., 'from:JamesWynnReal (BTC OR Bitcoin)', NOT 'from: JamesWynnReal'). Include operators like 'since:YYYY-MM-DD' (typically one day before today), hashtags (#BTC), and keywords."},
				},
				Required: []string{"query"},
			},
		},
	}
}

// Execute executes the tool. Signature follows Cogito's ToolDefinitionInterface expectations.
// Expects params to include {"query": "..."} per the schema definition.
func (t *TwitterSearch) Execute(params map[string]any) (string, error) {
	var query string
	if q, ok := params["query"].(string); ok {
		query = q
	} else {
		b, _ := json.Marshal(params)
		query = string(b)
	}

	args := twitter.NewSearchArguments()
	args.Query = query

	docs, err := t.Client.SearchTwitterWithArgs(args)
	if err != nil {
		// Return error as a structured result string so the LLM can see what happened
		// This allows the agent to continue with partial data
		errorResult := map[string]any{
			"error":     true,
			"query":     query,
			"error_msg": err.Error(),
			"documents": []any{},
		}
		b, _ := json.Marshal(errorResult)
		return string(b), nil
	}

	// Return full documents - lean structure with useful metadata (username, created_at, likes, etc.)
	// Embedding and Score are omitempty so won't be serialized
	b, _ := json.Marshal(docs)
	return string(b), nil
}
