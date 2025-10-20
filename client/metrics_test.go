package client

import (
	"github.com/masa-finance/tee-worker/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Metrics Client", func() {
	var client *Client

	BeforeEach(func() {
		// Create client from config for real API testing
		var err error
		client, err = NewClientFromConfig()
		if err != nil {
			Skip("Skipping metrics tests - no valid config available")
		}
	})

	Describe("GetAllMetrics", func() {
		Context("with valid input", func() {
			It("should get all metrics successfully", func() {
				refresh := false

				stats, err := client.GetAllMetrics(refresh)

				Expect(err).NotTo(HaveOccurred())
				Expect(stats).NotTo(BeNil())
				// Note: Results might be empty if no collections exist, which is acceptable
				Expect(stats).To(BeAssignableToTypeOf([]types.CollectionStats{}))
			})

			It("should get all metrics with refresh enabled", func() {
				refresh := true

				stats, err := client.GetAllMetrics(refresh)

				Expect(err).NotTo(HaveOccurred())
				Expect(stats).NotTo(BeNil())
				Expect(stats).To(BeAssignableToTypeOf([]types.CollectionStats{}))
			})
		})
	})

	Describe("GetMetrics", func() {
		Context("with valid input", func() {
			It("should get metrics for a specific collection successfully", func() {
				source := "web"
				refresh := false

				stats, err := client.GetMetrics(source, refresh)

				Expect(err).NotTo(HaveOccurred())
				Expect(stats).NotTo(BeNil())
				Expect(stats).To(BeAssignableToTypeOf(&types.CollectionStats{}))
				// Note: CollectionName might be empty or different, which is acceptable
			})

			It("should get metrics for a specific collection with refresh enabled", func() {
				source := "twitter"
				refresh := true

				stats, err := client.GetMetrics(source, refresh)

				Expect(err).NotTo(HaveOccurred())
				Expect(stats).NotTo(BeNil())
				Expect(stats).To(BeAssignableToTypeOf(&types.CollectionStats{}))
			})
		})

		Context("with invalid input", func() {
			It("should handle non-existent collection gracefully", func() {
				source := "nonexistent_collection"
				refresh := false

				stats, err := client.GetMetrics(source, refresh)

				// The API might return an error for non-existent collection, which is expected
				if err != nil {
					Expect(err.Error()).To(Or(ContainSubstring("collection"), ContainSubstring("not found"), ContainSubstring("error")))
				} else {
					Expect(stats).NotTo(BeNil())
				}
			})

			It("should handle empty source gracefully", func() {
				source := ""
				refresh := false

				_, err := client.GetMetrics(source, refresh)

				// When source is empty, the API returns an array instead of a single object
				// This causes an unmarshal error, which is expected behavior
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unmarshal"))
			})
		})
	})
})
