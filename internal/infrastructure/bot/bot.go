package bot

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"

	"github.com/DistributedShenanigans/doto/config"
	"github.com/DistributedShenanigans/doto/internal/infrastructure/bot/handlers"
	dotoapi "github.com/DistributedShenanigans/doto/internal/infrastructure/clients/doto"
	"github.com/go-telegram/bot"
)

type BotService struct {
	ApiBot     *bot.Bot
	ApiService dotoapi.ClientWithResponsesInterface
}

func NewBotService(config *config.Config, apiService dotoapi.ClientWithResponsesInterface) (*BotService, error) {
	const op = "bot.NewBotService"

	b, err := bot.New(
		config.BotToken,
		bot.WithWorkers(runtime.NumCPU()),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	addHandler := handlers.NewAddHandler(apiService)
	updateHandler := handlers.NewUpdateHandler(apiService)
	deleteHandler := handlers.NewDeleteHandler(apiService)

	b.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommand, handlers.StartHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "add", bot.MatchTypeCommand, addHandler.Handle)
	b.RegisterHandler(bot.HandlerTypeMessageText, "update", bot.MatchTypeCommand, updateHandler.Handle)
	b.RegisterHandler(bot.HandlerTypeMessageText, "delete", bot.MatchTypeCommand, deleteHandler.Handle)

	return &BotService{
		ApiBot:     b,
		ApiService: apiService,
	}, nil
}

func (b *BotService) Start(ctx context.Context) {
	slog.Info("starting bot service")

	b.ApiBot.Start(ctx)
}
