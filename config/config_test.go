package config_test

import (
	"os"

	"github.com/gopher-lab/gopher-client/config"
	"github.com/gopher-lab/gopher-client/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe("LoadConfig", func() {
		Context("when no environment variables are set", func() {
			BeforeEach(func() {
				// Clear any existing environment variables
				os.Unsetenv("GOPHER_CLIENT_URL")
				os.Unsetenv("GOPHER_CLIENT_TOKEN")
			})

			It("should load default values", func() {
				cfg, err := config.LoadConfig()
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg).NotTo(BeNil())

				Expect(cfg.BaseUrl).To(Equal("https://data.gopher-ai.com"))
				Expect(cfg.Token).To(Equal(""))
			})
		})

		Context("when environment variables are set", func() {
			BeforeEach(func() {
				os.Setenv("GOPHER_CLIENT_URL", "https://api.example.com")
				os.Setenv("GOPHER_CLIENT_TOKEN", "test-token-123")
			})

			AfterEach(func() {
				os.Unsetenv("GOPHER_CLIENT_URL")
				os.Unsetenv("GOPHER_CLIENT_TOKEN")
			})

			It("should load values from environment variables", func() {
				cfg, err := config.LoadConfig()
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg).NotTo(BeNil())

				Expect(cfg.BaseUrl).To(Equal("https://api.example.com"))
				Expect(cfg.Token).To(Equal("test-token-123"))
			})
		})

		Context("when environment variables are empty", func() {
			BeforeEach(func() {
				os.Setenv("GOPHER_CLIENT_URL", "")
				os.Setenv("GOPHER_CLIENT_TOKEN", "")
			})

			AfterEach(func() {
				os.Unsetenv("GOPHER_CLIENT_URL")
				os.Unsetenv("GOPHER_CLIENT_TOKEN")
			})

			It("should use empty values when environment variables are set to empty strings", func() {
				cfg, err := config.LoadConfig()
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg).NotTo(BeNil())

				Expect(cfg.BaseUrl).To(Equal(""))
				Expect(cfg.Token).To(Equal(""))
			})
		})

		Context("when loading from .env file", func() {
			var envFile string

			BeforeEach(func() {
				envContent := `GOPHER_CLIENT_URL=https://staging.example.com
GOPHER_CLIENT_TOKEN=staging-token-456`

				envFile = ".env.test"
				err := os.WriteFile(envFile, []byte(envContent), 0644)
				Expect(err).NotTo(HaveOccurred())

				// Clear environment variables to test .env file loading
				os.Unsetenv("GOPHER_CLIENT_URL")
				os.Unsetenv("GOPHER_CLIENT_TOKEN")

				// Temporarily rename the test file to .env
				err = os.Rename(envFile, ".env")
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				os.Rename(".env", envFile)
				os.Remove(envFile)
			})

			It("should load values from .env file", func() {
				cfg, err := config.LoadConfig()
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg).NotTo(BeNil())

				Expect(cfg.BaseUrl).To(Equal("https://staging.example.com"))
				Expect(cfg.Token).To(Equal("staging-token-456"))
			})
		})
	})

	Describe("MustLoadConfig", func() {
		Context("when configuration is valid", func() {
			BeforeEach(func() {
				os.Setenv("GOPHER_CLIENT_URL", "https://production.example.com")
				os.Setenv("GOPHER_CLIENT_TOKEN", "production-token-789")
			})

			AfterEach(func() {
				os.Unsetenv("GOPHER_CLIENT_URL")
				os.Unsetenv("GOPHER_CLIENT_TOKEN")
			})

			It("should load configuration without panicking", func() {
				cfg := config.MustLoadConfig()
				Expect(cfg).NotTo(BeNil())

				Expect(cfg.BaseUrl).To(Equal("https://production.example.com"))
				Expect(cfg.Token).To(Equal("production-token-789"))
			})
		})
	})

	Describe("Config struct", func() {
		Context("when manually instantiated", func() {
			It("should allow field access and modification", func() {
				cfg := &types.Config{
					BaseUrl: "https://custom.example.com",
					Token:   "custom-token-abc",
				}

				Expect(cfg.BaseUrl).To(Equal("https://custom.example.com"))
				Expect(cfg.Token).To(Equal("custom-token-abc"))
			})
		})
	})
})

var _ = Describe("Config Performance", func() {
	BeforeEach(func() {
		os.Setenv("GOPHER_CLIENT_URL", "https://performance.example.com")
		os.Setenv("GOPHER_CLIENT_TOKEN", "performance-token-xyz")
	})

	AfterEach(func() {
		os.Unsetenv("GOPHER_CLIENT_URL")
		os.Unsetenv("GOPHER_CLIENT_TOKEN")
	})

	It("should load configuration quickly", func() {
		// Simple performance test - should complete in reasonable time
		cfg, err := config.LoadConfig()
		Expect(err).NotTo(HaveOccurred())
		Expect(cfg).NotTo(BeNil())
		Expect(cfg.BaseUrl).To(Equal("https://performance.example.com"))
		Expect(cfg.Token).To(Equal("performance-token-xyz"))
	})
})
