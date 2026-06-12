package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Id    int
	Name  string
	Age   int
	Email string
}

var users = []User{
	{
		Id:    1,
		Name:  "Kamal",
		Age:   25,
		Email: "kamal@mail.com",
	},
	{
		Id:    2,
		Name:  "Jamal",
		Age:   20,
		Email: "jamal@mail.com",
	},
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /createUser", createHandler)
	mux.HandleFunc("GET /users", getUsersHandler)

	fmt.Println("Server is running at port 5000")
	err := http.ListenAndServe(":5000", mux)
	if err != nil {
		fmt.Println("Server error", err)
	}

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to go server")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Server is up and healthy")
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User Create")
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data, _ := json.Marshal(users)

	w.Write(data)

}
