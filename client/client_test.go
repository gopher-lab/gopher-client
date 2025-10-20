package client

import (
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	Describe("Constructor Functions", func() {
		Context("NewClient", func() {
			It("should create a client with provided URL and token", func() {
				baseURL := "https://api.example.com"
				token := "test-token-123"

				client := NewClient(baseURL, token)

				Expect(client).NotTo(BeNil())
				Expect(client.BaseURL).To(Equal(baseURL))
				Expect(client.Token).To(Equal(token))
			})

			It("should create a client with empty token", func() {
				baseURL := "https://api.example.com"
				token := ""

				client := NewClient(baseURL, token)

				Expect(client).NotTo(BeNil())
				Expect(client.BaseURL).To(Equal(baseURL))
				Expect(client.Token).To(Equal(""))
			})
		})

		Context("NewClientFromConfig", func() {
			BeforeEach(func() {
				// Set test environment variables
				os.Setenv("GOPHER_CLIENT_URL", "https://test.example.com")
				os.Setenv("GOPHER_CLIENT_TOKEN", "test-token-456")
			})

			AfterEach(func() {
				// Clean up environment variables
				os.Unsetenv("GOPHER_CLIENT_URL")
				os.Unsetenv("GOPHER_CLIENT_TOKEN")
			})

			It("should create a client from config successfully", func() {
				client, err := NewClientFromConfig()

				Expect(err).NotTo(HaveOccurred())
				Expect(client).NotTo(BeNil())
				Expect(client.BaseURL).To(Equal("https://test.example.com"))
				Expect(client.Token).To(Equal("test-token-456"))
			})
		})

		Context("MustNewClientFromConfig", func() {
			BeforeEach(func() {
				os.Setenv("GOPHER_CLIENT_URL", "https://must-test.example.com")
				os.Setenv("GOPHER_CLIENT_TOKEN", "must-test-token")
			})

			AfterEach(func() {
				os.Unsetenv("GOPHER_CLIENT_URL")
				os.Unsetenv("GOPHER_CLIENT_TOKEN")
			})

			It("should create a client without panicking", func() {
				client := MustNewClientFromConfig()

				Expect(client).NotTo(BeNil())
				Expect(client.BaseURL).To(Equal("https://must-test.example.com"))
				Expect(client.Token).To(Equal("must-test-token"))
			})
		})
	})

	Describe("Helper Functions", func() {
		Context("getErrorFromResponse", func() {
			It("should return nil for response without error", func() {
				response := `{"data": "some data"}`

				err := getErrorFromResponse([]byte(response))

				Expect(err).To(BeNil())
			})

			It("should return error for response with error field", func() {
				response := `{"error": "something went wrong"}`

				err := getErrorFromResponse([]byte(response))

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("job errored: something went wrong"))
			})

			It("should return nil for invalid JSON", func() {
				response := `invalid json`

				err := getErrorFromResponse([]byte(response))

				Expect(err).To(BeNil())
			})
		})
	})

	Describe("HTTP Request Methods", func() {
		var client *Client

		BeforeEach(func() {
			client = NewClient("https://api.example.com", "test-token")
		})

		Context("doRequest", func() {
			It("should handle invalid URL", func() {
				// This is a basic test - in a real scenario you'd use httptest.Server
				// For simplicity, we're just testing error handling with invalid URL
				requestBody := []byte(`{"query": "test"}`)

				// Since we can't easily mock http.Client without more complex setup,
				// we'll test that the method doesn't panic with invalid URL
				_, err := client.doRequest("invalid-url", requestBody)

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("failed to do POST request"))
			})
		})

		Context("doStatusRequest", func() {
			It("should handle invalid URL", func() {
				url := "invalid-url"

				_, err := client.doStatusRequest(url)

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("failed to do GET request"))
			})
		})

		Context("doResultRequest", func() {
			It("should handle invalid URL", func() {
				url := "invalid-url"
				var receiver interface{}

				err := client.doResultRequest(url, receiver)

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("failed to do GET request"))
			})
		})

		Context("doImmediateRequest", func() {
			It("should handle invalid URL", func() {
				url := "invalid-url"
				requestBody := []byte(`{"query": "test"}`)
				var receiver interface{}

				err := client.doImmediateRequest(url, requestBody, receiver)

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("failed to do POST request"))
			})
		})
	})

	Describe("Public API Methods", func() {
		var client *Client

		BeforeEach(func() {
			client = NewClient("https://api.example.com", "test-token")
		})

		Context("GetJobStatus", func() {
			It("should construct correct URL", func() {
				jobID := "test-job-123"
				expectedURL := "https://api.example.com/v1/search/live/status/test-job-123"

				// Test URL construction by checking error message
				_, err := client.GetJobStatus(jobID)

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("failed to do GET request"))
				// The error should contain the constructed URL
				Expect(err.Error()).To(ContainSubstring(expectedURL))
			})
		})

		Context("GetResult", func() {
			It("should construct correct URL", func() {
				jobID := "test-job-456"
				expectedURL := "https://api.example.com/v1/search/live/result/test-job-456"
				var receiver interface{}

				// Test URL construction by checking error message
				err := client.GetResult(jobID, receiver)

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("failed to do GET request"))
				// The error should contain the constructed URL
				Expect(err.Error()).To(ContainSubstring(expectedURL))
			})
		})

		Context("WaitForJobCompletion", func() {
			It("should timeout when job doesn't complete", func() {
				jobID := "test-job-timeout"
				timeout := 100 * time.Millisecond

				// This should timeout quickly since we're using an invalid URL
				_, err := client.WaitForJobCompletion(jobID, timeout)

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("timed out after"))
			})
		})

		Context("WaitForJobCompletionWithDefaultTimeout", func() {
			BeforeEach(func() {
				os.Setenv("GOPHER_CLIENT_TIMEOUT", "500ms")
			})

			AfterEach(func() {
				os.Unsetenv("GOPHER_CLIENT_TIMEOUT")
			})

			It("should use default timeout from config", func() {
				jobID := "test-job-default-timeout"

				// This should timeout quickly since we're using an invalid URL
				_, err := client.WaitForJobCompletionWithDefaultTimeout(jobID)

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("timed out after"))
			})
		})
	})

	Describe("Client Structure", func() {
		Context("when created", func() {
			It("should have correct field values", func() {
				baseURL := "https://example.com"
				token := "abc123"

				client := NewClient(baseURL, token)

				Expect(client.BaseURL).To(Equal(baseURL))
				Expect(client.Token).To(Equal(token))
			})
		})
	})
})
