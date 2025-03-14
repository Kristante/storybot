package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

func CreateDatabasePool() *pgxpool.Pool {
	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("Unable to connect to database with error: %v\n", err)
		os.Exit(1)
	}
	return pool
}

func SelectHandleFromDatabase(pool *pgxpool.Pool, message string) (string, error) {
	var result string
	err := pool.QueryRow(context.Background(), "SELECT answer FROM handlers WHERE handle = $1", message).Scan(&result)

	if err != nil {
		fmt.Printf("Проблемы с чтением из базы данных: %v\n", err)
		return "", err
	}
	return result, nil
}

func AddHandleFromDatabase(pool *pgxpool.Pool, handleMessage string, answerMessage string) error {
	_, err := pool.Exec(context.Background(), "INSERT into handlers (handle, answer) values ($1, $2)", handleMessage, answerMessage)
	if err != nil {
		return errors.New("произошла ошибка при добавлении handle")
	}
	return nil
}

func RemoveHandleFromDatabase(pool *pgxpool.Pool, handleMessage string) error {
	_, err := pool.Exec(context.Background(), "DELETE from handlers WHERE handle = $1", handleMessage)
	if err != nil {
		return errors.New("данный handle отсутствует")
	}
	return nil
}
