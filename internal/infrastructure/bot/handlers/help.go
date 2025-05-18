package handlers

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func HelpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "handlers.HelpHandler"
	slog.Info(op, "chat_id", update.Message.Chat.ID, "user_id", update.Message.From.ID)

	message := ""

	commands, err := b.GetMyCommands(ctx, nil)
	if err != nil {
		slog.Error(op, "failed to get commands", slog.Any("error", err))
	}

	if len(commands) > 0 {
		message = "_Here are the available commands:_\n\n"
		for _, command := range commands {
			message += "/" + command.Command + " - " + command.Description + "\n"
		}
	} else {
		message = "_No commands available._"
	}

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      message,
		ParseMode: models.ParseModeMarkdownV1,
	}); err != nil {
		slog.Error(op, "failed to send message", slog.Any("error", err))
	}
}
