package client

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Analysis Client", func() {
	var client *Client

	BeforeEach(func() {
		// Create client from config for real API testing
		var err error
		client, err = NewClientFromConfig()
		if err != nil {
			Skip("Skipping analysis tests - no valid config available")
		}
	})

	Describe("AnalyzeDataWithArgs", func() {
		Context("with valid input", func() {
			It("should analyze data successfully", func() {
				data := []string{
					"I am really excited about the economy",
					"I am bullish on cryptocurrency",
					"The market looks promising today",
				}
				prompt := "Analyze the sentiment of these tweets"
				model := "openai/gpt-4o-mini"

				response, err := client.AnalyzeDataWithArgs(data, prompt, model, false, nil, "")

				Expect(err).NotTo(HaveOccurred())
				Expect(response).NotTo(BeNil())
				Expect(response.Analysis).NotTo(BeEmpty())
				Expect(response.Reasoning).NotTo(BeEmpty())
				Expect(response.ModelUsed).To(Equal(model))
				Expect(response.TokensUsed).To(BeNumerically(">", 0))
				Expect(response.JobUUID).NotTo(BeEmpty())
			})
		})

		Context("with empty input", func() {
			It("should handle empty data array", func() {
				data := []string{}
				prompt := "Analyze this data"

				response, err := client.AnalyzeDataWithArgs(data, prompt, "", false, nil, "")

				// The API might return an error for empty data, which is expected
				if err != nil {
					Expect(err.Error()).To(ContainSubstring("Invalid input"))
				} else {
					Expect(response).NotTo(BeNil())
				}
			})
		})
	})

	Describe("AnalyzeData", func() {
		Context("with simple prompt", func() {
			It("should analyze data with default settings", func() {
				data := []string{
					"I love this product!",
					"This is amazing",
					"Great service",
				}
				prompt := "What is the sentiment?"

				response, err := client.AnalyzeData(data, prompt)

				Expect(err).NotTo(HaveOccurred())
				Expect(response).NotTo(BeNil())
				Expect(response.Analysis).NotTo(BeEmpty())
				Expect(response.Reasoning).NotTo(BeEmpty())
				Expect(response.ModelUsed).To(Equal("openai/gpt-4o-mini")) // Default model
				Expect(response.TokensUsed).To(BeNumerically(">", 0))
				Expect(response.JobUUID).NotTo(BeEmpty())
			})
		})
	})
})
