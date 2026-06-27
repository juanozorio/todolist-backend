package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/juanozorio/task-api/internal/domain"
	"github.com/juanozorio/task-api/internal/service"
)

type taskHandler struct {
	svc service.TaskService
}

func newTaskHandler(svc service.TaskService) *taskHandler {
	return &taskHandler{svc: svc}
}

// create godoc
//
//	@Summary      Create a task
//	@Description  Creates a new task with the given description
//	@Tags         tasks
//	@Accept       json
//	@Produce      json
//	@Param        body  body      domain.CreateTaskInput  true  "Task payload"
//	@Success      201   {object}  domain.Task
//	@Failure      400   {object}  errorResponse
//	@Failure      500   {object}  errorResponse
//	@Router       /tasks [post]
func (h *taskHandler) create(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	task, err := h.svc.Create(r.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidPayload) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respond(w, http.StatusCreated, task)
}

// update godoc
//
//	@Summary      Update a task
//	@Description  Updates description and/or completion status of a task
//	@Tags         tasks
//	@Accept       json
//	@Produce      json
//	@Param        id    path      string                  true  "Task ID (UUID)"
//	@Param        body  body      domain.UpdateTaskInput  true  "Update payload"
//	@Success      200   {object}  domain.Task
//	@Failure      400   {object}  errorResponse
//	@Failure      404   {object}  errorResponse
//	@Failure      500   {object}  errorResponse
//	@Router       /tasks/{id} [put]
func (h *taskHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var input domain.UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	task, err := h.svc.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			respondError(w, http.StatusNotFound, "task not found")
			return
		}
		if errors.Is(err, domain.ErrInvalidPayload) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respond(w, http.StatusOK, task)
}

// getAll godoc
//
//	@Summary      List all tasks
//	@Description  Returns all tasks ordered by creation date (newest first)
//	@Tags         tasks
//	@Produce      json
//	@Success      200  {array}   domain.Task
//	@Failure      500  {object}  errorResponse
//	@Router       /tasks [get]
func (h *taskHandler) getAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.svc.GetAll(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	respond(w, http.StatusOK, tasks)
}

// getByStatus godoc
//
//	@Summary      List tasks by status
//	@Description  Returns tasks filtered by their completion status
//	@Tags         tasks
//	@Produce      json
//	@Param        completed  query     bool  true  "Filter by completion status (true or false)"
//	@Success      200        {array}   domain.Task
//	@Failure      400        {object}  errorResponse
//	@Failure      500        {object}  errorResponse
//	@Router       /tasks/status [get]
func (h *taskHandler) getByStatus(w http.ResponseWriter, r *http.Request) {
	completedStr := r.URL.Query().Get("completed")
	if completedStr == "" {
		respondError(w, http.StatusBadRequest, "query param 'completed' is required (true or false)")
		return
	}

	completed, err := strconv.ParseBool(completedStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "query param 'completed' must be true or false")
		return
	}

	tasks, err := h.svc.GetByStatus(r.Context(), completed)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	respond(w, http.StatusOK, tasks)
}

// getByID godoc
//
//	@Summary      Get a task by ID
//	@Description  Returns a single task by its UUID
//	@Tags         tasks
//	@Produce      json
//	@Param        id   path      string  true  "Task ID (UUID)"
//	@Success      200  {object}  domain.Task
//	@Failure      404  {object}  errorResponse
//	@Failure      500  {object}  errorResponse
//	@Router       /tasks/{id} [get]
func (h *taskHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			respondError(w, http.StatusNotFound, "task not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respond(w, http.StatusOK, task)
}
