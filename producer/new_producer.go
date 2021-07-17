package producer

import (
	"github.com/leaq-ru/org-producer/state"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog"
)

func NewProducer(logger zerolog.Logger, stanConn stan.Conn, stateModel state.Model) Producer {
	return Producer{
		logger:     logger,
		stanConn:   stanConn,
		stateModel: stateModel,
	}
}
