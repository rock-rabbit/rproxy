package rproxy

import (
	"net/http"
	"strings"
)

// FiltreResource 过滤静态资源
func FiltreResource(req *http.Request) bool {
	// 过滤列表
	filtre := []string{".png", ".jpg", ".jpeg", ".gif", ".css", ".js"}
	for _, v := range filtre {
		if strings.HasSuffix(req.URL.Path, v) {
			return false
		}
	}
	return true
}
