package types

type Config struct {
	BaseUrl string `envconfig:"GOPHER_CLIENT_URL" default:"https://data.gopher-ai.com"`
	Token   string `envconfig:"GOPHER_CLIENT_TOKEN"`
}
