package category

import (
	"context"
	"main/internal/entity"
	category "main/internal/services/category"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service category.Service
}

func NewController(service category.Service) Controller {
	return Controller{service: service}
}

func (ac Controller) AdminCreateCategory(c *gin.Context) {
	var data category.Create

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.Background()
	categoryId, err := ac.service.Create(ctx, data, authHeader)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok!", "id": categoryId})
}

func (ac Controller) AdminGetCategoryById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
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
	category, err := ac.service.GetById(ctx, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok!", "data": category})
}

func (ac Controller) AdminGetCategoryList(c *gin.Context) {
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

	orderQ := query["order"]
	if len(orderQ) > 0 {
		filter.Order = &orderQ[0]
	}

	ctx := context.Background()
	lang := c.GetHeader("Accept-Language")

	categories, total, err := ac.service.GetList(ctx, filter, lang)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok!", "data": categories, "count": total})
}

func (ac Controller) AdminUpdateCategory(c *gin.Context) {
	var data category.Update
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
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
	if err := ac.service.Update(ctx, id, data, authHeader); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok!"})
}

func (ac Controller) AdminDeleteCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
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
	if err := ac.service.Delete(ctx, id, authHeader); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok!"})
}
