package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"go-http-crud/db"
	"go-http-crud/modle"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to go server")
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Server is up and healthy")
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	var newUser modle.User

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

	err = db.Db.QueryRow(context.Background(), query, newUser.Username, newUser.Age, newUser.Email).Scan(&newUser.Id)

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

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {

	query := "SELECT id, username, age, email FROM users"
	rows, err := db.Db.Query(context.Background(), query)

	if err != nil {
		fmt.Println("Error fetching users from database", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to fetch users")
		return
	}
	defer rows.Close()

	var users []modle.User
	for rows.Next() {
		var user modle.User
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

func GetSingleUsersHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid user id")
		return
	}

	var user modle.User

	query := "SELECT id, username, age, email FROM users WHERE id = $1"
	err = db.Db.QueryRow(context.Background(), query, id).Scan(&user.Id, &user.Username, &user.Age, &user.Email)

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

func UpdateUsersHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid user id")
		return
	}

	var updatedUser modle.User

	err = json.NewDecoder(r.Body).Decode(&updatedUser)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid requset body")
		return
	}

	query := "UPDATE users SET username = $1, age = $2, email = $3 WHERE id = $4 RETURNING id"
	err = db.Db.QueryRow(context.Background(), query, updatedUser.Username, updatedUser.Age, updatedUser.Email, id).Scan(&updatedUser.Id)

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

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid user id")
		return
	}

	query := "DELETE FROM users WHERE id = $1 RETURNING id"

	cmdTag, err := db.Db.Exec(context.Background(), query, id)

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
