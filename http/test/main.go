package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
}

func main() {
	http.HandleFunc("/v4/oa/access_token", handleTokenRequest)

	port := 8085
	fmt.Printf("Server starting on port %d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func handleTokenRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Log request headers
	fmt.Println("Headers:")
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", name, value)
		}
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Log form values
	fmt.Println("Form data:")
	fmt.Printf("refresh_token: %s\n", r.FormValue("refresh_token"))
	fmt.Printf("app_id: %s\n", r.FormValue("app_id"))
	fmt.Printf("grant_type: %s\n", r.FormValue("grant_type"))

	// Check secret key header
	secretKey := r.Header.Get("secret_key")
	if secretKey != "" {
		http.Error(w, "Invalid secret key", http.StatusUnauthorized)
		return
	}

	// Create response
	response := TokenResponse{
		AccessToken:  "ce",
		RefreshToken: "MT",
		ExpiresIn:    "90000",
	}

	// Set content type and encode response
	w.Header().Set("Content-Type", "text/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
