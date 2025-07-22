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
	
	// entries, err := db.GetEntriesByListID(r.Context(), 4)
	// if err != nil {
	// 	fmt.Printf("error getting entries: %v\n", err)
	// }
	// fmt.Println(entries)
	
	entry1 := db.Entry {
		Name: "database_entry",
		Score: 100,
		ListID: 5,
	}
	newID, err := db.CreateEntry(r.Context(), entry1)
	if err != nil {
		fmt.Printf("error creating enttry: %v", err)
	}
	fmt.Println(newID)

	entry2 := db.Entry {
		Name: "updated_database_entry",
		Score: 200,
		ListID: 5,
		EntryID: 10,
	}
	err = db.UpdateEntry(r.Context(), entry2)
	if err != nil {
		fmt.Printf("error updating entry: %v", err)
	}

	err = db.DeleteEntry(r.Context(), 11, 5)




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
