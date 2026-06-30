package path

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetCurrentPath 返回可执行文件所在目录。
func GetCurrentPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("获取可执行文件路径失败: %w", err)
	}
	return filepath.Dir(exe), nil
}
