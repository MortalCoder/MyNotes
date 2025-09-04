package notes

import (
	"context"
	"database/sql"
)

type Repo struct{ db *sql.DB }

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, userID int64, title, body string) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO notes (user_id, title, body) VALUES ($1,$2,$3) RETURNING id`,
		userID, title, body,
	).Scan(&id)
	return id, err
}

func (r *Repo) Get(ctx context.Context, userID, noteID int64) (*Note, error) {
	n := &Note{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, title, body FROM notes WHERE id=$1 AND user_id=$2`,
		noteID, userID,
	).Scan(&n.ID, &n.UserID, &n.Title, &n.Body)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return n, err
}

func (r *Repo) Update(ctx context.Context, userID, noteID int64, title, body string) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`UPDATE notes SET title=$1, body=$2 WHERE id=$3 AND user_id=$4`,
		title, body, noteID, userID,
	)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *Repo) Delete(ctx context.Context, userID, noteID int64) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM notes WHERE id=$1 AND user_id=$2`,
		noteID, userID,
	)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *Repo) List(ctx context.Context, userID int64, limit, offset int) ([]Note, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, title, body
		   FROM notes
		  WHERE user_id=$1
		  ORDER BY id DESC
		  LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Body); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}
