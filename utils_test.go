package rproxy_test

import (
	"fmt"
	"testing"

	"github.com/rock-rabbit/rproxy"
)

func TestGetProxy(t *testing.T) {
	fmt.Println(rproxy.GetProxy())
}

func TestOpenProxy(t *testing.T) {
	fmt.Println(rproxy.SetProxy(true, "127.0.0.1:8877"))
}

func TestStopProxy(t *testing.T) {
	fmt.Println(rproxy.SetProxy(false, "127.0.0.1:8877"))
}
