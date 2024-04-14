package api

import (
	"log"
	"msc24x/showdown/config"

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

	if config.INSTANCE_TYPE != config.T_WORKER {
		router.POST("/judge", Judge)
	}

	if config.INSTANCE_TYPE == config.T_MANAGER {
		router.POST("/workers/register", RegisterWorker)
	}

	router.POST("/tmp", Tmp)
	router.GET("/stats", GetStats)
}
