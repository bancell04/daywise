package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/bancell04/daywise/backend/db"
	"github.com/bancell04/daywise/backend/models"
	"github.com/gorilla/mux"
)

func UploadTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if task.ID != nil {
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

		id := task.ID
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
	} else {
		query := `
			INSERT INTO TASKS (title, category_id, start, "end")
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`
		var newID int
		err = db.Pool.QueryRow(context.Background(), query,
			task.Title, task.Category, task.Start, task.End,
		).Scan(&newID)

		if err != nil {
			http.Error(w, "Failed to insert/update task: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"id":       newID,
			"title":    task.Title,
			"category": task.Category,
			"start":    task.Start,
			"end":      task.End,
		}

		json.NewEncoder(w).Encode(response)
	}
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
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

func GetTasksByDay(w http.ResponseWriter, r *http.Request) {
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

func DeleteTask(w http.ResponseWriter, r *http.Request) {
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
