// provides support for listening to requests.
package api

import (
	"log"

	"github.com/msc24x/showdown/internal/api/urls"

	"github.com/gin-gonic/gin"
)

const (
	HTTP_SERVER_ERR   = 500
	HTTP_BAD_REQ      = 400
	HTTP_TOO_MANY_REQ = 429
)

func WriteError(c *gin.Context, code int, msg string) {
	c.Writer.WriteHeader(code)
	c.Writer.Write([]byte(msg))
}

func WriteServerError(c *gin.Context, msg string) {
	c.Writer.WriteHeader(500)
	c.Writer.Write([]byte(msg))
}

func WriteBadRequest(c *gin.Context, msg string) {
	c.Writer.WriteHeader(400)
	c.Writer.Write([]byte(msg))
}

// Initializes the API routes.
func AttachHandlers(router *gin.Engine) {
	log.Println("Attaching API handlers...")
	urls.AttachRouter(router)

	urls.POST("/judge", "judge", Judge)
	urls.POST("/workers/register", "workers-register", RegisterWorker)
	urls.GET("/status", "status", Status)

	urls.POST("/_debug/webhook", "debug-webhook", DebugWebhook)
}
