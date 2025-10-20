package client

import (
	"github.com/masa-finance/tee-worker/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Similarity Client", func() {
	var client *Client

	BeforeEach(func() {
		// Create client from config for real API testing
		var err error
		client, err = NewClientFromConfig()
		if err != nil {
			Skip("Skipping similarity tests - no valid config available")
		}
	})

	Describe("PerformSimilaritySearch", func() {
		Context("with valid input", func() {
			It("should perform similarity search successfully", func() {
				query := "artificial intelligence"
				sources := []types.Source{types.WebSource, types.TwitterSource}
				keywords := []string{"AI", "machine learning"}
				operator := "and"
				maxResults := 10

				results, err := client.PerformSimilaritySearch(query, sources, keywords, operator, maxResults)

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
				keywords := []string{"test"}
				operator := "and"
				maxResults := 5

				results, err := client.PerformSimilaritySearch(query, sources, keywords, operator, maxResults)

				// The API might return an error for empty query, which is expected
				if err != nil {
					Expect(err.Error()).To(ContainSubstring("Invalid input"))
				} else {
					Expect(results).NotTo(BeNil())
				}
			})
		})
	})
})
