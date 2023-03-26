package rproxy_test

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/rock-rabbit/rproxy"
)

type sspaiMiddle struct{}

func (s *sspaiMiddle) Name() string {
	return "sspai"
}

func (s *sspaiMiddle) Scope(req *http.Request) bool {
	// 过滤静态资源
	if !rproxy.FiltreResource(req) {
		return false
	}
	return req.URL.Host == "sspai.com" && req.Method == http.MethodGet
}

func (s *sspaiMiddle) Handle(res *http.Response, body []byte) {
	fmt.Println(res.Request.URL.String())
	// fmt.Println(string(body))
}

var _ rproxy.Middle = &sspaiMiddle{}

func TestMiddle(t *testing.T) {
	rproxy.RegisterMiddle(&sspaiMiddle{})
	rproxy.Run(":8080")
}

func TestMiniMiddle(t *testing.T) {
	rproxy.RegisterMiddle(rproxy.NewMiniMiddle("GET", "sspai.com", "/api/v1/recommend/page/get", func(res *http.Response, body []byte) {
		fmt.Println(res.Request.URL.String())
	}))
	ctl, err := rproxy.Run("")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("proxy addr:", ctl.Addr())
	ctl.Enable()

	// 关闭信号时关闭代理
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
	err = ctl.Disable()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("proxy disable")

	os.Exit(0)
}
