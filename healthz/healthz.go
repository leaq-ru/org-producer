package healthz

import "github.com/rs/zerolog"

type Healthz struct {
	Logger zerolog.Logger
	Port   string
}
