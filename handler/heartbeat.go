package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

// Heartbeat 健康检查接口。
// GET /heartbeat
func Heartbeat(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"uptime": time.Since(startTime).String(),
		"time":   time.Now().Format(time.RFC3339),
	})
}
