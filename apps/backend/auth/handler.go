package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"foreningsliv/backend/db"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// LoginHandler handles POST /auth/login.
// Validates email+password against the database and returns a JWT.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "invalid request body"})
		return
	}

	if req.Email == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "email and password are required"})
		return
	}

	// Look up auth + profile by email
	var profileID, name, storedPassword string
	err := db.DB.QueryRow(`
		SELECT p.id, p.name, a.password
		FROM auth a
		JOIN profile p ON p.id = a.profile_id
		WHERE a.email = $1
	`, req.Email).Scan(&profileID, &name, &storedPassword)

	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResponse{Error: "invalid email or password"})
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Error: "internal server error"})
		return
	}

	// Compare passwords (plaintext for now -- should use bcrypt in production)
	if storedPassword != req.Password {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResponse{Error: "invalid email or password"})
		return
	}

	// Generate JWT with profile ID and name
	token, err := GenerateToken(profileID, name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Error: "failed to generate token"})
		return
	}

	json.NewEncoder(w).Encode(loginResponse{Token: token})
}
