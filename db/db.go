package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var Db *pgx.Conn

func ConnectDB() {
	connString := os.Getenv("DB_STRING")

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to database successfully")
	Db = conn
}
