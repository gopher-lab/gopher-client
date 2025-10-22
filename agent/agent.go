package agent

import (
	"context"
	"errors"
	"time"

	"github.com/mudler/cogito"
	"github.com/mudler/cogito/structures"
	"github.com/sashabaranov/go-openai/jsonschema"

	"github.com/gopher-lab/gopher-client/client"
	"github.com/gopher-lab/gopher-client/config"
	"github.com/gopher-lab/gopher-client/types"
)

type Agent struct {
	llm cogito.LLM
	c   *client.Client
}

var (
	ErrOpenAITokenRequired = errors.New("must supply an OPENAI_TOKEN")
)

const (
	DefaultModel        = "gpt-5-nano"
	DefaultOpenAIApiUrl = "https://api.openai.com/v1"
	DefaultPromptSuffix = "If no date range is specified, search the last 7 days."
)

// New creates a new Agent with the provided OpenAI token and model. Model defaults to gpt-5-nano.
func New(c *client.Client, openAIToken string) (*Agent, error) {
	if openAIToken == "" {
		return nil, ErrOpenAITokenRequired
	}
	llm := cogito.NewOpenAILLM(DefaultModel, openAIToken, DefaultOpenAIApiUrl)
	return &Agent{llm: llm, c: c}, nil
}

// NewFromConfig creates a new Agent from config, defaulting the model to gpt-5-nano.
func NewFromConfig(c *client.Client) (*Agent, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	return New(c, cfg.OpenAIToken)
}

// Query runs the agent with the provided natural language instruction.
// It uses Cogito with the TwitterSearch tool and extracts a structured Output.
func (a *Agent) Query(ctx context.Context, query string) (*types.Output, error) {
	// Include instruction supplement: operators/help and default date guidance
	fullPrompt := query + "\n\n" + TwitterQueryInstructions + "\n\n" + DefaultPromptSuffix

	fragment := cogito.NewEmptyFragment().
		AddMessage("user", fullPrompt)

	improved, err := cogito.ContentReview(
		a.llm,
		fragment,
		// cogito.EnableDeepContext,
		// cogito.EnableToolReEvaluator,
		// cogito.EnableToolReasoner,
		cogito.WithIterations(1),
		cogito.WithMaxAttempts(1),
		cogito.WithTools(&TwitterSearch{Client: a.c}),
	)
	if err != nil {
		return nil, err
	}

	// Ask the model to return only JSON for topics with sentiments and influencers
	result, err := a.llm.Ask(ctx, improved.AddMessage("user", "Return now only a JSON object with fields: topics (ordered by most relevant, each with topic, sentiment as bullish/bearish/neutral, and top_influencers array)."))
	if err != nil {
		return nil, err
	}

	out := &types.Output{}

	// Define schema for structured extraction
	schema := jsonschema.Definition{
		Type:                 jsonschema.Object,
		AdditionalProperties: false,
		Properties: map[string]jsonschema.Definition{
			"topics": {
				Type:        jsonschema.Array,
				Description: "Trending topics with sentiment and influencers",
				Items: &jsonschema.Definition{
					Type:                 jsonschema.Object,
					AdditionalProperties: false,
					Properties: map[string]jsonschema.Definition{
						"topic":           {Type: jsonschema.String},
						"sentiment":       {Type: jsonschema.String, Description: "bullish, bearish, or neutral"},
						"top_influencers": {Type: jsonschema.Array, Items: &jsonschema.Definition{Type: jsonschema.String}},
					},
					Required: []string{"topic", "sentiment", "top_influencers"},
				},
			},
		},
		Required: []string{"topics"},
	}

	s := structures.Structure{
		Schema: schema,
		Object: out,
	}

	// Provide a timeout context for extraction
	ctxExtract, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	if err := result.ExtractStructure(ctxExtract, a.llm, s); err != nil {
		// If extraction fails, still return empty Output to avoid nil
		// no-op
	}

	return out, nil
}
