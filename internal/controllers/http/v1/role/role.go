package role

import (
	"context"
	"main/internal/entity"
	role "main/internal/services/role"
	"main/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service role.Service
}

func NewController(service role.Service) Controller {
	return Controller{service: service}
}

func (as *Controller) AdminGetList(c *gin.Context) {
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

	// lang := c.GetHeader("Accept-Language")
	// if lang == "" {
	// 	langType := "uz"
	// 	filter.Language = &langType
	// } else {
	// 	filter.Language = &lang
	// }
	//
	//
	ctx := context.Background()
	orderList, count, err := as.service.AdminRoleList(ctx, filter)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orderList, "message": "ok!", "count": count})

}
