package config

type Config struct {
	HTTP     http
	MongoDB  mongodb
	DaData   daData
	Service  service
	LogLevel string `envconfig:"LOGLEVEL"`
}

type http struct {
	Port string `envconfig:"HTTP_PORT"`
}

type mongodb struct {
	URL string `envconfig:"MONGODB_URL"`
}

type service struct {
	Org string `envconfig:"SERVICE_ORG"`
}

type daData struct {
	Tokens string `envconfig:"DADATA_TOKENS"`
}
