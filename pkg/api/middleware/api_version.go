package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

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
