package handlers

import (
	"context"
	"strings"

	"github.com/DistributedShenanigans/doto/internal/infrastructure/bot/errors"
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
	const op = "handlers.AddHandler.Handle"

	text := strings.TrimSpace(update.Message.Text)
	if len(text) <= 5 {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "Usage: `/add <task_description>`\n\nExample: `/add Buy groceries`",
			ParseMode: models.ParseModeMarkdownV1,
		}); err != nil {
			errors.HandleError(ctx, b, update.Message.Chat.ID, op, err, errors.ErrorMsgInternal)
		}
		return
	}

	description := strings.TrimSpace(text[5:])

	resp, err := h.ApiService.PostTasksWithResponse(ctx, &dotoapi.PostTasksParams{
		TgChatId: update.Message.Chat.ID,
	}, dotoapi.TaskCreation{
		Description: description,
		Status:      StatusPending,
	})

	if err != nil || resp.StatusCode() != 201 {
		errors.HandleError(ctx, b, update.Message.Chat.ID, op, err, errors.ErrorMsgInternal)
		return
	}

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "_Task added successfully!_\n\nYou can view your tasks with `/list`.",
		ReplyMarkup: &models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{
					{Text: "/list"},
				},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: true,
		},
		ParseMode: models.ParseModeMarkdownV1,
	}); err != nil {
		errors.HandleError(ctx, b, update.Message.Chat.ID, op, err, errors.ErrorMsgInternal)
	}
}
