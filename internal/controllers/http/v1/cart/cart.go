package cart

import (
	"context"
	"main/internal/entity"
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
	filter := entity.Filter{}
	query := c.Request.URL.Query()
	lang := c.GetHeader("Accept-Language")
	if lang == "" {
		lang = "uz"
	}
	filter.Language = &lang

	limitQ := query["limit"]
	if len(limitQ) > 0 {
		queryInt, err := strconv.Atoi(limitQ[0])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Limit must be number!",
			})
			return
		}

		filter.Limit = &queryInt
	}

	offsetQ := query["offset"]
	if len(offsetQ) > 0 {
		queryInt, err := strconv.Atoi(offsetQ[0])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "Offset must be number!",
			})

			return
		}
		filter.Offset = &queryInt
	}

	orderQ := query["order"]
	if len(orderQ) > 0 {
		filter.Order = &orderQ[0]
	}

	ctx := context.Background()
	cartItems, total, err := as.service.GetList(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cartItems, "message": "ok!", "count": total})
}
