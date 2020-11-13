package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// APIVersionResolver is a middleware which resolves the api version the request
// is coming for.
//
// It puts the version in context which can be accessed by next handlers.
func APIVersionResolver(c *gin.Context) {
	pathComponents := strings.Split(c.FullPath(), "/")

	// FullPath will be /api/version/subPath
	if pathComponents == nil || len(pathComponents) < 3 {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid path requested"})
		return
	}

	c.Set("version", pathComponents[2])

	c.Next()
}
