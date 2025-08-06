package auth

import (
	"net/http"
	"fmt"
	"encoding/json"
)

type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func GetGoogleUserInfo(w http.ResponseWriter, r *http.Request) (*GoogleUserInfo, error) {
	var req struct {
		AccessToken string `json:"access_token"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.AccessToken == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		fmt.Println("Error at json decoding stage")
		return nil, fmt.Errorf("Error! %w", err)
	}

	userInfoResp, err := http.Get("https://www.googleapis.com/oauth2/v3/userinfo?access_token=" + req.AccessToken)
	if err != nil || userInfoResp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to validate access token", http.StatusUnauthorized)
		fmt.Printf("Failed to fetch userinfo. Error: %v\n", err)
		return nil, fmt.Errorf("failed to fetch userinfo: %w", err)
	}
	defer userInfoResp.Body.Close()

	var userInfo GoogleUserInfo
	err = json.NewDecoder(userInfoResp.Body).Decode(&userInfo)
	if err != nil {
		http.Error(w, "Error parsing userinfo response", http.StatusInternalServerError)
		fmt.Println("Error decoding userinfo JSON:", err)
		return nil, fmt.Errorf("error decoding userinfo: %w", err)
	}

	if !userInfo.EmailVerified {
		http.Error(w, "Email not verified", http.StatusUnauthorized)
		fmt.Println("Email is not verified")
		return nil, fmt.Errorf("email not verified")
	}

	return &userInfo, nil
}
