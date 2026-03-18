package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/martins87/ngp-app/internal/repository"
)

type UserHandler struct {
	Repo *repository.UserRepository
}

func (h *UserHandler) GetUsers (c *gin.Context) {
	users, err := h.Repo.GetUsers(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, users)
}

// func GetUser() {}
// func CreateUser() {}
// func UpdateUser() {}
// func DeleteUser() {}
