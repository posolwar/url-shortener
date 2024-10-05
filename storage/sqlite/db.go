package sqlite

import (
	"context"
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(ctx context.Context, storagePath string) (*Storage, error) {
	const op = "storage.sqlite.NewStorage" // Имя текущей функции для логов и ошибок

	db, err := sql.Open("sqlite3", storagePath) // Подключаемся к БД
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// TODO: Позже вынести в полноценный механизм миграции
	_, err = db.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);`,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}
