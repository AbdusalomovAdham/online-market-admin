package product

import (
	"context"
	"main/internal/entity"
	product "main/internal/services/product"
	"main/internal/utils"
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
	var productData product.Create
	if err := c.ShouldBind(&productData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	images := form.File["images"]
	ctx := context.Background()
	pid := int32(1)
	if len(images) > 0 {
		imgFile, err := as.service.MultipleUpload(ctx, images, "../media/products", &pid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		productData.Images = imgFile
	}
	authHeader := c.GetHeader("Authorization")

	id, err := as.service.CreateProduct(c, productData, authHeader)
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
	lang := c.GetHeader("Accept-Language")

	if lang == "" {
		lang = "uz"
	}
	filter.Language = &lang

	categoryId := c.Query("category_id")
	if categoryId != "" {
		categoryInt, err := strconv.Atoi(categoryId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Category id must be number!",
			})
			return
		}

		categoryInt64 := int64(categoryInt)
		filter.CategoryId = &categoryInt64
	}

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

	search, err := utils.GetQuery(c, "search")
	if err != nil {
		return
	}
	filter.Search = search

	ctx := context.Background()

	list, count, err := as.service.GetList(ctx, filter)
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

func (as Controller) UpdateProduct(c *gin.Context) {
	productId := c.Param("id")
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Product id must be number!",
		})
		return
	}

	var data product.Update
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	images := form.File["images"]
	ctx := context.Background()

	authHeader := c.GetHeader("Authorization")
	err = as.service.UpdateProduct(ctx, productIdInt, data, authHeader, images)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
	})
}

func (as Controller) AdminDeleteProduct(c *gin.Context) {
	productId := c.Param("id")
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Product id must be number!",
		})
		return
	}

	ctx := context.Background()
	authHeader := c.GetHeader("Authorization")
	err = as.service.AdminDeleteProduct(ctx, int64(productIdInt), authHeader)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
	})
}
