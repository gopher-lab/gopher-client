package agent

import (
	"encoding/json"
	"fmt"

	"github.com/gopher-lab/gopher-client/client"
	"github.com/masa-finance/tee-worker/v2/api/args/twitter"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

// TwitterQueryInstructions contains long-form guidance and examples for constructing Twitter queries.
// If no date range is specified by the user, searches should consider the last 7 days.
const TwitterQueryInstructions = `
watching now 	containing both "watching" and "now". This is the default operator.
"happy hour" 	containing the exact phrase "happy hour".
love OR hate 	containing either "love" or "hate" (or both).
beer -root 	containing "beer" but not "root".
#haiku 	containing the hashtag "haiku".
from:interior 	sent from Twitter account "interior".
to:NASA 	a Tweet authored in reply to Twitter account "NASA".
@NASA 	mentioning Twitter account "NASA".
superhero since:2015-12-21 	containing "superhero" and sent since date "2015-12-21" (year-month-day).
puppy until:2015-12-21 	containing "puppy" and sent before the date "2015-12-21".

To search for the same day, you must subtract a day between since and until:
altcoin or bitcoin :) since:2025-03-23 until:2025-03-24

If no date range is specified, default to the last 1 day.
`

// TwitterSearch is a Cogito tool that bridges to the client's SearchTwitterWithArgs
type TwitterSearch struct {
	Client *client.Client
}

func (t *TwitterSearch) Name() string {
	return "search_twitter"
}

func (t *TwitterSearch) Description() string {
	return "Search Twitter using the provided query. Include operators, since/until. Defaults to last 1 day if none provided."
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
					"query": {Type: jsonschema.String, Description: "Twitter advanced search query (with operators, since/until)"},
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
		// If this is a timeout error, return a user-friendly message that the framework can use
		// The framework will convert this to a result string and continue execution
		if isTimeoutError(err) {
			return "", fmt.Errorf("twitter search timed out for query %s: %v", query, err)
		}
		return "", err
	}

	// Return full documents - lean structure with useful metadata (username, created_at, likes, etc.)
	// Embedding and Score are omitempty so won't be serialized
	b, _ := json.Marshal(docs)
	return string(b), nil
}
