package task

import (
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/metadata"

	"github.com/gorilla/mux"
	taskpb "github.com/rnd/kudu/golang/protogen/task"
	"github.com/unrolled/render"
)

var response = render.New()

// Task hold task information.
type Task struct {
	TaskName  string `json:"task_name"`
	URL       string `json:"url,omitempty"`
	Tags      string `json:"tags,omitempty"`
	Notes     string `json:"notes,omitempty"`
	NotesMD   string `json:"notes_md,omitempty"`
	Date      string `json:"date"`
	Completed bool   `json:"completed"`
}

// List fetch user task list.
func List(w http.ResponseWriter, r *http.Request) {
	var err error

	token := r.Header.Get("Authorization")
	kuduAuth := metadata.NewContext(
		r.Context(),
		metadata.Pairs("authorization", token),
	)

	// TODO: implement task filtering (today, this week, this month)
	vars := mux.Vars(r)
	res := &taskpb.TodayTasksResponse{}
	switch vars["period"] {
	case "today":
		res, err = kuduClient.TodayTasks(kuduAuth, &taskpb.TodayTasksRequest{})
	case "week":
		response.JSON(w, http.StatusUnprocessableEntity,
			map[string]string{"error": "Filter not implemented"})
		return
	case "month":
		response.JSON(w, http.StatusUnprocessableEntity,
			map[string]string{"error": "Filter not implemented"})
		return
	default:
		response.JSON(w, http.StatusUnprocessableEntity,
			map[string]string{"error": "Filter not implemented"})
		return
	}
	if err != nil {
		response.JSON(w, http.StatusInternalServerError,
			map[string]string{"error": err.Error()})
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"status": res.Status.String(),
		"tasks":  res.Tasks,
	})
}

// Add creates new task record for today date.
func Add(w http.ResponseWriter, r *http.Request) {
	var err error

	token := r.Header.Get("Authorization")
	kuduAuth := metadata.NewContext(
		r.Context(),
		metadata.Pairs("authorization", token),
	)

	if r.Body == nil {
		response.JSON(w, http.StatusBadGateway,
			map[string]string{"error": "Invalid request payload"})
		return
	}

	var task Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		response.JSON(w, http.StatusBadGateway,
			map[string]string{"error": "Invalid request payload"})
		return
	}
	defer r.Body.Close()

	var req taskpb.AddRequest
	req.Task = &taskpb.Task{
		TaskName: task.TaskName,
		Tags:     task.Tags,
		Notes:    task.Notes,
		Url:      task.URL,
	}
	res, err := kuduClient.AddTask(kuduAuth, &req)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError,
			map[string]string{"error": err.Error()})
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{
		"status":  res.Status.String(),
		"task_id": res.Id,
	})
}

// Remove deletes given task id from database.
func Remove(w http.ResponseWriter, r *http.Request) {
	var err error

	token := r.Header.Get("Authorization")
	kuduAuth := metadata.NewContext(
		r.Context(),
		metadata.Pairs("authorization", token),
	)

	vars := mux.Vars(r)
	res, err := kuduClient.RemoveTask(kuduAuth, &taskpb.RemoveRequest{Id: vars["id"]})
	if err != nil {
		response.JSON(w, http.StatusInternalServerError,
			map[string]string{"error": err.Error()})
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{
		"status": res.Status.String(),
	})
}
