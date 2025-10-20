package types

type Config struct {
	BaseUrl string `envconfig:"BASE_URL" default:"https://data.gopher-ai.com"`
}
