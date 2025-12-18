package wishlist

import (
	"context"
	"main/internal/services/wishlist"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service wishlist.Service
}

func NewController(service *wishlist.Service) Controller {
	return Controller{
		service: *service,
	}
}

func (as Controller) WishList(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.Background()
	wishlistItems, err := as.service.GetList(ctx, authHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!", "data": wishlistItems})
}

func (as Controller) Create(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var productId wishlist.Create
	if err := c.ShouldBind(&productId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ctx := context.Background()
	id, err := as.service.Create(ctx, productId, authHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "message": "ok!"})
}

func (as Controller) Delete(c *gin.Context) {
	paramsStr := c.Param("id")
	if paramsStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	wishlistId, err := strconv.Atoi(paramsStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid param id"})
		return
	}

	ctx := context.Background()
	if err := as.service.Delete(ctx, int64(wishlistId), authHeader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!"})
}
