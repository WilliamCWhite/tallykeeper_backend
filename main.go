package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the home page")
} 

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := GetGooglePayload(w, r)
	if (err != nil) {
		
		fmt.Printf("error from GetGooglePayload: %v", err)
		return
	}

	email := payload.Claims["email"].(string)
	fmt.Println("Email: ")
	fmt.Println(email)

	var userID = "27"

	tokenString, err := GenerateJWT(userID)
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
	w.Header().Set("Content-Type", "text/plain")
	userID := r.Context().Value("userID").(string)
	fmt.Println(userID)

	json.NewEncoder(w).Encode(map[string]string{
		"data_from_backend": "good job",
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter()
	r.Use(CORSResolver)
	r.Host("localhost:5173")

	// Routes
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/auth/google", AuthHandler)

	// Protected Router requiring authorization key
	pr := r.PathPrefix("/api").Subrouter() // all these routes start with api
	pr.Use(JWTVerifier)

	pr.HandleFunc("/tokentest", tokentestHandler)

	log.Fatal(http.ListenAndServe(":7070", r))
}
