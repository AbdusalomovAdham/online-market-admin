package products

import (
	"context"
	"log"
	"main/internal/entity"
	product "main/internal/services/product"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service product.Service
}

func NewController(service product.Service) Controller {
	return Controller{service: service}
}

func (as Controller) CreateProduct(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var data product.Create
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	images := form.File["images"]
	log.Println("Images uploaded successfully", images)
	ctx := context.Background()
	if len(images) > 0 {
		imgFile, err := as.service.MultipleUpload(ctx, images, "../media/products")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		data.Images = &imgFile
	}

	id, err := as.service.CreateProduct(c, data, authHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "ok!"})
}

func (as Controller) GetById(c *gin.Context) {
	productIdStr := c.Param("id")
	productId, err := strconv.ParseInt(productIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	product, err := as.service.GetById(ctx, productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!", "data": product})
}

func (as Controller) GetProductsList(c *gin.Context) {
	filter := entity.Filter{}
	query := c.Request.URL.Query()

	categoryId := c.Query("category_id")
	categoryInt, err := strconv.Atoi(categoryId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Category id must be number!",
		})
		return
	}

	categoryInt64 := int64(categoryInt)
	filter.CategoryId = &categoryInt64

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
	authHeader := c.GetHeader("Authorization")
	list, count, err := as.service.GetList(ctx, filter, authHeader)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"data": map[string]any{
			"results": list,
			"count":   count,
		},
	})
}
