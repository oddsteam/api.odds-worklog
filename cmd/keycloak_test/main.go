package main

import (
	"fmt"
	"log"
	"net/http"

	"gitlab.odds.team/worklog/api.odds-worklog/pkg/auth"
)

func main() {
	// Create a new validator
	validator := auth.NewKeycloakValidator(
		"http://localhost:9000", // Replace with your Keycloak server URL
		"odds",                  // Realm name
		"worklog",               // Client ID
	)

	// Create a test endpoint
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		token, err := auth.ExtractTokenFromHeader(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Validate the token
		claims, err := validator.ValidateToken(token)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Token is valid, return the claims
		fmt.Fprintf(w, "Token is valid!\n")
		fmt.Fprintf(w, "User: %s\n", claims.PreferredUsername)
		fmt.Fprintf(w, "Email: %s\n", claims.Email)
		fmt.Fprintf(w, "Roles: %v\n", claims.RealmAccess.Roles)
	})

	// Start the server
	fmt.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
