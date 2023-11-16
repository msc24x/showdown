package api

import "github.com/gin-gonic/gin"

func WriteServerError(c *gin.Context, msg string) {
	c.Writer.WriteHeader(500)
	c.Writer.Write([]byte(msg))
}

func WriteBadRequest(c *gin.Context, msg string) {
	c.Writer.WriteHeader(400)
	c.Writer.Write([]byte(msg))
}

// Initializes the API routes
func Init(router *gin.Engine) {

	router.POST("/judge", Judge)
	router.GET("/stats", GetStats)

}
