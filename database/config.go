package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB
var jwtSecretKey string

func Config() error {

	DBHOST := os.Getenv("DBHOST")
	DBPORT := os.Getenv("DBPORT")
	DBUSER := os.Getenv("DBUSER")
	DBPASS := os.Getenv("DBPASS")
	DBNAME := os.Getenv("DBNAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		DBHOST, DBPORT, DBUSER, DBPASS, DBNAME)

	fmt.Println("connection string : ", psqlInfo)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		// panic(err)
		return err
	}
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		// panic(err)
		return err

	}

	// initialise jwt secret key
	jwtSecretKey = os.Getenv("JWTSECRETKET")

	fmt.Println("Successfully connected!")

	return nil
}
