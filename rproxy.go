package rproxy

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rock-rabbit/rproxy/goproxy"
)

type Rproxy struct {
	// middles 代理中间件
	middles map[string]Middle
}

// RproxyService 代理服务
type RproxyService struct {
	Srv *http.Server
}

var (
	caFile  = filepath.Join(GetAppDatadir(), "rproxy-ca-cert.crt")
	keyFile = filepath.Join(GetAppDatadir(), "rproxy-ca-cert.key")
)

// NewGoproxy 创建代理服务
func (e *Rproxy) NewGoproxy() *goproxy.Proxy {
	ca, key, _ := e.LoadRootCA()
	return goproxy.New(
		goproxy.WithDelegate(e),
		goproxy.WithDecryptHTTPS(NewCretCache()),
		goproxy.WithRootCA(ca, key),
	)
}

// LoadRootCA 加载根证书
func (e *Rproxy) LoadRootCA() (ca []byte, key []byte, err error) {
	if !FileExists(caFile) || !FileExists(keyFile) {
		ca, key, err = GenerateCA()
		if err != nil {
			return
		}
		// 创建文件夹
		if err = os.MkdirAll(filepath.Dir(caFile), 0755); err != nil {
			return
		}
		if err = os.WriteFile(caFile, ca, 0644); err != nil {
			return
		}
		if err = os.WriteFile(keyFile, key, 0644); err != nil {
			return
		}
		return
	}
	// 加载根证书
	ca, err = os.ReadFile(caFile)
	if err != nil {
		return
	}
	key, err = os.ReadFile(keyFile)
	if err != nil {
		return
	}
	return
}

// Run 运行代理服务
func (e *Rproxy) Run(addr string) (rsrv *RproxyService, err error) {
	if addr == "" {
		addr = ":"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	srv := &http.Server{
		Addr:         ln.Addr().String(),
		Handler:      e.NewGoproxy(),
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}
	go func() {
		err := srv.Serve(ln)
		if err != nil {
			fmt.Println("rproxy run error:", err)
		}
	}()
	return &RproxyService{
		Srv: srv,
	}, nil
}

// Addr 获取代理服务地址
func (s *RproxyService) Addr() string {
	var (
		host     = s.Srv.Addr
		hostList = strings.Split(s.Srv.Addr, ":")
	)
	if len(hostList) != 0 {
		host = hostList[len(hostList)-1]
	}
	return "127.0.0.1:" + host
}

// Enable 开启代理服务
func (s *RproxyService) Enable() error {
	return SetProxy(true, s.Addr())
}

// Disable 关闭代理服务
func (s *RproxyService) Disable() error {
	return SetProxy(false, "")
}

// IsEnable 是否开启了当前代理服务
func (s *RproxyService) IsEnable() error {
	enable, server, err := GetProxy()
	if err != nil {
		return err
	}
	if !enable {
		return fmt.Errorf("proxy is disable")
	}
	if server != s.Addr() {
		return fmt.Errorf("proxy server is %s", server)
	}
	return nil
}

// RegisterMiddle 注册中间件
func (e *Rproxy) RegisterMiddle(middles ...Middle) error {
	for _, middle := range middles {
		if e.middles == nil {
			e.middles = make(map[string]Middle)
		}
		e.middles[middle.Name()] = middle
	}
	return nil
}

// Connect 连接前执行
func (e *Rproxy) Connect(ctx *goproxy.Context, rw http.ResponseWriter) {}

// Auth 权限验证
func (e *Rproxy) Auth(ctx *goproxy.Context, rw http.ResponseWriter) {}

// BeforeRequest 请求开始前执行
func (e *Rproxy) BeforeRequest(ctx *goproxy.Context) {
}

// BeforeResponse 请求结束后执行
func (e *Rproxy) BeforeResponse(ctx *goproxy.Context, resp *http.Response, err error) {
	if err != nil {
		return
	}
	var (
		isRead   = false
		body     []byte
		bodyData []byte
	)

	// 执行中间件
	for _, middle := range e.middles {
		if middle.Scope(ctx.Req) {
			// 触发适用范围时，读取 body 数据
			if !isRead {
				body, _ = io.ReadAll(resp.Body)
				var (
					newBody = bytes.NewReader(body)
					reader  io.ReadCloser
				)
				// 解压缩编码
				switch resp.Header.Get("Content-Encoding") {
				case "gzip":
					reader, _ = gzip.NewReader(newBody)
				case "deflate":
					reader = flate.NewReader(newBody)
				default:
					reader = resp.Body
				}
				bodyData, _ = io.ReadAll(reader)
				isRead = true
			}
			middle.Handle(resp, bodyData)
		}
	}

	if isRead {
		// 重新设置为原始的 body 给后续使用
		resp.Body = io.NopCloser(bytes.NewReader(body))
	}
}

// ParentProxy 设置上级代理
func (e *Rproxy) ParentProxy(req *http.Request) (*url.URL, error) {
	return nil, nil
}

// Finish 请求结束后执行
func (e *Rproxy) Finish(ctx *goproxy.Context) {}

// ErrorLog 记录错误日志
func (e *Rproxy) ErrorLog(err error) {}

// WebSocketSendMessage websocket 发送消息
func (h *Rproxy) WebSocketSendMessage(ctx *goproxy.Context, messageType *int, payload *[]byte) {}

// WebSockerReceiveMessage websocket 接收消息
func (h *Rproxy) WebSocketReceiveMessage(ctx *goproxy.Context, messageType *int, payload *[]byte) {}
