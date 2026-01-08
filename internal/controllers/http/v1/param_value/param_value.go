package param_value

import (
	"context"
	"main/internal/entity"
	"main/internal/services/param_value"
	"main/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service param_value.Service
}

func NewController(service param_value.Service) Controller {
	return Controller{service: service}
}

func (as Controller) AdminParamValueCreate(c *gin.Context) {
	var paramValueData param_value.Create

	if err := c.ShouldBind(&paramValueData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authHeader := c.GetHeader("Authorization")
	ctx := context.Background()
	id, err := as.service.ParamValueCreate(ctx, paramValueData, authHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "ok!"})
}

func (as Controller) AdminParamValueGetList(c *gin.Context) {
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

	list, count, err := as.service.ParamValueGetList(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!", "data": map[string]any{"count": count, "results": list}})
}

func (ac *Controller) AdminParamValueDelete(c *gin.Context) {
	paramValueId := c.Param("id")
	id, err := strconv.ParseInt(paramValueId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authHeader := c.GetHeader("Authorization")
	ctx := context.Background()
	err = ac.service.ParamValueDelete(ctx, id, authHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!"})
}

func (ac *Controller) AdminParamValueGetById(c *gin.Context) {
	paramId := c.Param("id")
	id, err := strconv.ParseInt(paramId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	paramValue, err := ac.service.ParamValueGetById(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!", "data": paramValue})
}

func (ac *Controller) AdminParamValueUpdate(c *gin.Context) {
	paramValueId := c.Param("id")
	id, err := strconv.ParseInt(paramValueId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var data param_value.Update
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authHeader := c.GetHeader("Authorization")
	ctx := context.Background()
	err = ac.service.ParamValueUpdate(ctx, id, data, authHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!"})
}
