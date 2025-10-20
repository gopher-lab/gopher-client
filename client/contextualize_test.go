package client

import (
	"github.com/gopher-lab/gopher-client/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Contextualize Client", func() {
	var client *Client

	BeforeEach(func() {
		// Create client from config for real API testing
		var err error
		client, err = NewClientFromConfig()
		if err != nil {
			Skip("Skipping contextualize tests - no valid config available")
		}
	})

	Describe("ContextualizeQuery", func() {
		Context("with valid input", func() {
			It("should contextualize query successfully", func() {
				currentQuery := "What about the competition?"
				chatHistory := []types.ChatHistoryItem{
					{Query: "Tell me about Tesla stock", Timestamp: "2024-01-15T10:00:00Z"},
					{Query: "How is their production doing?", Timestamp: "2024-01-15T10:05:00Z"},
				}
				maxHistoryItems := 3

				response, err := client.ContextualizeQuery(currentQuery, chatHistory, maxHistoryItems)

				Expect(err).NotTo(HaveOccurred())
				Expect(response).NotTo(BeNil())
				Expect(response.ContextualizedQuery).NotTo(BeEmpty())
				Expect(response.OriginalQuery).To(Equal(currentQuery))
				Expect(response.Reasoning).NotTo(BeEmpty())
			})
		})

		Context("with empty input", func() {
			It("should handle empty currentQuery", func() {
				currentQuery := ""
				chatHistory := []types.ChatHistoryItem{
					{Query: "Previous query", Timestamp: "2024-01-15T10:00:00Z"},
				}
				maxHistoryItems := 5

				response, err := client.ContextualizeQuery(currentQuery, chatHistory, maxHistoryItems)

				// The API might return an error for empty currentQuery, which is expected
				if err != nil {
					Expect(err.Error()).To(ContainSubstring("Bad request"))
				} else {
					Expect(response).NotTo(BeNil())
				}
			})
		})
	})
})
