package producer

import (
	"github.com/leaq-ru/org-producer/state"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog"
)

type Producer struct {
	logger     zerolog.Logger
	stanConn   stan.Conn
	stateModel state.Model
}
