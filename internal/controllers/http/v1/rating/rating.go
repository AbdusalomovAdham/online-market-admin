package rating

import (
	"context"
	rating "main/internal/services/rating"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service rating.Service
}

func NewController(service rating.Service) Controller {
	return Controller{service: service}
}

func (as Controller) CreateRating(c *gin.Context) {
	var data rating.Create

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	productIdStr := c.Param("id")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Product id must be number"})
		return
	}
	data.ProductId = int64(productId)

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}
	ctx := context.Background()
	ratingId, err := as.service.CreateRating(ctx, data, authHeader)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messsage": "ok!", "id": ratingId})

}
