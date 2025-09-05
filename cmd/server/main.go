package main

import (
	"fmt"
	"net/http"
)

// This function handles requests to the "/" path.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Go E-commerce backend!")
}

func main() {
	http.HandleFunc("/", homeHandler) // Register the handler function
	fmt.Println("Server running on http://localhost:8080")
	// Start the HTTP server on port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server failed:", err)
	}
}
