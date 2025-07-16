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
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value("userID").(string)
	fmt.Println(userID)

	lists, err := GetListsByUserID(r.Context(), 2)
	if err != nil {
		fmt.Printf("error getting lists: %v\n", err)
	}
	fmt.Println(lists)


	// LIST CREATE WORKS
	// var list1 List
	// list1.Title = "Database Created list"
	// list1.UserID = 2
	// new_id, err := CreateList(r.Context(), list1)
	// if err != nil {
	// 	fmt.Printf("error creating list: %v", err)
	// } else {
	// 	fmt.Println(new_id)
	// }


	// LIST UPDATE WORKS
	// list1 := List {
	// 	ListID: 6,
	// 	Title: "Database Updated List",
	// 	UserID: 2,
	// }
	// err = UpdateList(r.Context(), list1)
	// if err != nil {
	// 	fmt.Printf("error updating list: %v", err)
	// }

	// LIST DELETE WORKS
	// err = DeleteList(r.Context(), 6, 2)
	// if err != nil {
	// 	fmt.Printf("error deleting list: %v", err)
	// }

	json.NewEncoder(w).Encode(lists)
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

	initializeDB()

	log.Fatal(http.ListenAndServe(":7070", r))
}
