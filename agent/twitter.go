package agent

import (
	"encoding/json"

	"github.com/gopher-lab/gopher-client/client"
	"github.com/masa-finance/tee-worker/api/args/twitter"
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
list:NASA/astronauts-in-space-now 	sent from a Twitter account in the NASA list astronauts-in-space-now
to:NASA 	a Tweet authored in reply to Twitter account "NASA".
@NASA 	mentioning Twitter account "NASA".
puppy filter:media 	containing "puppy" and an image or video.
puppy -filter:retweets 	containing "puppy", filtering out retweets
puppy filter:native_video 	containing "puppy" and an uploaded video, Amplify video, Periscope, or Vine.
puppy filter:periscope 	containing "puppy" and a Periscope video URL.
puppy filter:vine 	containing "puppy" and a Vine.
puppy filter:images 	containing "puppy" and links identified as photos, including third parties such as Instagram.
puppy filter:twimg 	containing "puppy" and a pic.twitter.com link representing one or more photos.
hilarious filter:links 	containing "hilarious" and linking to URL.
puppy url:amazon 	containing "puppy" and a URL with the word "amazon" anywhere within it.
superhero since:2015-12-21 	containing "superhero" and sent since date "2015-12-21" (year-month-day).
puppy until:2015-12-21 	containing "puppy" and sent before the date "2015-12-21".
movie -scary :) 	containing "movie", but not "scary", and with a positive attitude.
flight :( 	containing "flight" and with a negative attitude.
traffic ? 	containing "traffic" and asking a question.

Example:

altcoin or bitcoin :)

To search for the same day, you must subtract a day between since and until:
altcoin or bitcoin :) since:2025-03-23 until:2025-03-24

If no date range is specified, default to the last 7 days.
`

// TwitterSearch is a Cogito tool that bridges to the client's SearchTwitterWithArgs
type TwitterSearch struct {
	Client *client.Client
}

func (t *TwitterSearch) Name() string {
	return "search_twitter"
}

func (t *TwitterSearch) Description() string {
	return "Search Twitter using the provided query. Include operators, since/until. Defaults to last 7 days if none provided."
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

// Run executes the tool. Signature follows Cogito's Tool interface expectations.
// Expects params to include either {"query": "..."} or raw query under a heuristic.
func (t *TwitterSearch) Run(params map[string]any) (string, error) {
	var query string
	if q, ok := params["query"].(string); ok {
		query = q
	} else if q, ok := params["input"].(string); ok {
		query = q
	} else {
		b, _ := json.Marshal(params)
		query = string(b)
	}

	args := twitter.NewSearchArguments()
	args.Query = query

	docs, err := t.Client.SearchTwitterWithArgs(args)
	if err != nil {
		return "", err
	}

	b, _ := json.Marshal(docs)
	return string(b), nil
}
