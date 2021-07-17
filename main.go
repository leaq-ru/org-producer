package main

import (
	"context"
	"github.com/leaq-ru/org-producer/config"
	"github.com/leaq-ru/org-producer/healthz"
	"github.com/leaq-ru/org-producer/logger"
	"github.com/leaq-ru/org-producer/mongo"
	"github.com/leaq-ru/org-producer/producer"
	"github.com/leaq-ru/org-producer/stan"
	"github.com/leaq-ru/org-producer/state"
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

	db, err := mongo.NewConn(ctx, cfg.ServiceName, cfg.MongoDB.URL)

	go func() {
		logg.Must(healthz.NewHealthz(logg.ZL, cfg.HTTP.Port).Serve())
	}()

	logg.Must(producer.NewProducer(logg.ZL, sc, state.NewModel(db)).Do(ctx))
}
