package users

import (
	"database/sql"
	"errors"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(email, hash string) error {
	_, err := r.db.Exec(
		`INSERT INTO users (email, hashed_password) VALUES ($1,$2)`,
		email, hash,
	)
	return err
}

func (r *Repo) ByEmail(email string) (*User, error) {
	u := &User{}
	err := r.db.QueryRow(
		`SELECT id, email, hashed_password FROM users WHERE email=$1`,
		email,
	).Scan(&u.ID, &u.Email, &u.Hash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return u, err
}
