package api

import (
	"msc24x/showdown/config"

	"github.com/gin-gonic/gin"
)

// Middleware to enable usage of request header Access-Token.
func AccessToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Access-Token")

		if token != config.ACCESS_TOKEN {
			ctx.String(401, "Unauthorized")
			ctx.Abort()
			return
		}
	}
}
