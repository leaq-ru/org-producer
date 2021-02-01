package config

type Config struct {
	ServiceName string
	HTTP        http
	STAN        stan
	NATS        nats
	MongoDB     mongoDB
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

type mongoDB struct {
	URL string `envconfig:"MONGODB_URL"`
}
