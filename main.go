package main

import (
	"context"
	"github.com/chn555/property-service/internal/rest/property"
	"github.com/chn555/property-service/pkg/db/mongo"
	property2 "github.com/chn555/property-service/pkg/property"
	"log/slog"

	"github.com/chn555/property-service/internal/config"
	"github.com/chn555/property-service/internal/rest"
	"os"
)

func main() {

	cfg, err := config.LoadConfig(context.Background())
	if err != nil {
		slog.Error("failed to load config", slog.String("err", err.Error()))
		os.Exit(1)
	}

	mongoClient, err := mongo.NewClient(context.TODO(), cfg.MongoConfig)
	if err != nil {
		slog.Error("failed to connect to mongo db", slog.String("err", err.Error()))
		os.Exit(1)
	}

	con := mongo.NewEventState(mongoClient, cfg.MongoEventStateConfig)

	e := rest.NewServer(
		property.NewRestHandler(property2.NewHandler(con)).RegisterHandlers,
	)

	if err := e.Start(":1323"); err != nil {
	}
	e.Logger.Fatal(e.Start(":1323"))
}
