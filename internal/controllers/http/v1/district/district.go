package district

import (
	"context"
	"main/internal/entity"
	"main/internal/services/district"
	"main/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *district.Service
}

func NewController(service *district.Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (ac *Controller) AdminDistrictGetList(c *gin.Context) {
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
	districtList, count, err := ac.service.GetDistricts(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"data": map[string]any{
			"results": districtList,
			"count":   count,
		},
	})
}

func (ac *Controller) AdminDistrictGetListByRegionId(c *gin.Context) {
	var filter entity.Filter

	lang := c.GetHeader("Accept-language")
	if lang == "" {
		lang = "uz"
	}
	filter.Language = &lang

	ctx := context.Background()
	regionId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Region ID must be number"})
		return
	}

	districtList, err := ac.service.GetDistrictsByRegionId(ctx, filter, regionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"data": map[string]any{
			"results": districtList,
		},
	})
}
