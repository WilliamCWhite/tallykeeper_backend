package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/WilliamCWhite/tallykeeper_backend/auth"
	"github.com/WilliamCWhite/tallykeeper_backend/db"
)

// This endpoint uses google sign in to either log in as a user or create a new user then log in as them, providing a JWT token in the response
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := auth.GetGooglePayload(w, r)
	if err != nil {
		fmt.Printf("error from GetGooglePayload: %v", err)
		http.Error(w, "Failed to communicate with Google login service", http.StatusInternalServerError)
		return
	}

	email := payload.Claims["email"].(string)

	userID, err := db.GetUserIDByEmail(r.Context(), email)
	if err != nil {
		fmt.Printf("error from GetUserIDByEmail: %v", err)
		http.Error(w, "Failed to log in user", http.StatusInternalServerError)
		return
	}
	// userID is -1 if the email doesn't exist, so create new user
	if userID == -1 {
		fmt.Printf("User %v didn't exist, so creating new user", email)

		userID, err = db.CreateUserFromEmail(r.Context(), email)
		if err != nil {
			fmt.Printf("error from CreateUserFromEmail: %v", err)
			http.Error(w, "Failed to create new user", http.StatusInternalServerError)
		}
	}

	tokenString, err := auth.GenerateJWT(userID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	fmt.Printf("logged in as %v\n", email)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}
