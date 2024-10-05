package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"url-shortener/storage"

	"github.com/mattn/go-sqlite3"
)

// SaveURL - Сохранение ссылки и алиаса в БД
func (s *Storage) SaveURL(ctx context.Context, url *url.URL, alias string) (id int64, err error) {
	const op = "storage.sqlite.SaveURL" // Имя текущей функции для логов и ошибок

	res, err := s.db.ExecContext(
		ctx,
		`INSERT OR IGNORE INTO url(alias, url) VALUES(?,?)`,
		alias,
		url.String(),
	)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// GetURL - Получение ссылки по алиасу из БД
func (s *Storage) GetURL(ctx context.Context, alias string) (*url.URL, error) {
	const op = "storage.sqlite.GetURL" // Имя текущей функции для логов и ошибок

	var urlStr string

	err := s.db.QueryRowContext(
		ctx,
		`SELECT url FROM url WHERE alias =?`,
		alias,
	).Scan(&urlStr)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
	} else if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return url.Parse(urlStr)
}

// DeleteURL - Удаление ссылки по алиасу из БД
func (s *Storage) DeleteURL(ctx context.Context, alias string) error {
	const op = "storage.sqlite.DeleteURL" // Имя текущей функции для логов и ошибок

	res, err := s.db.ExecContext(
		ctx,
		`DELETE FROM url WHERE alias =?`,
		alias,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
	}

	return nil
}
