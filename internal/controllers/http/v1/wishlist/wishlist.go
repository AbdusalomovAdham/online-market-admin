package wishlist

import (
	"context"
	"main/internal/entity"
	"main/internal/services/wishlist"
	"main/internal/utils"
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

func (as Controller) AdminWishistGetList(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
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

	lang := c.GetHeader("Accept-langueage")
	if lang == "" {
		lang = "uz"
		filter.Language = &lang
	} else {
		filter.Language = &lang
	}

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.Background()
	wishlistItems, count, err := as.service.AdminWishistGetList(ctx, authHeader, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!", "data": wishlistItems, "count": count})
}

func (as Controller) AdminWishlistCreate(c *gin.Context) {
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
	id, err := as.service.AdminWishlistCreate(ctx, productId, authHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "message": "ok!"})
}

func (as Controller) AdminWishlistDelete(c *gin.Context) {
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
	if err := as.service.AdminWishlistDelete(ctx, int64(wishlistId), authHeader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!"})
}
