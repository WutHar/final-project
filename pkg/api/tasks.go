package api

import (
	"final-project/pkg/db"
	"net/http"
	"strconv"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	limit := 50
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	tasks, err := db.Tasks(limit)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}
