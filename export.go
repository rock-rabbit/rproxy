package rproxy

import (
	"github.com/rock-rabbit/rproxy/goproxy"
)

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
func Run(addrs ...string) (rsrv *RproxyService, err error) {
	addr := ":"
	if len(addrs) != 0 {
		addr = addrs[0]
	}
	return std.Run(addr)
}

// NewGoproxy 创建代理服务
func NewGoproxy() *goproxy.Proxy {
	return std.NewGoproxy()
}
