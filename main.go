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
	ride := r.Group("/ride")
	// 静态文件（图标等）
	ride.Static("/static", "./web/static")
	// 公开接口：健康检查、下载页、创建用户
	ride.GET("/heartbeat", handler.Heartbeat)
	ride.GET("/download-page", handler.DownloadPage)
	ride.POST("/user", handler.CreateUser)
	// 需要 APP 版本校验 + X-User-Id 认证的接口
	auth := ride.Group("")
	auth.Use(middlerware.AppVersionCheck())
	auth.Use(middlerware.UserIDAuth())
	auth.GET("/point", handler.Point)
	ride.GET("/download", func(c *gin.Context) {
		currentPath, err := path.GetCurrentPath()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		filePath, _ := handler.FindApk(currentPath)
		if filePath == "" {
			c.String(http.StatusNotFound, "File not found")
			return
		}
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

		// 用文件修改时间戳作为版本标识，防止客户端缓存旧 apk
		version := fileInfo.ModTime().Unix()
		c.Header("Content-Disposition", "attachment; filename=骑行日记.apk")
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
		// 禁止缓存，确保每次都重新下载最新 apk
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Header("ETag", fmt.Sprintf("%d-%d", version, fileInfo.Size()))

		_, err = io.Copy(c.Writer, file)
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal server error")
			return
		}
	})
	cache.NewMemory([]cache.Cache{table.GetPoints()}).Refresh()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	err := r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
