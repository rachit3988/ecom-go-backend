package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	// "ecom-go-backend/internal/models/models"
	"github.com/jackc/pgx/v5"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var Db *pgx.Conn // Set this from main.go or pass db instance properly

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Hash the user password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Insert user into DB
	_, err = Db.Exec(context.Background(),
		"INSERT INTO users (email, password) VALUES ($1, $2)", req.Email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User registered successfully")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Login endpoint not implemented yet.")
}
