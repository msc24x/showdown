package api

import (
	"log"
	"msc24x/showdown/urls"

	"github.com/gin-gonic/gin"
)

func WriteServerError(c *gin.Context, msg string) {
	c.Writer.WriteHeader(500)
	c.Writer.Write([]byte(msg))
}

func WriteBadRequest(c *gin.Context, msg string) {
	c.Writer.WriteHeader(400)
	c.Writer.Write([]byte(msg))
}

// Initializes the API routes
func AttachHandlers(router *gin.Engine) {
	log.Println("Attaching API handlers...")
	urls.AttachRouter(router)

	urls.POST("/judge", "judge", Judge)
	urls.POST("/workers/register", "workers-register", RegisterWorker)
	urls.GET("/status", "status", Status)

	urls.POST("/_debug/webhook", "debug-webhook", DebugWebhook)
}
