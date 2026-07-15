package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"task-api/internal/model"
	"task-api/internal/service"
)

type TaskHandler struct {
	serv *service.TaskService
}

func NewTaskHandler(serv *service.TaskService) *TaskHandler {
	taskHandler := TaskHandler{
		serv: serv,
	}

	return &taskHandler
}

func (t *TaskHandler) HandlePing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

type createTaskRequest struct {
	Title string `json:"title"`
}

func parseTaskFilter(r *http.Request) (model.TaskFilter, error) {
	filter := model.TaskFilter{
		Limit:  20,
		Offset: 0,
	}

	query := r.URL.Query()

	doneStr := query.Get("done")
	if doneStr != "" {
		done, errDone := strconv.ParseBool(doneStr)
		if errDone != nil {
			return model.TaskFilter{}, errors.New("invalid done filter")
		}

		filter.Done = &done
	}

	limitStr := query.Get("limit")
	if limitStr != "" {
		limit, errLimit := strconv.Atoi(limitStr)
		if errLimit != nil {
			return model.TaskFilter{}, errors.New("invalid limit filter")
		}

		if limit <= 0 {
			return model.TaskFilter{}, errors.New("limit must be greater than 0")
		}

		if limit > 100 {
			return model.TaskFilter{}, errors.New("limit must be less than or equal to 100")
		}

		filter.Limit = limit
	}

	offsetStr := query.Get("offset")
	if offsetStr != "" {
		offset, errOffset := strconv.Atoi(offsetStr)
		if errOffset != nil {
			return model.TaskFilter{}, errors.New("invalid offset filter")
		}

		if offset < 0 {
			return model.TaskFilter{}, errors.New("offset must be greater than or equal to 0")
		}

		filter.Offset = offset
	}

	return filter, nil
}

func (t *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method == http.MethodGet {
		filter, errFilter := parseTaskFilter(r)
		if errFilter != nil {
			writeError(w, http.StatusBadRequest, errFilter.Error())
			return
		}

		list, errList := t.serv.List(ctx, filter)
		if errList != nil {
			handleError(w, errList)
			return
		}

		writeJSON(w, http.StatusOK, list)
	} else if r.Method == http.MethodPost {
		var req createTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}

		task, err := t.serv.Create(ctx, req.Title)
		if err != nil {
			handleError(w, err)
			return
		}

		writeJSON(w, http.StatusCreated, task)

	} else {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
}

type updateTaskRequest struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func (t *TaskHandler) HandleTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	ctx := r.Context()

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	if r.Method == http.MethodGet {
		task, err := t.serv.GetByID(ctx, id)
		if err != nil {
			handleError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, task)

	} else if r.Method == http.MethodPatch {
		var req updateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}

		task, err := t.serv.Update(r.Context(), id, req.Title, req.Done)
		if err != nil {
			handleError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, task)
	} else if r.Method == http.MethodDelete {
		err := t.serv.Delete(ctx, id)
		if err != nil {
			handleError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	} else {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
}
