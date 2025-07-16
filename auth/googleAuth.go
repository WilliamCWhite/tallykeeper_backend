package auth

import (
	"net/http"
	"fmt"
	"os"
	"encoding/json"
	"context"

	"google.golang.org/api/idtoken"
)
func GetGooglePayload(w http.ResponseWriter, r *http.Request) (*idtoken.Payload, error) {
	
	var req struct {
		IDToken string `json:"id_token"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	fmt.Println(req.IDToken)
	if err != nil || req.IDToken == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		fmt.Println("Error at json decoding stage")
		return nil, fmt.Errorf("Error! %w", err)
	}

	// replace with your client ID
	payload, err := idtoken.Validate(context.Background(), req.IDToken, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		http.Error(w, "Invalid ID token", http.StatusUnauthorized)
		fmt.Println("Error at token validation stage")
		return nil, err
	}
	return payload, nil
} 
