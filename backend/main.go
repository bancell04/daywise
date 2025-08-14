package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

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

	// If task.ID is zero, let PostgreSQL generate it
	query := `
        INSERT INTO tasks (id, title, category_id, start, "end")
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (id) DO UPDATE
        SET title = EXCLUDED.title,
            category_id = EXCLUDED.category_id,
            start = EXCLUDED.start,
            "end" = EXCLUDED.end
        RETURNING id
    `

	id := task.ID // use zero for new tasks
	err = db.Pool.QueryRow(context.Background(), query,
		id, task.Title, task.Category, task.Start, task.End,
	).Scan(&id)

	if err != nil {
		http.Error(w, "Failed to insert/update task: "+err.Error(), http.StatusInternalServerError)
		return
	}

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
	rows, err := db.Pool.Query(context.Background(), `SELECT * FROM tasks`)
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
	startStr := vars["start"]
	endStr := vars["end"]

	rows, err := db.Pool.Query(
		context.Background(),
		`SELECT id, title, category_id, start, "end" FROM tasks WHERE start >= $1 AND "end" <= $2`,
		startStr,
		endStr,
	)
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Category, &t.Start, &t.End); err != nil {
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

func getCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Pool.Query(context.Background(), `SELECT * FROM categories`)
	if err != nil {
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	categories := []models.Category{}
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name, &c.Color)
		if err != nil {
			log.Fatal(err)
		}
		categories = append(categories, c)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(categories); err != nil {
		http.Error(w, "Failed to encode categories", http.StatusInternalServerError)
		return
	}
}

func postCategories(w http.ResponseWriter, r *http.Request) {
	var categories []models.Category
	err := json.NewDecoder(r.Body).Decode(&categories)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	existingRows, err := db.Pool.Query(context.Background(), `SELECT id FROM categories`)
	if err != nil {
		http.Error(w, "Failed to fetch existing categories", http.StatusInternalServerError)
		return
	}
	defer existingRows.Close()

	existingIDs := map[int]bool{}
	for existingRows.Next() {
		var id int
		existingRows.Scan(&id)
		existingIDs[id] = true
	}

	incomingIDs := map[int]bool{}

	for _, c := range categories {
		if c.ID != nil {
			_, err := db.Pool.Exec(context.Background(), `
				INSERT INTO categories (id, name, color)
				VALUES ($1, $2, $3)
				ON CONFLICT (id) DO UPDATE
				SET name = EXCLUDED.name,
					color = EXCLUDED.color;
			`, *c.ID, c.Name, c.Color)
			if err != nil {
				http.Error(w, "Failed to upsert category: "+err.Error(), http.StatusInternalServerError)
				return
			}
			incomingIDs[*c.ID] = true
		} else {
			var newID int
			err := db.Pool.QueryRow(context.Background(), `
				INSERT INTO categories (name, color)
				VALUES ($1, $2)
				RETURNING id
			`, c.Name, c.Color).Scan(&newID)
			if err != nil {
				http.Error(w, "Failed to insert category: "+err.Error(), http.StatusInternalServerError)
				return
			}
			incomingIDs[newID] = true
		}
	}

	for id := range existingIDs {
		if !incomingIDs[id] {
			_, err := db.Pool.Exec(context.Background(), `DELETE FROM categories WHERE id=$1`, id)
			if err != nil {
				http.Error(w, "Failed to delete category: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
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

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	result, err := db.Pool.Exec(context.Background(), `DELETE FROM tasks WHERE id=$1`, id)
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected() == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}

func main() {
	db.Connect()
	db.Setup()
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/ping", pingHandler).Methods("GET", "OPTIONS")

	r.HandleFunc("/task", uploadTask).Methods("POST", "OPTIONS")

	r.HandleFunc("/task/{id}", deleteTask).Methods("DELETE", "OPTIONS")

	r.HandleFunc("/tasks", getLogs).Methods("GET", "OPTIONS")

	r.HandleFunc("/tasks/{start}/to/{end}", getLogsByDay).Methods("GET", "OPTIONS")

	r.HandleFunc("/categories", getCategories).Methods("GET", "OPTIONS")

	r.HandleFunc("/categories", postCategories).Methods("POST", "OPTIONS")

	r.HandleFunc("/db-reset", resetDatabaseHandler).Methods("GET", "OPTIONS")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "X-Admin-Secret"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
