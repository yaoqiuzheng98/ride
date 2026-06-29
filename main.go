package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"ride/db/cache"
	"ride/db/table"
	"ride/handler"
	"ride/middlerware"
	"ride/pkg/path"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// 公开接口：创建用户
	r.POST("/user", handler.CreateUser)
	// 需要 APP 版本校验 + X-User-Id 认证的接口
	auth := r.Group("")
	auth.Use(middlerware.AppVersionCheck())
	auth.Use(middlerware.UserIDAuth())
	auth.GET("/point", handler.Point)
	r.GET("/download", func(c *gin.Context) {
		currentPath, err := path.GetCurrentPath()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		filePath := fmt.Sprintf("%s/ride.apk", currentPath)
		file, err := os.Open(filePath)
		if err != nil {
			c.String(http.StatusNotFound, "File not found")
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal server error")
			return
		}

		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileInfo.Name()))
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

		_, err = io.Copy(c.Writer, file)
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal server error")
			return
		}
	})
	cache.NewMemory([]cache.Cache{table.GetPoints()}).Refresh()
	err := r.Run(":9999")
	if err != nil {
		panic(err)
	}
}
