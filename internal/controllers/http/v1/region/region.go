package region

import (
	"context"
	"main/internal/entity"
	"main/internal/services/region"
	"main/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *region.Service
}

func NewController(service *region.Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (ac *Controller) AdminRegionGetList(c *gin.Context) {
	var filter entity.Filter
	query := c.Request.URL.Query()

	limitQ := query["limit"]
	if len(limitQ) > 0 {
		limit, err := strconv.Atoi(limitQ[0])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Limit must be number"})
			return
		}
		filter.Limit = &limit
	}

	offset := query["offset"]
	if len(offset) > 0 {
		offset, err := strconv.Atoi(offset[0])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Offset must be number"})
			return
		}
		filter.Offset = &offset
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
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	filter.Search = search

	lang := c.GetHeader("Accept-language")
	if lang == "" {
		lang = "uz"
	}
	filter.Language = &lang
	ctx := context.Background()
	regionList, count, err := ac.service.GetRegions(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"data": map[string]any{
			"results": regionList,
			"count":   count,
		},
	})
}
