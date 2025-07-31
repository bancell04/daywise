package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bancell04/daywise/backend/db"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
}

func main() {
	db.Connect()
	db.Setup()
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/ping", pingHandler).Methods("GET")

	handler := cors.Default().Handler(r)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
