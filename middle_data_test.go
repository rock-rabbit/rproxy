package rproxy_test

import (
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/rock-rabbit/rproxy"
)

func TestMiniMiddle(t *testing.T) {
	rproxy.RegisterDataMiddle(rproxy.NewMiniDataMiddle("GET", "sspai.com", "/api/v1/combo/recommend/page/get", func(res *http.Response, body []byte) {
		t.Log(res.Request.URL.String())
		t.Log(string(body))
		// 重置res body
		if v, ok := res.Body.(io.ReadSeeker); ok {
			v.Seek(0, 0)
		}
		nbody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(nbody))
	}))
	ctl, err := rproxy.Run()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("proxy addr:", ctl.Addr())
	ctl.Enable()

	// 关闭信号时关闭代理
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
	err = ctl.Disable()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("proxy disable")

	os.Exit(0)
}
