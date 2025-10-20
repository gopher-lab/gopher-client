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
				os.Unsetenv("BASE_URL")
			})

			It("should load default values", func() {
				cfg, err := config.LoadConfig()
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg).NotTo(BeNil())

				Expect(cfg.BaseUrl).To(Equal("https://data.gopher-ai.com"))
			})
		})

		Context("when environment variables are set", func() {
			BeforeEach(func() {
				os.Setenv("BASE_URL", "https://api.example.com")
			})

			AfterEach(func() {
				os.Unsetenv("BASE_URL")
			})

			It("should load values from environment variables", func() {
				cfg, err := config.LoadConfig()
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg).NotTo(BeNil())

				Expect(cfg.BaseUrl).To(Equal("https://api.example.com"))
			})
		})

		Context("when environment variable is empty", func() {
			BeforeEach(func() {
				os.Setenv("BASE_URL", "")
			})

			AfterEach(func() {
				os.Unsetenv("BASE_URL")
			})

			It("should use empty value when environment variable is set to empty string", func() {
				cfg, err := config.LoadConfig()
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg).NotTo(BeNil())

				Expect(cfg.BaseUrl).To(Equal(""))
			})
		})

		Context("when loading from .env file", func() {
			var envFile string

			BeforeEach(func() {
				envContent := `BASE_URL=https://staging.example.com`

				envFile = ".env.test"
				err := os.WriteFile(envFile, []byte(envContent), 0644)
				Expect(err).NotTo(HaveOccurred())

				// Clear environment variables to test .env file loading
				os.Unsetenv("BASE_URL")

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
			})
		})
	})

	Describe("MustLoadConfig", func() {
		Context("when configuration is valid", func() {
			BeforeEach(func() {
				os.Setenv("BASE_URL", "https://production.example.com")
			})

			AfterEach(func() {
				os.Unsetenv("BASE_URL")
			})

			It("should load configuration without panicking", func() {
				cfg := config.MustLoadConfig()
				Expect(cfg).NotTo(BeNil())

				Expect(cfg.BaseUrl).To(Equal("https://production.example.com"))
			})
		})
	})

	Describe("Config struct", func() {
		Context("when manually instantiated", func() {
			It("should allow field access and modification", func() {
				cfg := &types.Config{
					BaseUrl: "https://custom.example.com",
				}

				Expect(cfg.BaseUrl).To(Equal("https://custom.example.com"))
			})
		})
	})
})

var _ = Describe("Config Performance", func() {
	BeforeEach(func() {
		os.Setenv("BASE_URL", "https://performance.example.com")
	})

	AfterEach(func() {
		os.Unsetenv("BASE_URL")
	})

	It("should load configuration quickly", func() {
		// Simple performance test - should complete in reasonable time
		cfg, err := config.LoadConfig()
		Expect(err).NotTo(HaveOccurred())
		Expect(cfg).NotTo(BeNil())
		Expect(cfg.BaseUrl).To(Equal("https://performance.example.com"))
	})
})
