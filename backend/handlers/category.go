package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bancell04/daywise/backend/db"
	"github.com/bancell04/daywise/backend/models"
)

func GetCategories(w http.ResponseWriter, r *http.Request) {
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

func PostCategories(w http.ResponseWriter, r *http.Request) {
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
