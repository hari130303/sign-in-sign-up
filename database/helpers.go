package database

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func ErrorResponse(w http.ResponseWriter, e error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{"error": e.Error()}

	jsonByte, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(jsonByte)
}

func SuccessResponse(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	jsonByte, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(jsonByte)
}

func HashPassword(password string) (string, error) {
	// Recommended cost for production is at least 12.
	// Adjust based on your security requirements and system capabilities.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
