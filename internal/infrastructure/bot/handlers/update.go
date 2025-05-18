package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"sort"

	"github.com/DistributedShenanigans/doto/internal/infrastructure/bot/errors"
	dotoapi "github.com/DistributedShenanigans/doto/internal/infrastructure/clients/doto"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
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

	resp, err := h.ApiService.GetTasksWithResponse(ctx, &dotoapi.GetTasksParams{
		TgChatId: update.Message.Chat.ID,
	})
	if err != nil || resp.StatusCode() != 200 {
		errors.HandleError(ctx, b, update.Message.Chat.ID, op, err, errors.ErrorMsgInternal)
		return
	}

	tasks := *resp.JSON200
	if len(tasks) == 0 {
		errors.HandleError(ctx, b, update.Message.Chat.ID, op, nil, errors.ErrorMsgNoTasks)

		return
	}

	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].Status == tasks[j].Status {
			return tasks[i].Id < tasks[j].Id
		}
		return tasks[i].Status == StatusPending
	})

	var messageText string
	var keyboardButtons [][]models.InlineKeyboardButton

	for i, task := range tasks {
		statusEmoji := "🔴"

		if task.Status == StatusInProgress {
			statusEmoji = "🟡"
		}

		if task.Status == StatusDone {
			statusEmoji = "🟢"
		}

		messageText += fmt.Sprintf("%d. %s %s\n", i+1, task.Description, statusEmoji)

		row := []models.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("%d", i+1),
				CallbackData: fmt.Sprintf("update_%s", task.Id),
			},
		}
		keyboardButtons = append(keyboardButtons, row)
	}

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        messageText,
		ParseMode:   models.ParseModeMarkdownV1,
		ReplyMarkup: &models.InlineKeyboardMarkup{InlineKeyboard: keyboardButtons},
	}); err != nil {
		slog.Error(op, "failed to send message", slog.Any("error", err))
	}
}

func (h *UpdateHandler) HandleCallback(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "handlers.UpdateHandler.HandleCallback"

	slog.Debug("handling callback", "op", op, "chat_id", update.CallbackQuery.From.ID, "data", update.CallbackQuery.Data)

	taskID := update.CallbackQuery.Data[7:]

	slog.Debug("retrieving task ID", "task_id", taskID)

	taskResp, err := h.ApiService.GetTasksWithResponse(ctx, &dotoapi.GetTasksParams{
		TgChatId: update.CallbackQuery.From.ID,
	})
	if err != nil || taskResp.StatusCode() != 200 {
		errors.HandleError(ctx, b, update.CallbackQuery.From.ID, op, err, errors.ErrorMsgInternal)
		return
	}

	slog.Debug(op, "task response", slog.Any("response", taskResp.JSON200))

	tasks := *taskResp.JSON200
	var currentTask *dotoapi.Task
	var newStatus string

	for _, task := range tasks {
		if task.Id == taskID {
			currentTask = &task
			break
		}
	}

	if currentTask == nil {
		slog.Error(op, "task not found", slog.String("task_id", taskID))
		errors.HandleError(ctx, b, update.CallbackQuery.From.ID, op, nil, errors.ErrorMsgInternal)
		return
	}

	slog.Debug(op, "current task", slog.Any("task", currentTask))

	switch currentTask.Status {
	case StatusPending:
		newStatus = StatusInProgress
	case StatusInProgress:
		newStatus = StatusDone
	case StatusDone:
		newStatus = StatusPending
	default:
		slog.Error(op, "unknown task status", slog.String("status", currentTask.Status))
		return
	}

	currentTask.Status = newStatus

	resp, err := h.ApiService.PutTasksTaskIdWithResponse(ctx, taskID, &dotoapi.PutTasksTaskIdParams{
		TgChatId: update.CallbackQuery.From.ID,
	}, dotoapi.TaskStatusUpdate{
		Status: newStatus,
	})
	if err != nil || resp.StatusCode() != 200 {
		slog.Error(op, "failed to update task", slog.Any("error", err))
		return
	}

	slog.Debug(op, "task updated successfully", slog.Any("response", resp.JSON200))

	messageText := fmt.Sprintf("*Task updated successfully!*\n\n%s", currentTask.Description)
	switch newStatus {
	case StatusPending:
		messageText += "\n\n*New Status:* Pending 🔴"
	case StatusInProgress:
		messageText += "\n\n*New Status:* In Progress 🟡"
	default:
		messageText += "\n\n*New Status:* Done 🟢"
	}

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    update.CallbackQuery.From.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
		Text:      messageText,
		ParseMode: models.ParseModeMarkdownV1,
	}); err != nil {
		slog.Error(op, "failed to edit message", slog.Any("error", err))

		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.CallbackQuery.From.ID,
			Text:      messageText,
			ParseMode: models.ParseModeMarkdownV1,
		})
	}
}
