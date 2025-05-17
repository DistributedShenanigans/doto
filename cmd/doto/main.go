package main

import (
	"context"
	"flag"
	"os"

	dotoapi "github.com/DistributedShenanigans/doto/api"
	"github.com/DistributedShenanigans/doto/config"
	"github.com/DistributedShenanigans/doto/internal/infrastructure/repository/tasks"
	dotosrv "github.com/DistributedShenanigans/doto/internal/infrastructure/servers/doto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	configFileName := flag.String("config", "./config/default-config.yaml", "path to config file")

	flag.Parse()

	cfg, err := config.New(*configFileName)
	if err != nil {
		os.Exit(1)
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.Database.ToDSN()))
	if err != nil {
		os.Exit(1)
	}

	repo := tasks.New(client.Database("doto"), "tasks")

	srv := dotosrv.New(cfg.Serving, dotoapi.New(repo))

	srv.ListenAndServe()
}
