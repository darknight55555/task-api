package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"task-api/internal/model"
	"task-api/internal/repository"
	"task-api/internal/service"
)

func newTestHandler() *TaskHandler {
	repo := repository.NewMemoryTaskRepository()
	serv := service.NewTaskService(repo)
	return NewTaskHandler(serv)
}

func TestHandlePingGET(t *testing.T) {
	h := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rr := httptest.NewRecorder()

	h.HandlePing(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected body %d, got %d", http.StatusOK, rr.Code)
	}

	if rr.Body.String() != "pong" {
		t.Fatalf("expected body %q, got %q", "pong", rr.Body.String())
	}
}

func TestHandlePingWrongMethod(t *testing.T) {
	h := newTestHandler()

	req := httptest.NewRequest(http.MethodPost, "/ping", nil)
	rr := httptest.NewRecorder()

	h.HandlePing(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestHandleTasksGET(t *testing.T) {
	h := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rr := httptest.NewRecorder()

	h.HandleTasks(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestHandleTasksPost(t *testing.T) {
	h := newTestHandler()

	body := strings.NewReader(`{"title":"learn go"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	rr := httptest.NewRecorder()

	h.HandleTasks(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
}

func TestHandleTasksPOSTEmptyTitle(t *testing.T) {
	h := newTestHandler()

	body := strings.NewReader(`{"title":""}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	rr := httptest.NewRecorder()

	h.HandleTasks(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleTasksByIDGETInvalidID(t *testing.T) {
	h := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/tasks/abc", nil)
	rr := httptest.NewRecorder()

	h.HandleTaskByID(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleTasksByIDGETNotFound(t *testing.T) {
	h := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/tasks/999", nil)
	rr := httptest.NewRecorder()

	h.HandleTaskByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestHandleTasksByIDGETSuccess(t *testing.T) {
	h := newTestHandler()

	body := strings.NewReader(`{"title":"learn go"}`)
	reqPost := httptest.NewRequest(http.MethodPost, "/tasks", body)
	rrPost := httptest.NewRecorder()

	h.HandleTasks(rrPost, reqPost)

	if rrPost.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rrPost.Code)
	}

	reqGet := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	rrGet := httptest.NewRecorder()

	h.HandleTaskByID(rrGet, reqGet)

	if rrGet.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rrGet.Code)
	}
}

func TestHandleTasksByIDPatchSuccess(t *testing.T) {
	h := newTestHandler()

	bodyPost := strings.NewReader(`{"title":"learn go"}`)
	reqPost := httptest.NewRequest(http.MethodPost, "/tasks", bodyPost)
	rrPost := httptest.NewRecorder()

	h.HandleTasks(rrPost, reqPost)

	if rrPost.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rrPost.Code)
	}

	bodyPatch := strings.NewReader(`{"title":"learn hard","done":true}`)
	reqPatch := httptest.NewRequest(http.MethodPatch, "/tasks/1", bodyPatch)
	rrPatch := httptest.NewRecorder()

	h.HandleTaskByID(rrPatch, reqPatch)

	if rrPatch.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rrPatch.Code)
	}
}

func TestHandleTasksByIDPatchBadRequestInvalidTitle(t *testing.T) {
	h := newTestHandler()

	bodyPost := strings.NewReader(`{"title":"learn go"}`)
	reqPost := httptest.NewRequest(http.MethodPost, "/tasks", bodyPost)
	rrPost := httptest.NewRecorder()

	h.HandleTasks(rrPost, reqPost)

	if rrPost.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rrPost.Code)
	}

	bodyPatch := strings.NewReader(`{"title":"","done":true}`)
	reqPatch := httptest.NewRequest(http.MethodPatch, "/tasks/1", bodyPatch)
	rrPatch := httptest.NewRecorder()

	h.HandleTaskByID(rrPatch, reqPatch)

	if rrPatch.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rrPatch.Code)
	}
}

func TestHandleTasksByIDDeleteSuccess(t *testing.T) {
	h := newTestHandler()

	bodyPost := strings.NewReader(`{"title":"learn go"}`)
	reqPost := httptest.NewRequest(http.MethodPost, "/tasks", bodyPost)
	rrPost := httptest.NewRecorder()

	h.HandleTasks(rrPost, reqPost)

	if rrPost.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rrPost.Code)
	}

	reqDel := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	rrDel := httptest.NewRecorder()

	h.HandleTaskByID(rrDel, reqDel)

	if rrDel.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rrDel.Code)
	}
}

