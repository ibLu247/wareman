package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func ConnectDB() {
	_, err := pgx.Connect(context.Background(), "postgres://postgres:password@localhost:5432/postgres")
	if err != nil {
		fmt.Println("Ошибка подключения к серверу PostgreSQL")
	}
}

func DisconnectDB() {
	Conn.Close(context.Background())
}
