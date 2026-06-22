package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username,omitempty"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
}

var users = []User{
	{
		Id:       1,
		Username: "Kamal",
		Age:      25,
		Email:    "kamal@mail.com",
	},
	{
		Id:       2,
		Username: "Jamal",
		Age:      20,
		Email:    "jamal@mail.com",
	},
}

func connectDB() {
	connString := "postgres://postgres:raju92@localhost:5432/go-crud"

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to database successfully")
	db = conn
}

func main() {
	connectDB()
	defer db.Close(context.Background())
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /createUser", createHandler)
	mux.HandleFunc("GET /users", getUsersHandler)
	mux.HandleFunc("GET /users/{id}", getSingleUsersHandler)
	mux.HandleFunc("PUT /users/{id}", updateUsersHandler)
	mux.HandleFunc("DELETE /users/{id}", deleteUserHandler)
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
	var newUser User

	err := json.NewDecoder(r.Body).Decode(&newUser)

	if err != nil {
		fmt.Println("Error decoding request body", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid request body")
		return
	}
	fmt.Printf("Received new user: %+v\n", newUser)

	// newUser.Id = len(users) + 1
	// users = append(users, newUser)

	query := "INSERT INTO users (username, age, email) VALUES ($1, $2, $3) RETURNING id"

	err = db.QueryRow(context.Background(), query, newUser.Username, newUser.Age, newUser.Email).Scan(&newUser.Id)

	if err != nil {
		fmt.Println("Error inserting user into database", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to create user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)

}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {

	query := "SELECT id, username, age, email FROM users"
	rows, err := db.Query(context.Background(), query)

	if err != nil {
		fmt.Println("Error fetching users from database", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to fetch users")
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Username, &user.Age, &user.Email)
		if err != nil {
			fmt.Println("Error scanning user row", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Failed to scan user")
			return
		}
		users = append(users, user)
	}

	err = rows.Err()

	if err != nil {
		fmt.Println("Error occurred while iterating over rows", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to fetch users")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)

}

func getSingleUsersHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid user id")
		return
	}

	var user User

	query := "SELECT id, username, age, email FROM users WHERE id = $1"
	err = db.QueryRow(context.Background(), query, id).Scan(&user.Id, &user.Username, &user.Age, &user.Email)

	if err == pgx.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "User not found")
		return
	}

	if err != nil {
		fmt.Println("Error fetching user from database", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to fetch user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func updateUsersHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid user id")
		return
	}

	var updatedUser User

	err = json.NewDecoder(r.Body).Decode(&updatedUser)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid requset body")
		return
	}

	query := "UPDATE users SET username = $1, age = $2, email = $3 WHERE id = $4 RETURNING id"
	err = db.QueryRow(context.Background(), query, updatedUser.Username, updatedUser.Age, updatedUser.Email, id).Scan(&updatedUser.Id)

	if err == pgx.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "User not found")
		return
	}

	if err != nil {
		fmt.Println("Error updating user in database", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to update user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid user id")
		return
	}

	query := "DELETE FROM users WHERE id = $1 RETURNING id"

	cmdTag, err := db.Exec(context.Background(), query, id)

	if err != nil {
		fmt.Println("Error deleting user from database", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to delete user")
		return
	}

	if cmdTag.RowsAffected() == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "User not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintln(w, "User deleted successfully")
}
