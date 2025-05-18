package handlers

import (
	"context"
	"log/slog"

	dotoapi "github.com/DistributedShenanigans/doto/internal/infrastructure/clients/doto"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type AddHandler struct {
	ApiService dotoapi.ClientWithResponsesInterface
}

func NewAddHandler(apiService dotoapi.ClientWithResponsesInterface) *AddHandler {
	return &AddHandler{
		ApiService: apiService,
	}
}

func (h *AddHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "handlers.AddHandler"
	slog.Info(op, "chat_id", update.Message.Chat.ID, "user_id", update.Message.From.ID)

	resp, err := h.ApiService.PostTasksWithResponse(ctx, &dotoapi.PostTasksParams{
		TgChatId: update.Message.Chat.ID,
	}, dotoapi.TaskCreation{
		Description: update.Message.Text,
		Status:      "pending",
	})
	if err != nil || resp.StatusCode() != 201 {
		slog.Error(op, "failed to create task", slog.Any("error", err))
		return
	}

	// Send a welcome message
	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "_Task added successfully!_ Use /list to see your tasks.",
		ParseMode: models.ParseModeMarkdownV1,
	}); err != nil {
		slog.Error(op, "failed to send message", slog.Any("error", err))
	}
}
