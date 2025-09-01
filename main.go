package main

import (
	"fmt"
	"gotask/database"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	//call config function to create postgres db connection
	err := database.Config()
	if err != nil {
		fmt.Println("Error during connect with postgres database : ", err)
		return
	}

	m := mux.NewRouter()
	m.HandleFunc("/user/register", database.UserRegister).Methods("POST")
	m.HandleFunc("/login", database.Login).Methods("POST")

	http.ListenAndServe(":8088", m)
}
