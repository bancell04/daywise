package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/bancell04/daywise/backend/db"
	"github.com/bancell04/daywise/backend/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

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

	r.HandleFunc("/task", handlers.UploadTask).Methods("POST", "OPTIONS")
	r.HandleFunc("/task/{id}", handlers.DeleteTask).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/tasks", handlers.GetTasks).Methods("GET", "OPTIONS")
	r.HandleFunc("/tasks/{start}/to/{end}", handlers.GetTasksByDay).Methods("GET", "OPTIONS")
	r.HandleFunc("/categories", handlers.GetCategories).Methods("GET", "OPTIONS")
	r.HandleFunc("/categories", handlers.PostCategories).Methods("POST", "OPTIONS")
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
