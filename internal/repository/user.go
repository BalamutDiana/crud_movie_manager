package repository

import (
	"context"
	"database/sql"

	"github.com/BalamutDiana/crud_movie_manager/internal/domain"
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{db}
}

func (r *Users) Create(ctx context.Context, user domain.User) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (name, email, password, registered_at) values ($1, $2, $3, $4)",
		user.Name, user.Email, user.Password, user.RegisteredAt)

	return err
}

func (r *Users) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRowContext(ctx, "SELECT id, name, email, registered_at FROM users WHERE email=$1 AND password=$2", email, password).
		Scan(&user.ID, &user.Name, &user.Email, &user.RegisteredAt)

	return user, err
}

func (r *Users) CheckUserExist(ctx context.Context, email string) (bool, error) {
	var user domain.User
	
	err := r.db.QueryRowContext(ctx, "SELECT id FROM users WHERE email=$1", email).
		Scan(&user.ID)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return true, err
	}
	return true, nil
}
