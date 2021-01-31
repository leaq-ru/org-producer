package main

import (
	"context"
	"github.com/nnqq/scr-org-producer/config"
	"github.com/nnqq/scr-org-producer/healthz"
	"github.com/nnqq/scr-org-producer/logger"
	"github.com/nnqq/scr-org-producer/producer"
	"github.com/nnqq/scr-org-producer/stan"
	"log"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	sc, err := stan.NewConn(cfg.ServiceName, cfg.STAN.ClusterID, cfg.NATS.URL)
	logg.Must(err)

	go func() {
		logg.Must(healthz.NewHealthz(logg.ZL, cfg.HTTP.Port).Serve())
	}()

	logg.Must(producer.NewProducer(logg.ZL, sc).Do(ctx))
}
