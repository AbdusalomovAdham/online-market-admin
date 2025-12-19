package user

import (
	"context"
	"main/internal/entity"
	user "main/internal/services/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service user.Service
}

func NewController(service user.Service) Controller {
	return Controller{service: service}
}

func (as Controller) AdminCreateUser(c *gin.Context) {
	var data user.Create

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authHeader := c.GetHeader("Authorization")

	file, err := c.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx := context.Background()
	if file != nil {
		avatarFile, err := as.service.Upload(ctx, file, "../media/avatar")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		data.Avatar = &avatarFile.Path
	}

	id, err := as.service.AdminCreateUser(ctx, data, authHeader)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!", "id": id})
}

func (as Controller) AdminGetList(c *gin.Context) {
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
	users, count, err := as.service.GetAll(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!", "data": users, "count": count})
}

func (as Controller) AdminGetById(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ID must be number!",
		})
		return
	}

	ctx := context.Background()
	user, err := as.service.GetById(ctx, idInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!", "data": user})
}

func (as Controller) AdminUpdateUser(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID must be number!",
		})
		return
	}

	var data user.Update
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx := context.Background()
	file, err := c.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if file != nil {
		avatarFile, err := as.service.Upload(ctx, file, "../media/avatar")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		data.Avatar = &avatarFile.Path
	}

	authHeader := c.GetHeader("Authorization")
	err = as.service.AdminUserUpdate(ctx, idInt, data, authHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!"})
}

func (as Controller) AdminDeleteUser(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID must be number!",
		})
		return
	}

	authHeader := c.GetHeader("Authorization")
	ctx := context.Background()
	err = as.service.AdminUserDelete(ctx, int64(idInt), authHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!"})
}
