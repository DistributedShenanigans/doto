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
	"github.com/go-telegram/bot/models"
)

type BotService struct {
	ApiBot     *bot.Bot
	ApiService dotoapi.ClientWithResponsesInterface
}

func NewBotService(
	config *config.Config,
	apiService dotoapi.ClientWithResponsesInterface,
) (*BotService, error) {
	const op = "bot.NewBotService"

	opts := []bot.Option{
		bot.WithWorkers(runtime.NumCPU()),
		bot.WithDefaultHandler(handlers.DefaultHandler),
	}

	b, err := bot.New(
		config.BotToken,
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	addHandler := handlers.NewAddHandler(apiService)
	updateHandler := handlers.NewUpdateHandler(apiService)
	deleteHandler := handlers.NewDeleteHandler(apiService)
	listHandler := handlers.NewListHandler(apiService)

	// Commands
	b.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommand, handlers.StartHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "help", bot.MatchTypeCommand, handlers.HelpHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "add", bot.MatchTypeCommand, addHandler.Handle)
	b.RegisterHandler(bot.HandlerTypeMessageText, "list", bot.MatchTypeCommand, listHandler.Handle)
	b.RegisterHandler(bot.HandlerTypeMessageText, "update", bot.MatchTypeCommand, updateHandler.Handle)
	b.RegisterHandler(bot.HandlerTypeMessageText, "delete", bot.MatchTypeCommand, deleteHandler.Handle)

	// Callbacks
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "update", bot.MatchTypePrefix, updateHandler.HandleCallback)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "delete", bot.MatchTypePrefix, deleteHandler.HandleCallback)

	_, err = b.SetMyCommands(context.Background(), &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{
				Command:     "start",
				Description: "Start the bot",
			},
			{
				Command:     "add",
				Description: "Add a new task, e.g. `/add <task_description>`",
			},
			{
				Command:     "list",
				Description: "List all tasks",
			},
			{
				Command:     "update",
				Description: "Update an existing task",
			},
			{
				Command:     "delete",
				Description: "Delete an existing task",
			},
			{
				Command:     "help",
				Description: "Show help",
			},
		},
	})
	if err != nil {
		slog.Warn(op, "failed to set commands", slog.Any("error", err))
	}

	return &BotService{
		ApiBot:     b,
		ApiService: apiService,
	}, nil
}

func (b *BotService) Start(ctx context.Context) {
	slog.Info("starting bot service")

	b.ApiBot.Start(ctx)
}
