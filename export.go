package rproxy

import "github.com/ouqiang/goproxy"

var std = NewRproxy()

func NewRproxy() *Rproxy {
	return &Rproxy{
		middles: make(map[string]Middle),
	}
}

// RegisterMiddle 注册中间件
func RegisterMiddle(middle Middle) error {
	return std.RegisterMiddle(middle)
}

// Run 运行代理服务
func Run(addr string) error {
	return std.Run(addr)
}

// NewGoproxy 创建代理服务
func NewGoproxy() *goproxy.Proxy {
	return std.NewGoproxy()
}
