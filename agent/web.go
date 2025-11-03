package agent

import (
	"encoding/json"

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
// Expects params to include either {"url": "..."} or raw url under a heuristic.
func (t *WebSearch) Execute(params map[string]any) (string, error) {
	var url string
	if q, ok := params["url"].(string); ok {
		url = q
	} else if q, ok := params["url"].(string); ok {
		url = q
	} else {
		b, _ := json.Marshal(params)
		url = string(b)
	}

	args := web.NewScraperArguments()
	args.URL = url

	docs, err := t.Client.ScrapeWebWithArgs(args)
	if err != nil {
		return "", err
	}

	// Extract both content and metadata.markdown from documents
	type docResult struct {
		Content  string `json:"content"`
		Markdown string `json:"markdown,omitempty"`
	}

	results := make([]docResult, 0, len(docs))
	for _, d := range docs {
		markdown := ""
		if d.Metadata != nil {
			if md, ok := d.Metadata["markdown"].(string); ok {
				markdown = md
			}
		}
		results = append(results, docResult{
			Content:  d.Content, // llm summary
			Markdown: markdown,
		})
	}

	b, _ := json.Marshal(results)
	return string(b), nil
}
