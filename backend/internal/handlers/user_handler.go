package handlers

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/martins87/ngp-app/internal/dto"
	"github.com/martins87/ngp-app/internal/models"
	"github.com/martins87/ngp-app/internal/repository"
)

type UserHandler struct {
	Repo *repository.UserRepository
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@")
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

func (h *UserHandler) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
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

func (h *UserHandler) UpdateUser(c *gin.Context) {
	// Parse ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}

	// Bind input DTO
	var input dto.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	// Input validation
	if input.Name == nil && input.Email == nil {
		c.JSON(400, gin.H{"error": "no fields to update"})
		return
	}

	if input.Name != nil && *input.Name == "" {
		c.JSON(400, gin.H{"error": "name cannot be empty"})
		return
	}

	if input.Email != nil && !isValidEmail(*input.Email) {
		c.JSON(400, gin.H{"error": "invalid email"})
		return
	}

	// Call service/repo
	user, err := h.Repo.UpdateUser(c.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	// Return updated user
	c.JSON(200, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Parse id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}

	// Call service/repo
	err = h.Repo.DeleteUser(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}

		// Log internal error
		log.Println("DeleteUser error:", err)

		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	// Return success code
	c.Status(200)
}
