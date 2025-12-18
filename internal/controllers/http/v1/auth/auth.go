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

func (as Controller) SendOtp(c *gin.Context) {
	var phoneNumer auth.PhoneRequest

	if err := c.ShouldBindJSON(&phoneNumer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	token, err := as.service.SendOTP(ctx, phoneNumer.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!", "token": token})
}

func (as Controller) ConfirmOTP(c *gin.Context) {
	var dataOTP auth.GetOTP
	if err := c.ShouldBind(&dataOTP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	firstName, token, err := as.service.ConfirmOTP(ctx, dataOTP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!", "firstName": firstName, "token": token})
}

func (as Controller) UpdateInfo(c *gin.Context) {
	var getInfo auth.GetInfo

	if err := c.ShouldBindJSON(&getInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	token, err := as.service.CreateInfo(ctx, getInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!", "token": token})
}
