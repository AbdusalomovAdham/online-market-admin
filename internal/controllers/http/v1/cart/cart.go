package cart

import (
	"context"
	cart "main/internal/services/cart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service cart.Service
}

func NewController(service cart.Service) Controller {
	return Controller{service: service}
}

func (as Controller) CreateCart(c *gin.Context) {
	var cart cart.Create
	if err := c.ShouldBind(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authHeader := c.GetHeader("Authorization")
	ctx := context.Background()
	cartId, err := as.service.Create(ctx, cart, authHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": cartId, "message": "ok!"})

}

func (as Controller) UpdateCartItemTotal(c *gin.Context) {

	cartItemIdStr := c.Param("id")
	if cartItemIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart id is missing"})
		return
	}

	orderId, err := strconv.Atoi(cartItemIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Order ID"})
		return
	}

	authHeader := c.GetHeader("Authorization")
	ctx := context.Background()
	if err := as.service.UpdateCartItemTotal(ctx, int64(orderId), authHeader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!"})
}

func (as Controller) DeleteCartItem(c *gin.Context) {
	cartItemIdStr := c.Param("id")
	if cartItemIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart id is missing"})
		return
	}

	cartItemId, err := strconv.Atoi(cartItemIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Cart Item id"})
		return
	}

	authHeader := c.GetHeader("Authorization")
	ctx := context.Background()

	if err := as.service.DeleteCartItem(ctx, int64(cartItemId), authHeader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!"})
}

func (as Controller) GetCartList(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	ctx := context.Background()

	cartItems, err := as.service.GetList(ctx, authHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cartItems, "message": "ok!"})
}
