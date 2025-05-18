package errors

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	ErrorMsgInternal = "Internal error. Please try again later."
	ErrorMsgNoTasks  = "No tasks found. Try `/add <task_description>` to add a task."
)

func HandleError(ctx context.Context, b *bot.Bot, chatID int64, op string, err error, msg string) {
	slog.Error(op, "error occurred",
		slog.Any("error", err),
		slog.Int64("chat_id", chatID),
	)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      msg,
		ParseMode: models.ParseModeMarkdownV1,
	}); err != nil {
		slog.Error(op, "failed to send error message", slog.Any("error", err))
	}
}
