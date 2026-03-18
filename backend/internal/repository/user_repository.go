package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
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

// func UpdateUser() (models.User, error) {}
// func DeleteUser() error {}
