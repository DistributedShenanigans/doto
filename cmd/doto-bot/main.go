package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/DistributedShenanigans/doto/config"
	"github.com/DistributedShenanigans/doto/internal/infrastructure/bot"
	dotoapi "github.com/DistributedShenanigans/doto/internal/infrastructure/clients/doto"
)

func main() {
	configFileName := flag.String("config", "./config/default-config.yaml", "path to config file")

	flag.Parse()

	cfg, err := config.New(*configFileName)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	slog.Info("config loaded", "config", *cfg)

	apiClient, err := dotoapi.NewClientWithResponses(
		fmt.Sprintf("http://%s:%d", cfg.Serving.Host, cfg.Serving.BotPort),
	)
	if err != nil {
		slog.Error("failed to create api client", "error", err)
		os.Exit(1)
	}

	slog.Info("api client created", "host", cfg.Serving.Host, "port", cfg.Serving.BotPort)

	bot, err := bot.NewBotService(cfg, apiClient)
	if err != nil {
		slog.Error("failed to create bot service", "error", err)
		os.Exit(1)
	}

	slog.Info("bot service created")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	bot.Start(ctx)
}
