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
	llm    cogito.LLM
	Client *client.Client
}

var (
	ErrOpenAITokenRequired = errors.New("must supply an OPENAI_TOKEN")
)

const (
	DefaultModel        = "gpt-5-nano"
	DefaultOpenAIApiUrl = "https://api.openai.com/v1"
	DefaultPromptSuffix = "If no date range is specified, search the last 7 days."
)

// QueryOptions configures the behavior of the Query method
type QueryOptions struct {
	Schema       jsonschema.Definition
	Instructions string
	PromptSuffix string
	FinalPrompt  string
}

// QueryOption modifies QueryOptions
type QueryOption func(*QueryOptions)

// WithSchema sets a custom schema for structured extraction
func WithSchema(schema jsonschema.Definition) QueryOption {
	return func(opts *QueryOptions) {
		opts.Schema = schema
	}
}

// WithInstructions sets custom query instructions
func WithInstructions(instructions string) QueryOption {
	return func(opts *QueryOptions) {
		opts.Instructions = instructions
	}
}

// WithPromptSuffix sets a custom prompt suffix
func WithPromptSuffix(suffix string) QueryOption {
	return func(opts *QueryOptions) {
		opts.PromptSuffix = suffix
	}
}

// WithFinalPrompt sets a custom final prompt instruction
func WithFinalPrompt(prompt string) QueryOption {
	return func(opts *QueryOptions) {
		opts.FinalPrompt = prompt
	}
}

// DefaultSchema returns the default schema for topic sentiment analysis
func DefaultSchema() jsonschema.Definition {
	return jsonschema.Definition{
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
						"topic":           {Type: jsonschema.String, Description: "Topic"},
						"reasoning":       {Type: jsonschema.String, Description: "Reasoning about the topic"},
						"sentiment":       {Type: jsonschema.String, Description: "bullish, bearish, or neutral"},
						"top_influencers": {Type: jsonschema.Array, Items: &jsonschema.Definition{Type: jsonschema.String}},
					},
					Required: []string{"topic", "sentiment", "top_influencers"},
				},
			},
		},
		Required: []string{"topics"},
	}
}

// DefaultFinalPrompt returns the default final prompt instruction
func DefaultFinalPrompt() string {
	return "Return now only a JSON object with fields: topics (ordered by most relevant, each with topic, sentiment as bullish/bearish/neutral, and top_influencers array)."
}

// New creates a new Agent with the provided OpenAI token and model. Model defaults to gpt-5-nano.
func New(c *client.Client, openAIToken string) (*Agent, error) {
	if openAIToken == "" {
		return nil, ErrOpenAITokenRequired
	}
	llm := cogito.NewOpenAILLM(DefaultModel, openAIToken, DefaultOpenAIApiUrl)
	return &Agent{llm: llm, Client: c}, nil
}

// NewFromConfig creates a new Agent from config, defaulting the model to gpt-5-nano.
func NewFromConfig(c *client.Client) (*Agent, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	return New(c, cfg.OpenAIToken)
}

// NewAgentFromConfig creates an Agent from config in a single call.
// This convenience function creates both the underlying Client and Agent automatically,
// eliminating the need to manually create a client first.
func NewAgentFromConfig() (*Agent, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	c, err := client.NewClientFromConfig()
	if err != nil {
		return nil, err
	}
	ag, err := New(c, cfg.OpenAIToken)
	if err != nil {
		return nil, err
	}
	return ag, nil
}

// Query runs the agent with the provided natural language instruction.
// It uses Cogito with the TwitterSearch and WebSearch tools and extracts a structured Output.
// Query options can be provided to customize schema, instructions, and prompts.
func (a *Agent) Query(ctx context.Context, query string, opts ...QueryOption) (*types.Output, error) {
	// Apply options with defaults
	options := &QueryOptions{
		Schema:       DefaultSchema(),
		Instructions: TwitterQueryInstructions,
		PromptSuffix: DefaultPromptSuffix,
		FinalPrompt:  DefaultFinalPrompt(),
	}

	for _, opt := range opts {
		opt(options)
	}

	// Build full prompt with query, instructions, and suffix
	// Make it explicit that tools MUST be used
	fullPrompt := query + "\n\nIMPORTANT: You MUST use the available tools (search_twitter and search_web) to gather data. Do not provide an answer without using these tools first."
	if options.Instructions != "" {
		fullPrompt += "\n\n" + options.Instructions
	}
	if options.PromptSuffix != "" {
		fullPrompt += "\n\n" + options.PromptSuffix
	}

	fragment := cogito.NewEmptyFragment().
		AddMessage("user", fullPrompt)

	// Execute tools with the LLM
	improved, err := cogito.ExecuteTools(
		a.llm,
		fragment,
		cogito.WithContext(ctx),
		cogito.WithIterations(3),    // Increase iterations to allow multiple tool calls
		cogito.WithMaxAttempts(3),   // Allow multiple attempts
		cogito.WithForceReasoning(), // Force LLM to reason about tool usage
		cogito.WithTools(&TwitterSearch{Client: a.Client}, &WebSearch{Client: a.Client}),
	)
	if err != nil {
		return nil, err
	}

	// Ask the model to return structured JSON according to the final prompt
	result, err := a.llm.Ask(ctx, improved.AddMessage("user", options.FinalPrompt))
	if err != nil {
		return nil, err
	}

	out := &types.Output{}

	s := structures.Structure{
		Schema: options.Schema,
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
