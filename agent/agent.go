package agent

import (
	"context"
	"errors"
	"fmt"
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
			"assets": {
				Type:        jsonschema.Array,
				Description: "Track the market sentiment of assets, such as Bitcoin, Ethereum, and other cryptocurrencies.",
				Items: &jsonschema.Definition{
					Type:                 jsonschema.Object,
					AdditionalProperties: false,
					Properties: map[string]jsonschema.Definition{
						"asset":     {Type: jsonschema.String, Description: "Asset name"},
						"reasoning": {Type: jsonschema.String, Description: "Brief reasoning about the sentiment of the asset"},
						"sentiment": {Type: jsonschema.Integer, Description: "Numeric sentiment score from 1-100, where 100 is the most bullish and 1 is the most bearish"},
					},
					Required: []string{"asset", "reasoning", "sentiment"},
				},
			},
		},
		Required: []string{"assets"},
	}
}

// DefaultFinalPrompt returns the default final prompt instruction
func DefaultFinalPrompt() string {
	return "Return now only a JSON object with fields that match the supplied schema."
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
	fullPrompt := query
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
		cogito.WithIterations(2),    // Allow multiple tool calls in sequence
		cogito.WithMaxAttempts(2),   // Allow multiple attempts for tool selection
		cogito.WithForceReasoning(), // Force LLM to reason about tool usage
		cogito.WithTools(&WebSearch{Client: a.Client}, &TwitterSearch{Client: a.Client}),
	)
	if err != nil {
		return nil, err
	}

	// Ask the model to return structured JSON according to the final prompt
	result, err := a.llm.Ask(ctx, improved.AddMessage("user", options.FinalPrompt))
	if err != nil {
		return nil, err
	}

	// Log the raw LLM response for debugging
	if lastMsg := result.LastMessage(); lastMsg != nil {
		fmt.Printf("DEBUG: LLM response before extraction: %s\n", lastMsg.Content)
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
		// If extraction fails, still return what we have (might be partial data)
		// Log the error but don't fail completely - we might have some data
		fmt.Printf("DEBUG: Extraction error (but returning partial results): %v\n", err)
		// Continue to return out even if extraction failed - it might have partial data
	}

	// Check if we got any data
	if len(out.Assets) == 0 {
		// If no topics extracted, log this for debugging
		fmt.Printf("DEBUG: No topics extracted. Tools called: %d\n", len(improved.Status.ToolsCalled))
	}

	return out, nil
}
