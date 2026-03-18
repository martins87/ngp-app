package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/martins87/ngp-app/internal/dto"
	"github.com/martins87/ngp-app/internal/models"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func (r *UserRepository) CreateTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);`
	_, err := r.DB.Exec(ctx, query)
	return err
}

func (r *UserRepository) GetUsers(ctx context.Context) ([]models.User, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, name, email, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) GetUser(ctx context.Context, id int) (models.User, error) {
	query := "SELECT id, name, email, created_at FROM users WHERE id = $1"
	var u models.User
	err := r.DB.QueryRow(ctx, query, id).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	return u, err
}

func (r *UserRepository) CreateUser(ctx context.Context, u models.User) (models.User, error) {
	query := `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id, name, email, created_at
	`
	err := r.DB.QueryRow(ctx, query, u.Name, u.Email).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	if err != nil {
		return models.User{}, err
	}

	return u, nil
}

func (r *UserRepository) UpdateUser(
	ctx context.Context,
	id int,
	input dto.UpdateUserInput,
) (models.User, error) {
	// Execute SQL
	query := `
		UPDATE users
		SET
			name = COALESCE($1, name),
			email = COALESCE($2, email)
		WHERE id = $3
		RETURNING id, name, email, created_at
	`

	// Map DB -> model
	var user models.User
	err := r.DB.QueryRow(
		ctx,
		query,
		input.Name,
		input.Email,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)

	// Detect 'not found'
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, err
	}

	return user, nil
}

// func DeleteUser() error {}
