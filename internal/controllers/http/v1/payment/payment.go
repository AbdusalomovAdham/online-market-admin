package payment

import (
	"context"
	"main/internal/entity"
	payment "main/internal/services/payment"
	"main/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service payment.Service
}

func NewController(service payment.Service) Controller {
	return Controller{service: service}
}

func (ac Controller) AdminPaymentCreate(c *gin.Context) {
	var paymentData payment.Create

	if err := c.ShouldBind(&paymentData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.Background()
	paymentId, err := ac.service.AdminPaymentCreate(ctx, paymentData, authHeader)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok!", "id": paymentId})
}

func (ac Controller) AdminPaymentGetById(c *gin.Context) {
	paymentIdStr := c.Param("id")
	paymentId, err := strconv.ParseInt(paymentIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.Background()
	paymentData, err := ac.service.AdminPaymentGetById(ctx, paymentId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok!", "data": paymentData})
}

func (ac Controller) AdminPaymentGetList(c *gin.Context) {
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

	ctx := context.Background()
	lang := c.GetHeader("Acceptf-Language")

	paymentList, total, err := ac.service.AdminPaymentGetList(ctx, filter, lang)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok!", "data": paymentList, "count": total})
}

func (ac Controller) AmdinPaymentUpdate(c *gin.Context) {
	var paymentData payment.Update
	if err := c.ShouldBind(&paymentData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paymentIdParamStr := c.Param("id")
	paymentIdParam, err := strconv.ParseInt(paymentIdParamStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.Background()
	if err := ac.service.AmdinPaymentUpdate(ctx, paymentIdParam, paymentData, authHeader); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok!"})
}

func (ac Controller) AdminPaymentDelete(c *gin.Context) {
	paymentIdParamStr := c.Param("id")
	paymentIdParam, err := strconv.ParseInt(paymentIdParamStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.Background()
	if err := ac.service.AdminPaymentDelete(ctx, paymentIdParam, authHeader); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok!"})
}
