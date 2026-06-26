package middlerware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Code") != "BonVoyage" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		// 继续处理请求
		c.Next()
	}
}
