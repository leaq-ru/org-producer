package main

import (
	"context"
	"github.com/nnqq/scr-org-producer/config"
	"github.com/nnqq/scr-org-producer/healthz"
	"github.com/nnqq/scr-org-producer/logger"
	"github.com/nnqq/scr-org-producer/mongo"
	"github.com/nnqq/scr-org-producer/reindex"
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

	db, err := mongo.NewConn(ctx, cfg.MongoDB.URL)
	logg.Must(err)

	orgClient, err := call.NewClients(cfg.Service.Org)
	logg.Must(err)

	go func() {
		logg.Must(healthz.NewHealthz(logg.ZL, cfg.HTTP.Port).Serve())
	}()

	logg.Must(reindex.NewReindex(logg.ZL, orgClient).Do(ctx))
}
