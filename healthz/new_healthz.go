package healthz

import "github.com/rs/zerolog"

func NewHealthz(logger zerolog.Logger, port string) Healthz {
	return Healthz{
		Logger: logger,
		Port:   port,
	}
}
