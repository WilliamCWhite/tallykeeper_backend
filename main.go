package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/WilliamCWhite/tallykeeper_backend/db"
	"github.com/WilliamCWhite/tallykeeper_backend/auth"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the home page")
} 

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := auth.GetGooglePayload(w, r)
	if (err != nil) {
		
		fmt.Printf("error from GetGooglePayload: %v", err)
		return
	}

	email := payload.Claims["email"].(string)
	fmt.Println("Email: ")
	fmt.Println(email)

	var userID = "27"

	tokenString, err := auth.GenerateJWT(userID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
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

	testEmail := "test_user_3"
	newId, err := db.CreateUserFromEmail(r.Context(), testEmail)
	if err != nil {
		fmt.Printf("error creating user: %v", err)
	}

	fmt.Println(newId)

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
	r.HandleFunc("/auth/google", AuthHandler)

	// Protected Router requiring authorization key
	pr := r.PathPrefix("/api").Subrouter() // all these routes start with api
	pr.Use(auth.JWTVerifier)

	pr.HandleFunc("/tokentest", tokentestHandler)

	db.InitializeDB()
	

	log.Fatal(http.ListenAndServe(":7070", r))
}
