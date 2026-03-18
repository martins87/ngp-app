package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/martins87/ngp-app/internal/models"
	"github.com/martins87/ngp-app/internal/repository"
)

type UserHandler struct {
	Repo *repository.UserRepository
}

func repositoryUser(input struct {
	Name  string
	Email string
}) models.User {
	return models.User{
		Name:  input.Name,
		Email: input.Email,
	}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.Repo.GetUsers(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, users)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.Repo.GetUser(c.Request.Context(), id)
	if err != nil {
		// TODO refine later with errors.Is
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	user, err := h.Repo.CreateUser(c.Request.Context(), repositoryUser(
		struct {
			Name  string
			Email string
		}(input),
	))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, user)
}

// func UpdateUser() {}
// func DeleteUser() {}