func TestHandleTasksPostReturnsCreatedTask(t *testing.T) {
	h := newTestHandler()

	bodyPost := strings.NewReader(`{"title":"learn go"}`)
	reqPost := httptest.NewRequest(http.MethodPost, "/tasks", bodyPost)
	rrPost := httptest.NewRecorder()

	h.HandleTasks(rrPost, reqPost)

	if rrPost.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rrPost.Code)
	}

	var task model.Task

	if err := json.NewDecoder(rrPost.Body).Decode(&task); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if task.ID != 1 {
		t.Fatalf("expected id %d, got %d", 1, task.ID)
	}

	if task.Title != "learn go" {
		t.Fatalf("expected title %q, got %q", "learn go", task.Title)
	}

	if task.Done != false {
		t.Fatalf("expected done %t, got %t", false, task.Done)
	}
}

func TestHandleTasksGETInvalidFilter(t *testing.T) {
	h := newTestHandler()

	reqGet := httptest.NewRequest(http.MethodGet, "/tasks?done=abc", nil)
	rrGet := httptest.NewRecorder()

	h.HandleTasks(rrGet, reqGet)

	if rrGet.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rrGet.Code)
	}
}

func TestHandleTasksGETDoneTrue(t *testing.T) {
	h := newTestHandler()

	bodyPost1 := strings.NewReader(`{"title":"exmpl1"}`)
	reqPost1 := httptest.NewRequest(http.MethodPost, "/tasks", bodyPost1)
	rrPost1 := httptest.NewRecorder()

	h.HandleTasks(rrPost1, reqPost1)

	if rrPost1.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rrPost1.Code)
	}

	bodyPost2 := strings.NewReader(`{"title":"exmpl2"}`)
	reqPost2 := httptest.NewRequest(http.MethodPost, "/tasks", bodyPost2)
	rrPost2 := httptest.NewRecorder()

	h.HandleTasks(rrPost2, reqPost2)

	if rrPost2.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rrPost2.Code)
	}

	bodyPatch := strings.NewReader(`{"title":"exmpl1 did it", "done":true}`)
	reqPatch := httptest.NewRequest(http.MethodPatch, "/tasks/1", bodyPatch)
	rrPatch := httptest.NewRecorder()

	h.HandleTaskByID(rrPatch, reqPatch)

	if rrPatch.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rrPatch.Code)
	}

	reqGet := httptest.NewRequest(http.MethodGet, "/tasks?done=true", nil)
	rrGet := httptest.NewRecorder()

	h.HandleTasks(rrGet, reqGet)

	if rrGet.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rrGet.Code)
	}

	var tasks []model.Task
	if err := json.NewDecoder(rrGet.Body).Decode(&tasks); err != nil {
		t.Fatal("failed to decode", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("expected %d task, got %d", 1, len(tasks))
	}

	if !tasks[0].Done {
		t.Fatalf("FILTER failed: expected filtered task done to be true, got: %+v, all tasks: %+v", tasks[0], tasks)
	}

}

func TestHandleTasksGETWithLimitOffset(t *testing.T) {
	h := newTestHandler()

	bodyPost1 := strings.NewReader(`{"title":"exmpl1"}`)
	reqPost1 := httptest.NewRequest(http.MethodPost, "/tasks", bodyPost1)
	rrPost1 := httptest.NewRecorder()

	h.HandleTasks(rrPost1, reqPost1)

	if rrPost1.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rrPost1.Code)
	}

	bodyPost2 := strings.NewReader(`{"title":"exmpl2"}`)
	reqPost2 := httptest.NewRequest(http.MethodPost, "/tasks", bodyPost2)
	rrPost2 := httptest.NewRecorder()

	h.HandleTasks(rrPost2, reqPost2)

	if rrPost2.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rrPost2.Code)
	}

	bodyPost3 := strings.NewReader(`{"title":"exmpl3"}`)
	reqPost3 := httptest.NewRequest(http.MethodPost, "/tasks", bodyPost3)
	rrPost3 := httptest.NewRecorder()

	h.HandleTasks(rrPost3, reqPost3)

	if rrPost3.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rrPost3.Code)
	}

	reqGet := httptest.NewRequest(http.MethodGet, "/tasks?limit=2&offset=1", nil)
	rrGet := httptest.NewRecorder()

	h.HandleTasks(rrGet, reqGet)

	if rrGet.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rrGet.Code)
	}

	var tasks []model.Task
	if err := json.NewDecoder(rrGet.Body).Decode(&tasks); err != nil {
		t.Fatal("failed to decode", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("expected %d tasks, got %d", 2, len(tasks))
	}

	if tasks[0].ID != 2 || tasks[1].ID != 3 {
		t.Fatalf("expected tasks with id 2 and 3, got %d and %d", tasks[0].ID, tasks[1].ID)
	}

}

func TestHandleTasksGETInvalidLimit(t *testing.T) {
	h := newTestHandler()

	reqGet := httptest.NewRequest(http.MethodGet, "/tasks?limit=abc", nil)
	rrGet := httptest.NewRecorder()

	h.HandleTasks(rrGet, reqGet)

	if rrGet.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rrGet.Code)
	}
}

