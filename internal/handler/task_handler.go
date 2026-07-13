package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
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

func (t *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		list := t.serv.List()

		writeJSON(w, http.StatusOK, list)
	} else if r.Method == http.MethodPost {
		var req createTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}

		task, err := t.serv.Create(req.Title)
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

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	if r.Method == http.MethodGet {
		task, err := t.serv.GetByID(id)
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

		task, err := t.serv.Update(req.Title, id, req.Done)
		if err != nil {
			handleError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, task)
	} else if r.Method == http.MethodDelete {
		err := t.serv.Delete(id)
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
