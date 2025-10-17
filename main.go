package main

import (
	"final-project/pkg/api"
	"final-project/pkg/date"
	"final-project/pkg/db"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	if err := db.Init("scheduler.db"); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}

	port := getPort()

	http.Handle("/", http.FileServer(http.Dir("./web")))

	http.HandleFunc("/api/nextdate", handleNextDate)
	http.HandleFunc("/api/task", api.TaskHandler)
	http.HandleFunc("/api/tasks", api.TasksHandler)
	http.HandleFunc("/api/task/done", api.TaskDoneHandler)

	log.Printf("Сервер запущен на порту %d", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func handleNextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Неверный формат даты 'now'", http.StatusBadRequest)
		return
	}

	result, err := date.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(result))
}

func getPort() int {
	portStr := os.Getenv("TODO_PORT")
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err == nil {
			return port
		}
		log.Printf("Неверный порт в TODO_PORT: %s, используем по умолчанию 7540", portStr)
	}
	return 7540
}
