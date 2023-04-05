package rproxy

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/andybalholm/brotli"
	"gopkg.in/elazarl/goproxy.v1"
)

// DataMiddle 数据中间代理接口
type DataMiddle interface {
	// Name 中间件名称
	Name() string
	// Scope 指定代理范围
	Scope(req *http.Request, ctx *goproxy.ProxyCtx) bool
	// Handle 处理请求
	Handle(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response
}

// MiniMiddle 实现微型中间件
type MiniDataMiddle struct {
	Method  string
	Host    string
	Path    string
	HandleF func(res *http.Response, body []byte)
}

var _ DataMiddle = &MiniDataMiddle{}

// NewMiniMiddle 实例化微型中间件
func NewMiniDataMiddle(method, host, path string, handleF func(res *http.Response, body []byte)) *MiniDataMiddle {
	return &MiniDataMiddle{
		Method:  method,
		Host:    host,
		Path:    path,
		HandleF: handleF,
	}
}

// Name 中间件名称
func (m *MiniDataMiddle) Name() string {
	return "mini_" + m.Host + m.Path
}

// Scope 指定代理范围
func (m *MiniDataMiddle) Scope(req *http.Request, ctx *goproxy.ProxyCtx) bool {
	if m.Method != "" {
		if req.Method != strings.ToUpper(m.Method) {
			return false
		}
	}
	if m.Host != "" && req.Host != m.Host {
		return false
	}
	if m.Path != "" && req.URL.Path != m.Path {
		return false
	}
	return true
}

// Handle 处理请求
func (m *MiniDataMiddle) Handle(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	if m.HandleF != nil {
		var (
			reader     io.Reader
			oldBody, _ = io.ReadAll(resp.Body)
		)
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, _ = gzip.NewReader(bytes.NewReader(oldBody))
		case "deflate":
			reader = flate.NewReader(bytes.NewReader(oldBody))
		case "br":
			reader = brotli.NewReader(bytes.NewReader(oldBody))
		default:
			reader = bytes.NewReader(oldBody)
		}
		body, err := io.ReadAll(reader)
		if err != nil {
			return resp
		}
		m.HandleF(resp, body)
		if v, ok := resp.Body.(io.ReadSeeker); ok {
			v.Seek(0, 0)
		} else {
			resp.Body = io.NopCloser(bytes.NewReader(oldBody))
		}
	}
	return resp
}
