package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

var Conn *pgx.Conn

func ConnectDB() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	var err error
	Conn, err = pgx.Connect(context.Background(), "postgres://postgres:password@localhost:5432/postgres")
	if err != nil {
		logger.Fatal("Не удалось подключиться к бд")
	}
}

func DisconnectDB() {
	Conn.Close(context.Background())
}
