package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	dotoapi "github.com/DistributedShenanigans/doto/api"
	"github.com/DistributedShenanigans/doto/config"
	"github.com/DistributedShenanigans/doto/internal/infrastructure/repository/tasks"
	dotosrv "github.com/DistributedShenanigans/doto/internal/infrastructure/servers/doto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	configFileName := flag.String("config", "./config/default-config.yaml", "path to config file")

	flag.Parse()

	cfg, err := config.New(*configFileName)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.Database.ToDSN()))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	repo := tasks.New(client.Database("doto"), "tasks")

	si := dotoapi.New(repo)
	handler := dotoapi.HandlerWithOptions(si, dotoapi.StdHTTPServerOptions{
		Middlewares: []dotoapi.MiddlewareFunc{dotoapi.MetricsMiddleware},
	})

	mux := http.NewServeMux()
	mux.Handle("/", handler)
	mux.Handle("/metrics", promhttp.Handler())

	srv := dotosrv.New(cfg.Serving, mux)

	srv.ListenAndServe()
}
