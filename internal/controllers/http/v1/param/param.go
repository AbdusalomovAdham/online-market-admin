package param

import (
	"context"
	"main/internal/entity"
	param "main/internal/services/param"
	"main/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service param.Service
}

func NewController(service param.Service) Controller {
	return Controller{service: service}
}

// Param create
func (as *Controller) AdminParamCreate(c *gin.Context) {
	var paramData param.Create
	authHeader := c.GetHeader("Authorization")

	if err := c.ShouldBindJSON(&paramData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	paramId, err := as.service.CreateParam(ctx, paramData, authHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": paramId, "message": "ok!"})
}

// Param by Id
func (as Controller) AdminParamGetById(c *gin.Context) {
	var paramId int64
	paramIdStr := c.Param("id")
	if paramIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "param id is required"})
		return
	}
	paramId, err := strconv.ParseInt(paramIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	param, err := as.service.GetParamById(c, paramId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": param, "message": "ok!"})
}

// Param get list
func (ac *Controller) AdminParamGetList(c *gin.Context) {
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
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Offset must be number!",
			})
			return
		}

		filter.Offset = &queryInt
	}

	order, err := utils.GetQuery(c, "order")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	filter.Order = order

	search, err := utils.GetQuery(c, "search")
	if err != nil {
		return
	}
	filter.Search = search

	ctx := context.Background()
	lang := c.GetHeader("Accept-Language")

	if lang == "" {
		lang = "uz"
		filter.Language = &lang
	} else {
		filter.Language = &lang
	}

	categories, count, err := ac.service.ParamGetList(ctx, filter)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok!", "data": map[string]any{"results": categories, "count": count}})
}

// Param delete
func (ac *Controller) AdminParamDelete(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	paramIdStr := c.Param("id")
	if paramIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "param id is required"})
		return
	}
	paramId, err := strconv.ParseInt(paramIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	if err := ac.service.DeleteParam(ctx, paramId, authHeader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!"})
}

// Param update
func (ac *Controller) AdminParamUpdate(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	paramIdStr := c.Param("id")
	if paramIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "param id is required"})
		return
	}
	paramId, err := strconv.ParseInt(paramIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var data param.UpdateParam
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	if err := ac.service.UpdateParam(ctx, paramId, data, authHeader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok!"})
}
