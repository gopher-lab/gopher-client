package agent

import (
	"encoding/json"
	"fmt"

	"github.com/gopher-lab/gopher-client/client"
	"github.com/masa-finance/tee-worker/v2/api/args/web"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

// WebSearch is a Cogito tool that bridges to the client's SearchTwitterWithArgs
type WebSearch struct {
	Client *client.Client
}

func (t *WebSearch) Name() string {
	return "search_web"
}

func (t *WebSearch) Description() string {
	return "Web search using the provided url"
}

// Tool describes the tool for the underlying LLM provider (OpenAI-compatible)
func (t *WebSearch) Tool() openai.Tool {
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        t.Name(),
			Description: t.Description(),
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"url": {Type: jsonschema.String, Description: "Web scrape url"},
				},
				Required: []string{"url"},
			},
		},
	}
}

// Execute executes the tool. Signature follows Cogito's ToolDefinitionInterface expectations.
// Expects params to include {"url": "..."} per the schema definition.
func (t *WebSearch) Execute(params map[string]any) (string, error) {
	var url string
	if q, ok := params["url"].(string); ok {
		url = q
	} else {
		b, _ := json.Marshal(params)
		url = string(b)
	}

	args := web.NewScraperArguments()
	args.URL = url

	docs, err := t.Client.ScrapeWebWithArgs(args)
	if err != nil {
		// If this is a timeout error, return a user-friendly message that the framework can use
		// The framework will convert this to a result string and continue execution
		if isTimeoutError(err) {
			return "", fmt.Errorf("web search timed out for URL %s: %v", url, err)
		}
		return "", err
	}

	// Return full documents - lean structure with useful metadata (title, canonicalUrl, markdown, etc.)
	// Embedding and Score are omitempty so won't be serialized
	b, _ := json.Marshal(docs)
	return string(b), nil
}
