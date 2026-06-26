package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"ride/db/cache"
	"ride/db/collection"
	"ride/handler"
	"ride/middlerware"
	"ride/pkg/path"
)

func main() {
	r := gin.Default()
	//r.Use(middlerware.Auth())
	r.GET("/point", middlerware.Auth(), handler.Point)
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
	cache.NewMemory([]cache.Cache{collection.GetPoints()}).Refresh()
	err := r.Run(":9999")
	if err != nil {
		panic(err)
	}
}
