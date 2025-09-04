package users

import (
	"context"
	"database/sql"
	"errors"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, email, hash string) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO users (email, hashed_password) VALUES ($1,$2) RETURNING id`,
		email, hash,
	).Scan(&id)
	return id, err
}

func (r *Repo) ByEmail(ctx context.Context, email string) (*User, error) {
	u := &User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, hashed_password FROM users WHERE email=$1`,
		email,
	).Scan(&u.ID, &u.Email, &u.Hash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return u, err
}
