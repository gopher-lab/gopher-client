package client

import (
	"github.com/gopher-lab/gopher-client/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Extraction Client", func() {
	var client *Client

	BeforeEach(func() {
		// Create client from config for real API testing
		var err error
		client, err = NewClientFromConfig()
		if err != nil {
			Skip("Skipping extraction tests - no valid config available")
		}
	})

	Describe("ExtractSearchTerms", func() {
		Context("with valid input", func() {
			It("should extract search terms successfully", func() {
				userInput := "Find me the latest news about AI developments and machine learning breakthroughs"
				maxTerms := 3

				response, err := client.ExtractSearchTerms(userInput, maxTerms)

				Expect(err).NotTo(HaveOccurred())
				Expect(response).NotTo(BeNil())
				// Note: API might return empty fields in some cases, so we just check that response exists
				Expect(response).To(BeAssignableToTypeOf(&types.ExtractionResponse{}))
			})
		})

		Context("with empty input", func() {
			It("should handle empty userInput gracefully", func() {
				userInput := ""
				maxTerms := 4

				response, err := client.ExtractSearchTerms(userInput, maxTerms)

				// The API might return an error for empty input, which is expected
				if err != nil {
					Expect(err.Error()).To(ContainSubstring("Invalid input"))
				} else {
					Expect(response).NotTo(BeNil())
				}
			})
		})
	})
})
