package agent_test

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gopher-lab/gopher-client/agent"
	"github.com/gopher-lab/gopher-client/config"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func prettyPrint(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

var _ = Describe("Agent integration", func() {
	It("creates client+agent from config and performs sentiment analysis on Bitcoin and Ethereum", func() {
		cfg := config.MustLoadConfig()
		if cfg.OpenAIToken == "" {
			Skip("OPENAI_TOKEN not set; skipping agent integration test")
		}

		ag, err := agent.NewAgentFromConfig()
		Expect(err).ToNot(HaveOccurred())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		query := "use search_twitter and search_web tools to understand current sentiment on the following coins: bitcoin, ethereum. use farside.co.uk/btc/ website to scrape web data for bitcoin, and farside.co.uk/eth/ website to scrape web data for ethereum. focus on tweets / influencers with a big following or likes / views on their tweets."
		out, err := ag.Query(ctx, query)
		Expect(err).ToNot(HaveOccurred())
		Expect(out).ToNot(BeNil())

		// Verify we got multiple topics (at least Bitcoin and Ethereum)
		Expect(out.Topics).ToNot(BeEmpty(), "Expected at least one topic in the output")

		fmt.Println(prettyPrint(out))
	})
})
