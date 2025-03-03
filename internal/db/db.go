package db

import (
	"context"
	"os"

	"github.com/ibLu247/wareman.git/internal/logger"
	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func ConnectDB() {
	var err error
	Conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Logger.Fatal("Не удалось подключиться к бд")
	}
}

func DisconnectDB() {
	Conn.Close(context.Background())
}