func TestHandleTasksGETInvalidOffset(t *testing.T) {
	h := newTestHandler()

	reqGet := httptest.NewRequest(http.MethodGet, "/tasks?offset=-1", nil)
	rrGet := httptest.NewRecorder()

	h.HandleTasks(rrGet, reqGet)

	if rrGet.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rrGet.Code)
	}
}

func TestHandleTasksGETLimitTooLarge(t *testing.T) {
	h := newTestHandler()

	reqGet := httptest.NewRequest(http.MethodGet, "/tasks?limit=101", nil)
	rrGet := httptest.NewRecorder()

	h.HandleTasks(rrGet, reqGet)

	if rrGet.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rrGet.Code)
	}
}

func TestHandleTasksGETLimitZero(t *testing.T) {
	h := newTestHandler()

	reqGet := httptest.NewRequest(http.MethodGet, "/tasks?limit=0", nil)
	rrGet := httptest.NewRecorder()

	h.HandleTasks(rrGet, reqGet)

	if rrGet.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rrGet.Code)
	}
}

func TestHandleTasksGETDoneFalse(t *testing.T) {
	h := newTestHandler()

	bodyPost1 := strings.NewReader(`{"title":"exmpl1"}`)
	reqPost1 := httptest.NewRequest(http.MethodPost, "/tasks", bodyPost1)
	rrPost1 := httptest.NewRecorder()

	h.HandleTasks(rrPost1, reqPost1)

	if rrPost1.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rrPost1.Code)
	}

	bodyPost2 := strings.NewReader(`{"title":"exmpl2"}`)
	reqPost2 := httptest.NewRequest(http.MethodPost, "/tasks", bodyPost2)
	rrPost2 := httptest.NewRecorder()

	h.HandleTasks(rrPost2, reqPost2)

	if rrPost2.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rrPost2.Code)
	}

	bodyPost3 := strings.NewReader(`{"title":"exmpl3"}`)
	reqPost3 := httptest.NewRequest(http.MethodPost, "/tasks", bodyPost3)
	rrPost3 := httptest.NewRecorder()

	h.HandleTasks(rrPost3, reqPost3)

	if rrPost3.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rrPost3.Code)
	}

	bodyPatch := strings.NewReader(`{"title":"exmple1 done", "done":true}`)
	reqPatch := httptest.NewRequest(http.MethodPatch, "/tasks/1", bodyPatch)
	rrPatch := httptest.NewRecorder()

	h.HandleTaskByID(rrPatch, reqPatch)

	if rrPatch.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rrPatch.Code)
	}

	reqGet := httptest.NewRequest(http.MethodGet, "/tasks?done=false", nil)
	rrGet := httptest.NewRecorder()

	h.HandleTasks(rrGet, reqGet)

	if rrGet.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rrGet.Code)
	}

	var tasks []model.Task
	if err := json.NewDecoder(rrGet.Body).Decode(&tasks); err != nil {
		t.Fatal("failed to decode", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}

	if tasks[0].Done != false || tasks[1].Done != false {
		t.Fatalf("expected all tasks are false, got tasks %t and %t", tasks[0].Done, tasks[1].Done)
	}
}

