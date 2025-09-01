package database

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func createToken(userId int, username string, mailId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id":  userId,
			"username": username,
			"mail-id":  mailId,
			"exp":      time.Now().Add(time.Hour * 1).Unix(),
		})

	// fmt.Print("ser : ", jwtSecretKey)
	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
