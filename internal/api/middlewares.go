package api

import (
	"github.com/msc24x/showdown/internal/config"

	"github.com/gin-gonic/gin"
)

// Middleware to enable usage of request header Access-Token.
func AccessToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Access-Token")
		wh_secret := ctx.GetHeader("Webhook-Secret")

		if token != config.ACCESS_TOKEN && wh_secret != config.ACCESS_TOKEN {
			ctx.String(401, "Unauthorized")
			ctx.Abort()
			return
		}
	}
}
