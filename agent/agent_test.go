package agent_test

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gopher-lab/gopher-client/agent"
	"github.com/gopher-lab/gopher-client/client"
	"github.com/gopher-lab/gopher-client/config"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func prettyPrint(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

var _ = Describe("Agent integration", func() {
	It("creates client+agent from config and answers a query", func() {
		cfg := config.MustLoadConfig()
		if cfg.OpenAIToken == "" {
			Skip("OPENAI_TOKEN not set; skipping agent integration test")
		}

		c, err := client.NewClientFromConfig()
		Expect(err).ToNot(HaveOccurred())

		ag, err := agent.NewFromConfig(c)
		Expect(err).ToNot(HaveOccurred())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		query := "You have to find out a list of the most trending topics related to crypto and altcoins worldwide on Twitter for the last 3 days. Use the search_twitter tool."
		out, err := ag.Query(ctx, query)
		Expect(err).ToNot(HaveOccurred())
		Expect(out).ToNot(BeNil())

		fmt.Println(prettyPrint(out))
	})
})
