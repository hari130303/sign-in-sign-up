package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"time"
)

func UserRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var inputPayload struct {
		UserName string `json:"username"`
		Password string `json:"password"`
		MailId   string `json:"mail_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&inputPayload); err != nil {
		ErrorResponse(w, fmt.Errorf("invalid request body"), http.StatusBadRequest)
		return
	}

	fmt.Printf("User request payload: %+v\n", inputPayload)

	// Validation
	if inputPayload.UserName == "" || inputPayload.Password == "" || inputPayload.MailId == "" {
		ErrorResponse(w, fmt.Errorf("all fields are required"), http.StatusBadRequest)
		return
	}

	if _, err := mail.ParseAddress(inputPayload.MailId); err != nil {
		ErrorResponse(w, fmt.Errorf("invalid email format"), http.StatusBadRequest)
		return
	}

	// Check if user exists
	var exist bool
	err := db.QueryRow("SELECT true FROM user_master WHERE user_name = $1", inputPayload.UserName).Scan(&exist)
	if err != nil && err != sql.ErrNoRows {
		ErrorResponse(w, fmt.Errorf("database error"), http.StatusInternalServerError)
		return
	}
	if exist {
		ErrorResponse(w, fmt.Errorf("username already exists"), http.StatusBadRequest)
		return
	}

	// Prepare return data
	var ReturnData struct {
		UserId    int       `json:"id"`
		UserName  string    `json:"name"`
		MailId    string    `json:"email"`
		CreatedAt time.Time `json:"timestamp"`
	}
	ReturnData.UserName = inputPayload.UserName
	ReturnData.MailId = inputPayload.MailId
	ReturnData.CreatedAt = time.Now().UTC()

	// Hash password
	hashedPassword, err := HashPassword(inputPayload.Password)
	if err != nil {
		ErrorResponse(w, fmt.Errorf("failed to hash password"), http.StatusInternalServerError)
		return
	}

	// Insert into DB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = db.QueryRowContext(ctx, `
		INSERT INTO user_master (user_name, password, mail_id, created_at) 
		VALUES ($1, $2, $3, $4) RETURNING user_id
	`, inputPayload.UserName, hashedPassword, inputPayload.MailId, ReturnData.CreatedAt).Scan(&ReturnData.UserId)

	if err != nil {
		ErrorResponse(w, fmt.Errorf("failed to insert user"), http.StatusInternalServerError)
		return
	}

	SuccessResponse(w, ReturnData, http.StatusCreated)
}

// func validMailAddress(address string) (string, bool) {
// 	addr, err := mail.ParseAddress(address)
// 	if err != nil {
// 		return "", false
// 	}
// 	return addr.Address, true
// }

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var inputPayload struct {
		MailId   string `json:"mail_id"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&inputPayload); err != nil {
		ErrorResponse(w, fmt.Errorf("invalid request body"), http.StatusBadRequest)
		return
	}

	fmt.Printf("User request payload: %+v\n", inputPayload)

	// Validation
	if inputPayload.Password == "" || inputPayload.MailId == "" {
		ErrorResponse(w, fmt.Errorf("all fields are required"), http.StatusBadRequest)
		return
	}

	if _, err := mail.ParseAddress(inputPayload.MailId); err != nil {
		ErrorResponse(w, fmt.Errorf("invalid email format"), http.StatusBadRequest)
		return
	}

	// Insert into DB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// var dbPassword to store the password stored in database
	var dbPassword string

	// user struct for token payload
	var user struct {
		UserId   int    `json:"user_id"`
		UserName string `json:"user_name"`
		MailId   string `json:"mail_id"`
	}
	err := db.QueryRowContext(ctx, `
		select user_id,user_name, password, mail_id from user_master where mail_id = $1
	`, inputPayload.MailId).Scan(&user.UserId, &user.UserName, &dbPassword, &user.MailId)

	if err != nil {
		if err == sql.ErrNoRows {
			ErrorResponse(w, fmt.Errorf("invalid credentials"), http.StatusUnauthorized)
			return
		}
		ErrorResponse(w, fmt.Errorf("database error"), http.StatusInternalServerError)
		return
	}

	ok := VerifyPassword(inputPayload.Password, dbPassword)
	if !ok {
		ErrorResponse(w, fmt.Errorf("invalid credentials"), http.StatusUnauthorized)
		return
	}

	token, err := createToken(user.UserId, user.UserName, inputPayload.MailId)
	if err != nil {
		fmt.Println("error during token creation : ", err)
		ErrorResponse(w, fmt.Errorf("failed to create token"), http.StatusInternalServerError)
		return
	}

	var responseData struct {
		Token string `json:"token"`
	}

	responseData.Token = token

	SuccessResponse(w, responseData, http.StatusOK)
}
