package handler

import (
	"errors"
	"net/http"
	"task-api/internal/repository"
	"task-api/internal/service"
)

func handleError(w http.ResponseWriter, err error) {
	if errors.Is(err, service.ErrInvalidTitle) {
		writeError(w, http.StatusBadRequest, service.ErrInvalidTitle.Error())
		return
	} else if errors.Is(err, service.ErrInvalidID) {
		writeError(w, http.StatusBadRequest, service.ErrInvalidID.Error())
		return
	} else if errors.Is(err, repository.ErrTaskNotFound) {
		writeError(w, http.StatusNotFound, repository.ErrTaskNotFound.Error())
		return
	} else {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
}
