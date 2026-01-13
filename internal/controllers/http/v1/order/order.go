package order

import (
	"context"
	"main/internal/entity"
	order "main/internal/services/order"
	"main/internal/utils"
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

func (as Controller) AdminOrderCreate(c *gin.Context) {
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
	if err := as.service.AdminOrderCreate(ctx, data, authHeader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "ok!"})
}

func (as Controller) AdminOrderGetList(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	var filter entity.Filter
	query := c.Request.URL.Query()

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

	order, err := utils.GetQuery(c, "order")
	if err != nil {
		return
	}
	filter.Order = order

	lang := c.GetHeader("Accept-Language")
	if lang == "" {
		langType := "uz"
		filter.Language = &langType
	} else {
		filter.Language = &lang
	}

	ctx := context.Background()
	orderList, count, err := as.service.AdminOrderGetList(ctx, authHeader, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": map[string]any{
		"results": orderList,
		"count":   count,
	}, "message": "ok!"})
}

func (as Controller) AdminOrderGetById(c *gin.Context) {
	var filter entity.Filter
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

	lang := c.GetHeader("Accept-Language")
	if lang == "" {
		langType := "uz"
		filter.Language = &langType
	} else {
		filter.Language = &lang
	}

	ctx := context.Background()
	order, err := as.service.AdminOrderGetById(ctx, int64(orderId), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order, "message": "ok!"})
}

func (as Controller) AdminOrderDelete(c *gin.Context) {
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
	if err := as.service.AdminOrderDelete(ctx, int64(orderId), authHeader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!"})
}

func (as Controller) AdminOrderUpdate(c *gin.Context) {
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

	var updateData order.Update
	if err := c.ShouldBind(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	if err := as.service.AdminOrderUpdate(ctx, updateData, int64(orderId), authHeader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!"})
}

func (as Controller) AdminOrderDeleteOrderItem(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	itemIdStr := c.Param("id")
	if itemIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order Id is missing"})
		return
	}

	itemId, err := strconv.Atoi(itemIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Item ID"})
		return
	}

	ctx := context.Background()
	if err := as.service.AdminOrderDeleteOrderItem(ctx, int64(itemId), authHeader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!"})
}
