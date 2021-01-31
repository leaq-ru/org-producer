package config

type Config struct {
	ServiceName string
	HTTP        http
	STAN        stan
	NATS        nats
	LogLevel    string `envconfig:"LOGLEVEL"`
}

type http struct {
	Port string `envconfig:"HTTP_PORT"`
}

type stan struct {
	ClusterID string `envconfig:"STAN_CLUSTERID"`
}

type nats struct {
	URL string `envconfig:"NATS_URL"`
}
