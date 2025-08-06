package main

import (
	"log"
	"net/http"

	"github.com/WilliamCWhite/tallykeeper_backend/auth"
	"github.com/WilliamCWhite/tallykeeper_backend/db"
	"github.com/WilliamCWhite/tallykeeper_backend/handlers"

	"github.com/gorilla/mux"
	// "github.com/joho/godotenv"
)

func main() {
	// Only necessary if not using dockerl
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	r := mux.NewRouter()
	r.Use(auth.CORSResolver)
	r.Host("localhost:5173")

	// Routes
	r.HandleFunc("/api/auth/google", handlers.LoginHandler)

	// Protected Router requiring authorization key
	pr := r.PathPrefix("/api").Subrouter() // all these routes start with api
	pr.Use(auth.JWTVerifier)

	pr.HandleFunc("/lists", handlers.ListsHandler)
	pr.HandleFunc("/entries/{listID}", handlers.EntriesHandler)

	db.InitializeDB()
	

	log.Fatal(http.ListenAndServe(":7070", r))
}
