package rproxy

var std = NewRproxy()

func NewRproxy() *Rproxy {
	return &Rproxy{
		dataMiddles: make(map[string]DataMiddle),
	}
}

// RegisterDataMiddle 注册数据中间件
func RegisterDataMiddle(middles ...DataMiddle) error {
	return std.RegisterDataMiddle(middles...)
}

// Run 运行代理服务
func Run(addrs ...string) (rsrv *RproxyService, err error) {
	addr := ":"
	if len(addrs) != 0 {
		addr = addrs[0]
	}
	return std.Run(addr)
}
