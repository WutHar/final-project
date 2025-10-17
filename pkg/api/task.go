package api

import (
	"encoding/json"
	"final-project/pkg/date"
	"final-project/pkg/db"
	"fmt"
	"net/http"
	"time"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, map[string]interface{}{})
}

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	if err := db.CompleteTask(id); err != nil {
		writeJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, map[string]interface{}{})
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "Ошибка декодирования JSON", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJSONError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	today := date.Today()

	if err := processTaskDate(&task, today); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{"id": id})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "Ошибка декодирования JSON", http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		writeJSONError(w, "Не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJSONError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	today := date.Today()
	if err := processTaskDate(&task, today); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, map[string]interface{}{})
}

func processTaskDate(task *db.Task, today string) error {
	if task.Date != "" && task.Date != "today" {
		_, err := time.Parse("20060102", task.Date)
		if err != nil {
			return fmt.Errorf("некорректный формат даты")
		}
	}

	if task.Date == "today" || task.Date == "" {
		task.Date = today
	}

	if task.Date < today && task.Date != "" && task.Date != "today" {
		task.Date = today
	}

	if task.Repeat != "" {
		now := time.Now()
		_, err := date.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("некорректное правило повторения")
		}
		task.Date = today
	}

	return nil
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
