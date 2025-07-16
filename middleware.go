package main

import (
	"net/http"
	"strings"
	"fmt"
	"context"
)

func CORSResolver(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		//unsure if necessary
		// w.Header().Set("Access-Control-Allow-Credentials")

		if (r.Method == "OPTIONS") {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	}) 
}

func JWTVerifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Expected format: "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // unchanged from trimprefix
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		claims, err := VerifyJWT(tokenString)
		if (err != nil) {
			http.Error(w, fmt.Sprintf("Unauthorized: %v", err), http.StatusUnauthorized)
			return
		}

		// serves to next handler with context
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
		// context can be read in next handler layer
		// as r.Context().Value("userID").(string)
	}) 
} 
