package main

import (
	"context"
	"fmt"
	"go-http-crud/db"
	"go-http-crud/handler"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		panic(err)
	}

	db.ConnectDB()
	defer db.Db.Close(context.Background())
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.RootHandler)
	mux.HandleFunc("GET /health", handler.HealthHandler)
	mux.HandleFunc("POST /createUser", handler.CreateHandler)
	mux.HandleFunc("GET /users", handler.GetUsersHandler)
	mux.HandleFunc("GET /users/{id}", handler.GetSingleUsersHandler)
	mux.HandleFunc("PUT /users/{id}", handler.UpdateUsersHandler)
	mux.HandleFunc("DELETE /users/{id}", handler.DeleteUserHandler)
	fmt.Println("Server is running at port 5000")
	err = http.ListenAndServe(":5000", mux)
	if err != nil {
		fmt.Println("Server error", err)
	}

}
