package dotoapi

import (
	"context"
	"encoding/json"
	"net/http"
)

type TasksRepository interface {
	Get(ctx context.Context, chatId int64) ([]Task, error)
	Add(ctx context.Context, chatId int64, task TaskCreation) error
	UpdateStatus(ctx context.Context, chatId int64, taskId string, update TaskStatusUpdate) (Task, error)
	Delete(ctx context.Context, chatId int64, taskId string) error
}

type Server struct {
	tasks TasksRepository
}

func New(tasks TasksRepository) ServerInterface {
	return &Server{
		tasks: tasks,
	}
}

func (s *Server) GetTasks(w http.ResponseWriter, r *http.Request, params GetTasksParams) {
	ctx := r.Context()

	tasks, err := s.tasks.Get(ctx, params.TgChatId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func (s *Server) PostTasks(w http.ResponseWriter, r *http.Request, params PostTasksParams) {
	ctx := r.Context()

	var model TaskCreation

	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.tasks.Add(ctx, params.TgChatId, model); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, model)
}

func (s *Server) DeleteTasksTaskId(
	w http.ResponseWriter,
	r *http.Request,
	taskId string,
	params DeleteTasksTaskIdParams,
) {
	ctx := r.Context()

	if err := s.tasks.Delete(ctx, params.TgChatId, taskId); err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}

func (s *Server) PutTasksTaskId(
	w http.ResponseWriter,
	r *http.Request,
	taskID string,
	params PutTasksTaskIdParams,
) {
	ctx := r.Context()

	var model TaskStatusUpdate

	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := s.tasks.UpdateStatus(ctx, params.TgChatId, taskID, model)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}
