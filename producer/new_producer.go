package producer

import (
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog"
)

func NewProducer(logger zerolog.Logger, stanConn stan.Conn) Producer {
	return Producer{
		logger:   logger,
		stanConn: stanConn,
	}
}
