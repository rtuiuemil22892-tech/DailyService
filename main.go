package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})
	http.HandleFunc("/tasks", tasksHand)
	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

var tasks = []Task{}
var nextID = 1

func tasksHand(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		_ = json.NewEncoder(w).Encode(tasks)

	case http.MethodPost:
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
			return
		}

		task.Title = strings.TrimSpace(task.Title)
		if task.Title == "" {
			http.Error(w, "title is required", http.StatusBadRequest)
			return
		}

		task.ID = nextID
		nextID++
		tasks = append(tasks, task)

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(tasks)

	default:
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "method not allowed",
		})
	}
}
