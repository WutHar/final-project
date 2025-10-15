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
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
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

	now := time.Now()

	if err := processTaskDate(&task, now); err != nil {
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
func processTaskDate(task *db.Task, now time.Time) error {
	today := now.Format("20060102")

	fmt.Printf("DEBUG: Input - date: %s, title: %s, repeat: %s\n", task.Date, task.Title, task.Repeat)
	if task.Date == "today" {
		task.Date = today
	}

	if task.Date == "" {
		task.Date = today
	}

	parsedDate, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("некорректный формат даты")
	}
	if task.Repeat != "" {

		if task.Date == today && task.Repeat == "d 1" {

			task.Date = today
		} else {

			nextDate, err := date.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("некорректное правило повторения: %v", err)
			}
			task.Date = nextDate
		}
	} else {

		if parsedDate.Before(now) {
			task.Date = today
		}
	}

	//fmt.Printf("DEBUG: Output - date: %s\n", task.Date)
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
