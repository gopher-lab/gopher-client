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
	"github.com/gopher-lab/gopher-client/log"
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
	DefaultPromptSuffix = "If no date range is specified, search the last 1 day"
	// DefaultDataCollectionInstructions provides explicit guidance about trying all sources
	DefaultDataCollectionInstructions = `
CRITICAL: You MUST use the available tools to gather data. Do NOT attempt to answer without using tools.

You have access to these tools:
1. search_web - Search and scrape web pages
2. search_twitter - Search Twitter for tweets from specific accounts

REQUIRED ACTIONS:
1. You MUST use search_web to fetch data from the provided websites
2. You MUST use search_twitter to gather sentiment from Twitter accounts
3. Do NOT skip using tools - they are required to complete this task

IMPORTANT: You must attempt to gather data from ALL available sources, even if some fail.
- Try ALL websites provided, even if some URLs timeout or return errors
- Execute Twitter searches for sentiment analysis, but IMPORTANT: Randomly sample Twitter accounts
  - Randomly select accounts from the provided list (typically 3-6 accounts for good coverage)
  - DO NOT exhaustively query all accounts - a random sample is sufficient
  - Query format: 'from:username (BTC OR Bitcoin OR ETH OR Ethereum OR SOL OR Solana)' - NO SPACE after 'from:' (e.g., 'from:JamesWynnReal', NOT 'from: JamesWynnReal')
  - Use hashtags and keywords: '#BTC OR #ETH OR bitcoin OR ethereum'
  - For faster execution, you can batch multiple queries using the 'queries' array parameter - this will execute them concurrently
  - Example single query: 'from:JamesWynnReal (BTC OR Bitcoin OR ETH OR Ethereum) since:2025-11-03'
  - Example batch queries: ["from:JamesWynnReal (BTC OR Bitcoin) since:2025-11-03", "from:CryptoWendyO (ETH OR Ethereum) since:2025-11-03", "from:PeterLBrandt (SOL OR Solana) since:2025-11-03"]
- Continue with remaining sources even if earlier sources fail
- Partial data is acceptable - gather what you can from each source
- Use multiple iterations to systematically collect data from all sources before synthesizing results
`
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
						"sentiment": {Type: jsonschema.Integer, Description: "Numeric sentiment score from 0-100. Scale: 0 = most bearish, 50 = neutral, 100 = most bullish. Each asset should have a distinct score based on its unique data and sentiment signals."},
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
	return "Return now only a JSON object with fields that match the supplied schema. IMPORTANT: Each asset must have a distinct sentiment score based on its unique data and signals - do not use the same sentiment value for different assets. Analyze each asset independently."
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
	fullPrompt += "\n\n" + DefaultDataCollectionInstructions
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
		cogito.WithIterations(6),    // Increased to allow more Twitter account sampling
		cogito.WithMaxAttempts(1),   // Allow multiple attempts for tool selection
		cogito.WithForceReasoning(), // Force LLM to reason about tool usage
		cogito.WithTools(&WebSearch{Client: a.Client}, &TwitterSearch{Client: a.Client}),
	)
	if err != nil {
		// Check if it's ErrNoToolSelected and provide a more helpful error
		if errors.Is(err, cogito.ErrNoToolSelected) {
			return nil, fmt.Errorf("LLM did not select any tools. This task requires using search_web and search_twitter tools to gather data: %w", err)
		}
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
		log.Error("Extraction error", err)
	}

	return out, nil
}
