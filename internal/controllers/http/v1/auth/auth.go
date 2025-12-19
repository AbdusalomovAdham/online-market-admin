package auth

import (
	"context"
	"main/internal/services/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service auth.Service
}

func NewController(service *auth.Service) Controller {
	return Controller{
		service: *service,
	}
}

func (as Controller) SignIn(c *gin.Context) {
	var signIn auth.SignIn

	if err := c.ShouldBindJSON(&signIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	token, err := as.service.SignIn(ctx, signIn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!", "token": token})
}
