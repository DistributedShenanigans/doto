package handlers

import (
	"context"
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/DistributedShenanigans/doto/internal/infrastructure/clients/doto"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type UpdateHandler struct {
	ApiService dotoapi.ClientWithResponsesInterface
}

func NewUpdateHandler(apiService dotoapi.ClientWithResponsesInterface) *UpdateHandler {
	return &UpdateHandler{
		ApiService: apiService,
	}
}

func (h *UpdateHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "handlers.UpdateHandler"
	slog.Info(op, "chat_id", update.Message.Chat.ID, "user_id", update.Message.From.ID)

	// Extract task ID and new description from the message text
	parts := strings.SplitN(update.Message.Text, " ", 2)
	if len(parts) != 2 {
		slog.Error(op, "invalid message format", slog.String("message", update.Message.Text))
		return
	}
	taskID := parts[0]
	newDescription := parts[1]

	// Validate the new description length
	if utf8.RuneCountInString(newDescription) > 100 {
		slog.Error(op, "description too long", slog.Int("length", utf8.RuneCountInString(newDescription)))
		return
	}

	// Call the API to update the task
	resp, err := h.ApiService.PutTasksTaskIdWithResponse(ctx, taskID, &dotoapi.PutTasksTaskIdParams{
		TgChatId: update.Message.Chat.ID,
	}, dotoapi.TaskStatusUpdate{
		Status: "done",
	})
	if err != nil || resp.StatusCode() != 200 {
		slog.Error(op, "failed to update task", slog.Any("error", err))
		return
	}

	// Send a confirmation message
	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "_Task updated successfully!_ Use /list to see your tasks.",
		ParseMode: models.ParseModeMarkdownV1,
	}); err != nil {
		slog.Error(op, "failed to send message", slog.Any("error", err))
	}
}
