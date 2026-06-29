package middlerware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 最低要求的 APP 版本号和更新说明
const minAppVersion = "1.0"
const updateMessage = "1. 新增标记点语音播报功能\n2. 新增模拟导航\n3. 修复地图定位问题\n请下载最新版本使用"

// AppVersionCheck 校验请求头 X-App-Version，低于最低版本返回 409 + 更新内容
func AppVersionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		appVersion := c.Request.Header.Get("X-App-Version")
		if appVersion == "" {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"message": updateMessage,
			})
			return
		}
		if compareVersion(appVersion, minAppVersion) < 0 {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"message": updateMessage,
			})
			return
		}
		c.Next()
	}
}

// compareVersion 比较两个语义化版本号字符串
// 返回: 1 表示 v1 > v2, -1 表示 v1 < v2, 0 表示相等
func compareVersion(v1, v2 string) int {
	s1 := splitVersion(v1)
	s2 := splitVersion(v2)
	maxLen := len(s1)
	if len(s2) > maxLen {
		maxLen = len(s2)
	}
	for i := 0; i < maxLen; i++ {
		var n1, n2 int
		if i < len(s1) {
			n1 = s1[i]
		}
		if i < len(s2) {
			n2 = s2[i]
		}
		if n1 > n2 {
			return 1
		}
		if n1 < n2 {
			return -1
		}
	}
	return 0
}

// splitVersion 将 "1.2.3" 拆分为 [1, 2, 3]
func splitVersion(v string) []int {
	var result []int
	current := 0
	for _, ch := range v {
		if ch == '.' {
			result = append(result, current)
			current = 0
		} else if ch >= '0' && ch <= '9' {
			current = current*10 + int(ch-'0')
		}
	}
	result = append(result, current)
	return result
}
