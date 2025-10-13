package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {

	port := getPort()

	http.Handle("/", http.FileServer(http.Dir("./web")))

	log.Printf("Сервер запущен на порту %d", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
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
