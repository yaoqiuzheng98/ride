package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"ride/pkg/path"
)

// versionRegexp 从 apk 文件名提取版本号，如 骑行日记v1.0.1.apk -> 1.0.1
var versionRegexp = regexp.MustCompile(`v(\d+\.\d+(?:\.\d+)?)\.apk$`)

// FindApk 在指定目录下查找 apk 文件，优先找带版本号的，否则找 骑行日记.apk
func FindApk(dir string) (string, string) {
	// 优先匹配 骑行日记v*.apk
	matches, _ := filepath.Glob(filepath.Join(dir, "骑行日记v*.apk"))
	if len(matches) > 0 {
		apkPath := matches[0]
		name := filepath.Base(apkPath)
		if m := versionRegexp.FindStringSubmatch(name); len(m) > 1 {
			return apkPath, m[1]
		}
		return apkPath, ""
	}
	// 兜底：骑行日记.apk
	apkPath := filepath.Join(dir, "骑行日记.apk")
	if _, err := os.Stat(apkPath); err == nil {
		return apkPath, "1.0"
	}
	return "", ""
}

// DownloadPage 渲染下载页面。
// GET /ride/download-page
func DownloadPage(ctx *gin.Context) {
	dir, err := path.GetCurrentPath()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	apkPath, version := FindApk(dir)
	iconPath := filepath.Join("web", "static", "icon.png")

	type pageData struct {
		Version    string
		SizeText   string
		Build      string
		IconExists bool
	}

	data := pageData{
		Version:    "未知",
		SizeText:   "未知",
		Build:      fmt.Sprintf("%d", 0),
		IconExists: fileExists(iconPath),
	}

	if apkPath != "" {
		if fileInfo, err := os.Stat(apkPath); err == nil {
			if version == "" {
				version = "1.0"
			}
			data.Version = version
			data.Build = fmt.Sprintf("%d", fileInfo.ModTime().Unix())
			data.SizeText = formatSize(fileInfo.Size())
		}
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

// keep strings import referenced
var _ = strings.Contains
