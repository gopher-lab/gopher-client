package agent

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gopher-lab/gopher-client/client"
	"github.com/masa-finance/tee-worker/v2/api/args/twitter"
	"github.com/masa-finance/tee-worker/v2/api/types"
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
Single query:
- from:JamesWynnReal (BTC OR Bitcoin OR ETH OR Ethereum) since:2025-11-03
- from:CryptoWendyO #BTC OR #ETH since:2025-11-03
- from:VitalikButerin (ethereum OR ETH) since:2025-11-03

Batch queries (executed concurrently):
- Use the "queries" parameter with an array of queries for faster parallel execution
- Example: ["from:JamesWynnReal (BTC OR Bitcoin) since:2025-11-03", "from:CryptoWendyO (ETH OR Ethereum) since:2025-11-03", "from:PeterLBrandt (SOL OR Solana) since:2025-11-03"]

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
	return "Search Twitter using the provided query or queries. Include operators like 'since:YYYY-MM-DD' (typically one day before today). Defaults to last 1 day if none provided. You can provide a single 'query' or multiple 'queries' as an array. For multiple queries, they will be executed concurrently. CRITICAL: Use 'from:username' format with NO SPACE after 'from:' (e.g., 'from:JamesWynnReal', NOT 'from: JamesWynnReal'). Randomly sample accounts - do not exhaustively query all accounts. Use hashtags and keywords like '#BTC OR #ETH OR bitcoin OR ethereum' to find relevant tweets."
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
					"query": {
						Type:        jsonschema.String,
						Description: "Twitter advanced search query. CRITICAL: Use 'from:username' format with NO SPACE after 'from:' (e.g., 'from:JamesWynnReal (BTC OR Bitcoin)', NOT 'from: JamesWynnReal'). Include operators like 'since:YYYY-MM-DD' (typically one day before today), hashtags (#BTC), and keywords.",
					},
					"queries": {
						Type:        jsonschema.Array,
						Description: "Array of Twitter advanced search queries to execute concurrently. Each query should follow the same format as 'query'. When provided, all queries will be executed in parallel for faster results.",
						Items: &jsonschema.Definition{
							Type: jsonschema.String,
						},
					},
				},
			},
		},
	}
}

// Execute executes the tool. Signature follows Cogito's ToolDefinitionInterface expectations.
// Supports both single query {"query": "..."} and multiple queries {"queries": ["...", "..."]}
// Multiple queries are executed concurrently using goroutines.
func (t *TwitterSearch) Execute(params map[string]any) (string, error) {
	// Check if multiple queries are provided
	if queries, ok := params["queries"].([]any); ok && len(queries) > 0 {
		return t.executeConcurrent(queries)
	}

	// Single query execution (backward compatible)
	var query string
	if q, ok := params["query"].(string); ok {
		query = q
	} else {
		b, _ := json.Marshal(params)
		query = string(b)
	}

	return t.executeSingle(query)
}

// executeSingle executes a single Twitter search query
func (t *TwitterSearch) executeSingle(query string) (string, error) {
	args := twitter.NewSearchArguments()
	args.Query = query

	docs, err := t.Client.SearchTwitterWithArgs(args)
	if err != nil {
		// Return error as a structured result string so the LLM can see what happened
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
	b, _ := json.Marshal(docs)
	return string(b), nil
}

// executeConcurrent executes multiple Twitter search queries concurrently using goroutines
func (t *TwitterSearch) executeConcurrent(queries []any) (string, error) {
	type queryResult struct {
		query     string
		documents []types.Document
		err       error
	}

	// Convert queries to strings
	queryStrings := make([]string, 0, len(queries))
	for _, q := range queries {
		if qStr, ok := q.(string); ok {
			queryStrings = append(queryStrings, qStr)
		}
	}

	if len(queryStrings) == 0 {
		return "[]", nil
	}

	// Channel to collect results
	results := make(chan queryResult, len(queryStrings))
	var wg sync.WaitGroup

	// Execute all queries concurrently
	for _, query := range queryStrings {
		wg.Add(1)
		go func(q string) {
			defer wg.Done()
			args := twitter.NewSearchArguments()
			args.Query = q

			docs, err := t.Client.SearchTwitterWithArgs(args)
			result := queryResult{
				query:     q,
				documents: docs,
				err:       err,
			}
			results <- result
		}(query)
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect all results
	allDocs := []types.Document{}
	errors := []map[string]any{}

	for result := range results {
		if result.err != nil {
			// Store error info but continue processing other results
			errors = append(errors, map[string]any{
				"error":     true,
				"query":     result.query,
				"error_msg": result.err.Error(),
			})
		} else {
			allDocs = append(allDocs, result.documents...)
		}
	}

	// Build response with aggregated results and any errors
	response := map[string]any{
		"documents":          allDocs,
		"total_queries":      len(queryStrings),
		"successful_queries": len(queryStrings) - len(errors),
		"failed_queries":     len(errors),
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	b, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("failed to marshal results: %w", err)
	}

	return string(b), nil
}
