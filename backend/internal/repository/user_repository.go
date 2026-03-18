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

func (r *UserRepository) GetUsers(ctx context.Context) ([]models.User, error){
	rows, err := r.DB.Query(ctx, "SELECT id, name, email, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
		users = append(users, u)
	}
	return users, nil
}

// func GetUser() {}
// func CreateUser() {}
// func UpdateUser() {}
// func DeleteUser() {}
