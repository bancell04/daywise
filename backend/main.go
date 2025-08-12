package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

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
	w.Header().Set("Content-Type", "application/json")

	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	query := `
        INSERT INTO tasks (title, category, start, "end")
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	var id int
	err = db.Pool.QueryRow(context.Background(), query,
		task.Title, task.Category, task.Start, task.End,
	).Scan(&id)

	if err != nil {
		http.Error(w, "Failed to insert task", http.StatusInternalServerError)
		return
	}

	// Send back the inserted task with its new ID
	response := map[string]interface{}{
		"id":       id,
		"title":    task.Title,
		"category": task.Category,
		"start":    task.Start,
		"end":      task.End,
	}

	json.NewEncoder(w).Encode(response)
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
		err := rows.Scan(&t.ID, &t.Title, &t.Category, &t.Start, &t.End)
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

func getLogsByDay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]

	rows, err := db.Pool.Query(
		context.Background(),
		`SELECT * FROM your_table WHERE start::date = $1`,
		date,
	)

	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var t models.Task
		err := rows.Scan(&t.ID, &t.Title, &t.Category, &t.Start, &t.End)
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

func resetDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	secret := r.Header.Get("X-Admin-Secret")
	if secret != os.Getenv("RESET_SECRET") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	_, err := db.Pool.Exec(context.Background(), `TRUNCATE TABLE tasks RESTART IDENTITY`)
	if err != nil {
		http.Error(w, "Failed to reset DB", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Database reset"))
}

func main() {
	db.Connect()
	db.Setup()
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/ping", pingHandler).Methods("GET")

	r.HandleFunc("/task", uploadTask).Methods("POST", "OPTIONS")

	r.HandleFunc("/tasks", getLogs).Methods("GET")

	r.HandleFunc("/tasks/{date}", getLogsByDay).Methods("GET")

	r.HandleFunc("/db-reset", resetDatabaseHandler).Methods("GET", "OPTIONS")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "X-Admin-Secret"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
