package types

type Config struct {
	BaseUrl string `envconfig:"GOPHER_CLIENT_URL" default:"https://data.gopher-ai.com/api"`
	Token   string `envconfig:"GOPHER_CLIENT_TOKEN"`
}
