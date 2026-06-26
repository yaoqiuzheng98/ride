package path

import (
	"fmt"
	"path"
	"runtime"
)

func GetCurrentPath() (string, error) {
	if _, file, _, ok := runtime.Caller(1); ok {
		return path.Dir(file), nil
	}
	return "", fmt.Errorf("获取路径失败")
}