func TestHandleTasksGETDoneTrueWithLimitOffset(t *testing.T) {
	h := newTestHandler()

	bodyPost1 := strings.NewReader(`{"title":"exmpl1"}`)
	reqPost1 := httptest.NewRequest(http.MethodPost, "/tasks", bodyPost1)
	rrPost1 := httptest.NewRecorder()

	h.HandleTasks(rrPost1, reqPost1)

	if rrPost1.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rrPost1.Code)
	}

	bodyPost2 := strings.NewReader(`{"title":"exmpl2"}`)
	reqPost2 := httptest.NewRequest(http.MethodPost, "/tasks", bodyPost2)
	rrPost2 := httptest.NewRecorder()

	h.HandleTasks(rrPost2, reqPost2)

	if rrPost2.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rrPost2.Code)
	}

	bodyPost3 := strings.NewReader(`{"title":"exmpl3"}`)
	reqPost3 := httptest.NewRequest(http.MethodPost, "/tasks", bodyPost3)
	rrPost3 := httptest.NewRecorder()

	h.HandleTasks(rrPost3, reqPost3)

	if rrPost3.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rrPost3.Code)
	}

	bodyPost4 := strings.NewReader(`{"title":"exmpl4"}`)
	reqPost4 := httptest.NewRequest(http.MethodPost, "/tasks", bodyPost4)
	rrPost4 := httptest.NewRecorder()

	h.HandleTasks(rrPost4, reqPost4)

	if rrPost4.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rrPost4.Code)
	}

	bodyPatch1 := strings.NewReader(`{"title":"exmple1", "done": true}`)
	reqPatch1 := httptest.NewRequest(http.MethodPatch, "/tasks/1", bodyPatch1)
	rrPatch1 := httptest.NewRecorder()

	h.HandleTaskByID(rrPatch1, reqPatch1)

	if rrPatch1.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rrPatch1.Code)
	}

	bodyPatch2 := strings.NewReader(`{"title":"exmple2", "done": true}`)
	reqPatch2 := httptest.NewRequest(http.MethodPatch, "/tasks/2", bodyPatch2)
	rrPatch2 := httptest.NewRecorder()

	h.HandleTaskByID(rrPatch2, reqPatch2)

	if rrPatch2.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rrPatch2.Code)
	}

	bodyPatch3 := strings.NewReader(`{"title":"exmple3", "done": true}`)
	reqPatch3 := httptest.NewRequest(http.MethodPatch, "/tasks/3", bodyPatch3)
	rrPatch3 := httptest.NewRecorder()

	h.HandleTaskByID(rrPatch3, reqPatch3)

	if rrPatch3.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rrPatch3.Code)
	}

	reqGet := httptest.NewRequest(http.MethodGet, "/tasks?done=true&limit=2&offset=1", nil)
	rrGet := httptest.NewRecorder()

	h.HandleTasks(rrGet, reqGet)

	if rrGet.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rrGet.Code)
	}

	var tasks []model.Task
	if err := json.NewDecoder(rrGet.Body).Decode(&tasks); err != nil {
		t.Fatal("failed to decode", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}

	if tasks[0].Done != true || tasks[1].Done != true {
		t.Fatalf("expected 2 tasks with true, got %t, %t", tasks[0].Done, tasks[1].Done)
	}

	if tasks[0].ID != 2 || tasks[1].ID != 3 {
		t.Fatalf("expected tasks with id 2 and 3, got %d and %d", tasks[0].ID, tasks[1].ID)
	}

}

func TestHandleTasksGETDefaultLimit(t *testing.T) {
	h := newTestHandler()

	bodyPost1 := strings.NewReader(`{"title":"exmpl1"}`)
	reqPost1 := httptest.NewRequest(http.MethodPost, "/tasks", bodyPost1)
	rrPost1 := httptest.NewRecorder()

	h.HandleTasks(rrPost1, reqPost1)

	if rrPost1.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rrPost1.Code)
	}

	bodyPost2 := strings.NewReader(`{"title":"exmpl2"}`)
	reqPost2 := httptest.NewRequest(http.MethodPost, "/tasks", bodyPost2)
	rrPost2 := httptest.NewRecorder()

	h.HandleTasks(rrPost2, reqPost2)

	if rrPost2.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rrPost2.Code)
	}

	bodyPost3 := strings.NewReader(`{"title":"exmpl3"}`)
	reqPost3 := httptest.NewRequest(http.MethodPost, "/tasks", bodyPost3)
	rrPost3 := httptest.NewRecorder()

	h.HandleTasks(rrPost3, reqPost3)

	if rrPost3.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rrPost3.Code)
	}

	reqGet := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rrGet := httptest.NewRecorder()

	h.HandleTasks(rrGet, reqGet)

	if rrGet.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rrGet.Code)
	}

	var tasks []model.Task
	if err := json.NewDecoder(rrGet.Body).Decode(&tasks); err != nil {
		t.Fatal("failed to decode", err)
	}

	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}
}

func TestHandleTasksGETOffsetOutOfRange(t *testing.T) {
	h := newTestHandler()

	bodyPost1 := strings.NewReader(`{"title":"exmpl1"}`)
	reqPost1 := httptest.NewRequest(http.MethodPost, "/tasks", bodyPost1)
	rrPost1 := httptest.NewRecorder()

	h.HandleTasks(rrPost1, reqPost1)

	if rrPost1.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rrPost1.Code)
	}

	bodyPost2 := strings.NewReader(`{"title":"exmpl2"}`)
	reqPost2 := httptest.NewRequest(http.MethodPost, "/tasks", bodyPost2)
	rrPost2 := httptest.NewRecorder()

	h.HandleTasks(rrPost2, reqPost2)

	if rrPost2.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rrPost2.Code)
	}

	reqGet := httptest.NewRequest(http.MethodGet, "/tasks?limit=10&offset=100", nil)
	rrGet := httptest.NewRecorder()

	h.HandleTasks(rrGet, reqGet)

	if rrGet.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rrGet.Code)
	}

	var tasks []model.Task
	if err := json.NewDecoder(rrGet.Body).Decode(&tasks); err != nil {
		t.Fatal("failed to decode", err)
	}

	if len(tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(tasks))
	}
}
