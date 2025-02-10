package storage

import (
	"context"
	"database/sql"

	"example.com/m/config"
	"example.com/m/pkg/logging"
	postgresDB "example.com/m/pkg/postgres"
	"github.com/jmoiron/sqlx"
)

type PostgresStorage struct {
	db  *sqlx.DB
	log *logging.Logger
}

func NewPostgresStorage(cfg *config.Config, log *logging.Logger) *PostgresStorage {
	dsn := cfg.GetDSN()
	s, err := postgresDB.New(dsn)
	if err != nil {
		log.Fatalf("Cannot connect to postgres")
	}
	return &PostgresStorage{s, log}
}

func (s *PostgresStorage) GetShortURL(ctx context.Context, longURL string) (string, error) {
	var shortURL string
	query := `SELECT short_url FROM url_mappings WHERE long_url = $1`
	err := s.db.GetContext(ctx, &shortURL, query, longURL)
	if err == sql.ErrNoRows {
		return "", ErrNotFound
	}
	return shortURL, err
}

func (s *PostgresStorage) GetLongURL(ctx context.Context, shortURL string) (string, error) {
	var longURL string
	query := `SELECT long_url FROM url_mappings WHERE short_url = $1`
	err := s.db.GetContext(ctx, &longURL, query, shortURL)
	if err == sql.ErrNoRows {
		return "", ErrNotFound
	}
	return longURL, err
}

func (s *PostgresStorage) GetLongURLByID(ctx context.Context, id uint32) (string, error) {
	var longURL string
	query := `SELECT long_url FROM url_mappings WHERE id = $1`
	err := s.db.GetContext(ctx, &longURL, query, id)
	if err == sql.ErrNoRows {
		return "", ErrNotFound
	}
	return longURL, err
}

func (s *PostgresStorage) SaveURL(ctx context.Context, id uint32, longURL, shortURL string) error {
	query := `INSERT INTO url_mappings (id, long_url, short_url) VALUES ($1, $2, $3) 
	          ON CONFLICT (id) DO NOTHING`
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = s.db.ExecContext(ctx, query, id, longURL, shortURL)

	return tx.Commit()
}

func (s *PostgresStorage) Close() error {
	return s.db.Close()
}
