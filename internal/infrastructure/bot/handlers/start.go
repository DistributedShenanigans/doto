package handlers

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "handlers.StartHandler"
	slog.Info(op, "chat_id", update.Message.Chat.ID, "user_id", update.Message.From.ID)

	// Send a welcome message
	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "_Welcome to the bot!_ Use /help to see available commands.",
		ReplyMarkup: &models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{
					{
						Text: "Help",
					},
				},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: true,
		},
		ParseMode: models.ParseModeMarkdownV1,
	}); err != nil {
		slog.Error(op, "failed to send message", slog.Any("error", err))
	}
}
