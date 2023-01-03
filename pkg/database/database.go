package database

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Client struct {
	db *sql.DB
}

// Open is the function to open a connection with database and set some poll configurations
func Open(ctx context.Context, connString string) (*Client, error) {
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(1 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	return &Client{db: db}, nil
}

// Close closes underlying connection pool
func (c *Client) Close() {
	c.db.Close()
}

// Query executes the provided query as a prepared statement
func (c *Client) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return c.db.QueryContext(ctx, query, args...)
}

// QueryRow executes the provided query as a prepared statement
func (c *Client) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return c.db.QueryRowContext(ctx, query, args...)
}

// Exec executes the provided query as a prepared statement
func (c *Client) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return c.db.ExecContext(ctx, query, args...)
}
