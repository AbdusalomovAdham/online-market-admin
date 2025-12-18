package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	auth Auth
}

func NewMiddleware(auth Auth) *Middleware {
	return &Middleware{auth: auth}
}

func (am *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("authorization")

		if authHeader == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Token invalid"})
			ctx.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if authHeader == tokenStr {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Token invalid"})
			ctx.Abort()
			return
		}

		claims, err := am.auth.IsValidToken(ctx.Request.Context(), tokenStr)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid or expired"})
			ctx.Abort()
			return
		}
		ctx.Set("claims", claims)
		ctx.Next()
	}
}
