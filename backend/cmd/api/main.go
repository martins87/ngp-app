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
	godotenv.Load()

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

	r.GET("/users", handler.GetUsers)
	r.GET("/users/:id", handler.GetUser)
	r.POST("/users", handler.CreateUser)

	r.Run(":8080")
}
