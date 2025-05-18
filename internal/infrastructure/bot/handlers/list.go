package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"sort"

	dotoapi "github.com/DistributedShenanigans/doto/internal/infrastructure/clients/doto"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type ListHandler struct {
	ApiService dotoapi.ClientWithResponsesInterface
}

func NewListHandler(apiService dotoapi.ClientWithResponsesInterface) *ListHandler {
	return &ListHandler{
		ApiService: apiService,
	}
}

func (h *ListHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "handlers.ListHandler"
	slog.Info(op, "chat_id", update.Message.Chat.ID, "user_id", update.Message.From.ID)

	resp, err := h.ApiService.GetTasksWithResponse(ctx, &dotoapi.GetTasksParams{
		TgChatId: update.Message.Chat.ID,
	})
	if err != nil || resp.StatusCode() != 200 {
		slog.Error(op, "failed to get tasks", slog.Any("error", err))
		return
	}

	tasks := *resp.JSON200
	if len(tasks) == 0 {
		slog.Warn(op, "no tasks found", slog.Int("chat_id", int(update.Message.Chat.ID)))

		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "No tasks found. Try adding some!",
			ReplyMarkup: &models.ReplyKeyboardMarkup{
				Keyboard: [][]models.KeyboardButton{
					{
						{Text: "/add"},
					},
				},
				ResizeKeyboard:  true,
				OneTimeKeyboard: true,
			},
			ParseMode: models.ParseModeMarkdownV1,
		}); err != nil {
			slog.Error(op, "failed to send message", slog.Any("error", err))
		}

		return
	}

	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].Status == tasks[j].Status {
			return tasks[i].Id < tasks[j].Id
		}
		return tasks[i].Status == StatusPending
	})

	var messageText string

	for i, task := range tasks {
		statusEmoji := "🔴"

		if task.Status == StatusInProgress {
			statusEmoji = "🟡"
		}

		if task.Status == StatusDone {
			statusEmoji = "🟢"
		}

		messageText += fmt.Sprintf("%d. %s %s\n", i+1, task.Description, statusEmoji)
	}

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      messageText,
		ParseMode: models.ParseModeMarkdownV1,
	}); err != nil {
		slog.Error(op, "failed to send message", slog.Any("error", err))
	}
}
