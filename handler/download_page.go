package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"ride/pkg/path"
)

// DownloadPage 渲染下载页面。
// GET /ride/download-page
func DownloadPage(ctx *gin.Context) {
	dir, err := path.GetCurrentPath()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	apkPath := filepath.Join(dir, "骑行日记.apk")
	iconPath := filepath.Join(dir, "icon.png")

	type pageData struct {
		Version    string
		SizeText   string
		Build      string
		IconExists bool
	}

	data := pageData{
		Version:    "1.0",
		Build:      fmt.Sprintf("%d", 0),
		IconExists: fileExists(iconPath),
	}

	if fileInfo, err := os.Stat(apkPath); err == nil {
		data.Version = fmt.Sprintf("%d", fileInfo.ModTime().Unix())
		data.Build = fmt.Sprintf("%d", fileInfo.ModTime().Unix())
		data.SizeText = formatSize(fileInfo.Size())
	} else {
		data.SizeText = "未知"
	}

	tmpl, err := template.ParseFiles("web/templates/download.html")
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(ctx.Writer, data)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func formatSize(bytes int64) string {
	const mb = 1024 * 1024
	if bytes >= mb {
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(mb))
	}
	const kb = 1024
	if bytes >= kb {
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(kb))
	}
	return fmt.Sprintf("%d B", bytes)
}
