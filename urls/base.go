package urls

import (
	"msc24x/showdown/internal/utils"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

var url_map = map[string]string{}
var url_rmap = map[string]string{}

func AttachRouter(r *gin.Engine) {
	router = r
}

func Url(name string) string {
	url, found := url_map[name]
	utils.BPanicIf(!found, "unable to parse '%s' to a url", name)

	return url
}

func RUrl(url string) string {
	name, found := url_rmap[url]
	utils.BPanicIf(!found, "unable to reverse '%s' to a url name", url)

	return name
}

func GET(url string, name string, handlers ...gin.HandlerFunc) gin.IRoutes {
	url_map[name] = url
	url_rmap[url] = name
	return router.GET(url, handlers...)
}

func POST(url string, name string, handlers ...gin.HandlerFunc) gin.IRoutes {
	url_map[name] = url
	url_rmap[url] = name
	return router.POST(url, handlers...)
}
