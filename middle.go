package rproxy

import (
	"net/http"
	"regexp"
	"strings"
)

// Middle 中间代理接口
type Middle interface {
	// Name 中间件名称
	Name() string
	// Scope 指定代理范围
	Scope(req *http.Request) bool
	// Handle 处理请求
	Handle(res *http.Response, body []byte)
}

// MiniMiddle 实现微型中间件
type MiniMiddle struct {
	Method  string
	Host    string
	Path    string
	HandleF func(res *http.Response, body []byte)
}

var _ Middle = &MiniMiddle{}

// NewMiniMiddle 实例化微型中间件
func NewMiniMiddle(method, host, path string, handleF func(res *http.Response, body []byte)) *MiniMiddle {
	return &MiniMiddle{
		Method:  method,
		Host:    host,
		Path:    path,
		HandleF: handleF,
	}
}

// Name 中间件名称
func (m *MiniMiddle) Name() string {
	return "mini_" + m.Host + m.Path
}

// Scope 指定代理范围
func (m *MiniMiddle) Scope(req *http.Request) bool {
	if m.Method != "" {
		if req.Method != strings.ToUpper(m.Method) {
			return false
		}
	}
	if m.Host != "" {
		matched, _ := regexp.MatchString(m.Host, req.URL.Host)
		if !matched {
			return false
		}
	}
	if m.Path != "" {
		matched, _ := regexp.MatchString(m.Path, req.URL.Path)
		if !matched {
			return false
		}
	}
	return true
}

// Handle 处理请求
func (m *MiniMiddle) Handle(res *http.Response, body []byte) {
	m.HandleF(res, body)
}
