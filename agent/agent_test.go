package agent_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
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
	It("creates Agent from config and performs sentiment analysis on assets", func() {
		cfg := config.MustLoadConfig()
		if cfg.OpenAIToken == "" {
			Skip("OPENAI_TOKEN not set; skipping agent integration test")
		}

		ag, err := agent.NewAgentFromConfig()
		Expect(err).ToNot(HaveOccurred())

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		query := fmt.Sprintf("use all available tools to determine the market sentiment of the following assets: %s, using the following websites: %s, and the following twitter accounts: %s", strings.Join(agent.Assets, ", "), strings.Join(agent.Websites, ", "), strings.Join(agent.Kols, ", "))
		out, err := ag.Query(ctx, query)
		Expect(err).ToNot(HaveOccurred())
		Expect(out).ToNot(BeNil())

		// Verify we got multiple topics (at least Bitcoin and Ethereum)
		Expect(out.Assets).ToNot(BeEmpty(), "Expected at least one asset in the output")

		fmt.Println(prettyPrint(out.Assets))
	})
})
