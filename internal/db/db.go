package db

import (
	"context"

	"github.com/ibLu247/wareman.git/internal/logger"
	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func ConnectDB() {
	var err error
	Conn, err = pgx.Connect(context.Background(), "postgres://postgres:password@localhost:5432/postgres")
	if err != nil {
		logger.Logger.Fatal("Не удалось подключиться к бд")
	}
}

func DisconnectDB() {
	Conn.Close(context.Background())
}
