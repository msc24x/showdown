package api

import (
	"log"

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
	router.POST("/judge", Judge)
	router.POST("/tmp", Tmp)
	router.GET("/stats", GetStats)
}
