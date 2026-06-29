package middlerware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ride/db/table"
)

// UserIDAuth 校验请求头 X-User-Id：必须带上且数据库中存在该用户ID，否则报错。
// 校验通过后将 user 存入 context（键 "user"）供后续 handler 使用。
func UserIDAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Request.Header.Get("X-User-Id")
		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing X-User-Id header"})
			return
		}
		user, err := table.FindUserByUserID(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid X-User-Id"})
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
