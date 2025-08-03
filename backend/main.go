package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bancell04/daywise/backend/db"
	"github.com/bancell04/daywise/backend/models"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
}

func uploadTask(w http.ResponseWriter, r *http.Request) {

}

func getLogs(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Pool.Query(context.Background(), `SELECT * FROM TASKS`)
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var t models.Task
		err := rows.Scan(&t.ID, &t.Title, &t.Category, &t.StartTime, &t.EndTime)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, t)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
		return
	}
}

func main() {
	db.Connect()
	db.Setup()
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/ping", pingHandler).Methods("GET")

	r.HandleFunc("/task", uploadTask).Methods("POST")

	r.HandleFunc("/tasks", getLogs).Methods("GET")

	handler := cors.Default().Handler(r)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
