package handlers

import (
	"context"
	"log/slog"

	dotoapi "github.com/DistributedShenanigans/doto/internal/infrastructure/clients/doto"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type DeleteHandler struct {
	ApiService dotoapi.ClientWithResponsesInterface
}

func NewDeleteHandler(apiService dotoapi.ClientWithResponsesInterface) *DeleteHandler {
	return &DeleteHandler{
		ApiService: apiService,
	}
}

func (h *DeleteHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "handlers.DeleteHandler"
	slog.Info(op, "chat_id", update.Message.Chat.ID, "user_id", update.Message.From.ID)

	// Extract task ID from the message text
	taskID := update.Message.Text

	// Call the API to delete the task
	resp, err := h.ApiService.DeleteTasksTaskIdWithResponse(ctx, taskID, &dotoapi.DeleteTasksTaskIdParams{
		TgChatId: update.Message.Chat.ID,
	})
	if err != nil || resp.StatusCode() != 200 {
		slog.Error(op, "failed to delete task", slog.Any("error", err))
		return
	}

	// Send a confirmation message
	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "_Task deleted successfully!_ Use /list to see your tasks.",
		ParseMode: models.ParseModeMarkdownV1,
	}); err != nil {
		slog.Error(op, "failed to send message", slog.Any("error", err))
	}
}
