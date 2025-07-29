package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/WilliamCWhite/tallykeeper_backend/auth"
	"github.com/WilliamCWhite/tallykeeper_backend/db"
	"github.com/WilliamCWhite/tallykeeper_backend/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the home page")
} 

func tokentestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value("userID").(string)
	fmt.Println(userID)

	lists, err := db.GetListsByUserID(r.Context(), 2)
	if err != nil {
		fmt.Printf("error getting lists: %v\n", err)
	}
	fmt.Println(lists)
	
	testresult, err := db.GetUserIDByEmail(r.Context(), "poopemail")
	fmt.Println(testresult)
	fmt.Println(err)



	json.NewEncoder(w).Encode(lists)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter()
	r.Use(auth.CORSResolver)
	r.Host("localhost:5173")

	// Routes
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/auth/google", handlers.LoginHandler)

	// Protected Router requiring authorization key
	pr := r.PathPrefix("/api").Subrouter() // all these routes start with api
	pr.Use(auth.JWTVerifier)

	pr.HandleFunc("/tokentest", tokentestHandler)
	pr.HandleFunc("/lists", handlers.ListsHandler)
	pr.HandleFunc("/entries/{listID}", handlers.EntriesHandler)

	db.InitializeDB()
	

	log.Fatal(http.ListenAndServe(":7070", r))
}
