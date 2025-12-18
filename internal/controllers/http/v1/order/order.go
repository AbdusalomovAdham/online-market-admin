package order

import (
	"context"
	order "main/internal/services/order"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service order.Service
}

func NewController(service order.Service) Controller {
	return Controller{service: service}
}

func (as Controller) CreateOrder(c *gin.Context) {
	var data order.Create
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	ctx := context.Background()
	if err := as.service.Create(ctx, data, authHeader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "ok!"})
}

func (as Controller) GetOrderList(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	ctx := context.Background()
	orderList, err := as.service.GetList(ctx, authHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orderList, "message": "ok!"})
}

func (as Controller) GetOrderById(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	orderIdStr := c.Param("id")
	if orderIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is missing"})
		return
	}

	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Order ID"})
		return
	}

	ctx := context.Background()
	order, err := as.service.GetById(ctx, int64(orderId), authHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order, "message": "ok!"})
}

func (as Controller) DeleteOrder(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	orderIdStr := c.Param("id")
	if orderIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order Id is missing"})
		return
	}

	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Order ID"})
		return
	}

	ctx := context.Background()
	if err := as.service.Delete(ctx, int64(orderId), authHeader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!"})
}
