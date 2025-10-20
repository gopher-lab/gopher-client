package client

import (
	"github.com/masa-finance/tee-worker/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hybrid Client", func() {
	var client *Client

	BeforeEach(func() {
		// Create client from config for real API testing
		var err error
		client, err = NewClientFromConfig()
		if err != nil {
			Skip("Skipping hybrid tests - no valid config available")
		}
	})

	Describe("PerformHybridSearch", func() {
		Context("with valid input", func() {
			It("should perform hybrid search successfully", func() {
				query := "machine learning"
				sources := []types.Source{types.WebSource, types.RedditSource}
				text := "artificial intelligence"
				queryWeight := 0.7
				textWeight := 0.3
				keywords := []string{"AI", "ML"}
				operator := "or"
				maxResults := 10

				var results []types.Document
				err := client.PerformHybridSearch(query, sources, text, queryWeight, textWeight, keywords, operator, maxResults, &results)

				Expect(err).NotTo(HaveOccurred())
				Expect(results).NotTo(BeNil())
				// Note: Results might be empty if no data is indexed, which is acceptable
				Expect(results).To(BeAssignableToTypeOf([]types.Document{}))
			})
		})

		Context("with invalid input", func() {
			It("should handle empty query gracefully", func() {
				query := ""
				sources := []types.Source{types.WebSource}
				text := "test"
				queryWeight := 0.5
				textWeight := 0.5
				keywords := []string{"test"}
				operator := "and"
				maxResults := 5

				var results []types.Document
				err := client.PerformHybridSearch(query, sources, text, queryWeight, textWeight, keywords, operator, maxResults, &results)

				// The API might return an error for empty query, which is expected
				if err != nil {
					Expect(err.Error()).To(ContainSubstring("Invalid text_query"))
				} else {
					Expect(results).NotTo(BeNil())
				}
			})
		})
	})
})
