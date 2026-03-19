package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/martins87/ngp-app/internal/db"
	"github.com/martins87/ngp-app/internal/handlers"
	"github.com/martins87/ngp-app/internal/repository"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	pool, err := db.NewPool()
	if err != nil {
		log.Fatal(err)
	}

	repo := &repository.UserRepository{DB: pool}

	if err := repo.CreateTable(context.Background()); err != nil {
		log.Fatal("Failed to create table:", err)
	}

	handler := &handlers.UserHandler{Repo: repo}

	r := gin.Default()

	r.GET("/health", handler.HealthCheck)
	r.GET("/users", handler.GetUsers)
	r.GET("/users/:id", handler.GetUser)
	r.POST("/users", handler.CreateUser)
	r.PATCH("/users/:id", handler.UpdateUser)
	r.DELETE("users/:id", handler.DeleteUser)

	r.Run(":8080")
}
