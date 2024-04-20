package api

import (
	"msc24x/showdown/config"

	"github.com/gin-gonic/gin"
)

// Middleware to enable usage of request header Access-Token.
func AccessToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Access-Token")
		wh_secret := ctx.GetHeader("Webhook-Secret")

		if token != config.ACCESS_TOKEN && wh_secret != config.WEBHOOK_SECRET {
			ctx.String(401, "Unauthorized")
			ctx.Abort()
			return
		}
	}
}
