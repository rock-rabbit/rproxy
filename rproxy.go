package rproxy

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/elazarl/goproxy.v1"
)

type Rproxy struct {
	// dataMiddles 代理中间件
	dataMiddles map[string]DataMiddle
}

// RproxyService 代理服务
type RproxyService struct {
	Srv *http.Server
}

var (
	CrtFile = filepath.Join(GetAppDatadir(), "rproxy-ca-cert.crt")
	KeyFile = filepath.Join(GetAppDatadir(), "rproxy-ca-cert.key")
)

// NewGoproxy 创建代理服务
func (e *Rproxy) NewGoproxy() (*goproxy.ProxyHttpServer, error) {
	// 初始化证书
	caCert, caKey, err := LoadRootCA()
	if err != nil {
		return nil, err
	}
	err = SetCA(caCert, caKey)
	if err != nil {
		return nil, err
	}
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	goproxy.UrlIs()

	// 设置 Middle
	for _, m := range e.dataMiddles {
		proxy.OnResponse(goproxy.ReqConditionFunc(m.Scope)).DoFunc(m.Handle)
	}

	return proxy, nil
}

// SetCA 设置根证书
func SetCA(crt []byte, key []byte) error {
	var err error
	caCert, caKey, err := LoadRootCA()
	if err != nil {
		return err
	}
	goproxyCa, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return err
	}
	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
		return err
	}
	goproxy.GoproxyCa = goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	return nil
}

// LoadRootCA 加载根证书
func LoadRootCA() (crt []byte, key []byte, err error) {
	if !FileExists(CrtFile) || !FileExists(KeyFile) {
		crt, key, err = GenerateCA()
		if err != nil {
			return
		}
		// 创建文件夹
		if err = os.MkdirAll(filepath.Dir(CrtFile), 0755); err != nil {
			return
		}
		if err = os.WriteFile(CrtFile, crt, 0644); err != nil {
			return
		}
		if err = os.WriteFile(KeyFile, key, 0644); err != nil {
			return
		}
		return
	}
	// 加载根证书
	crt, err = os.ReadFile(CrtFile)
	if err != nil {
		return
	}
	key, err = os.ReadFile(KeyFile)
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
	proxy, err := e.NewGoproxy()
	if err != nil {
		return nil, err
	}
	srv := &http.Server{
		Addr:         ln.Addr().String(),
		Handler:      proxy,
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

// RegisterDataMiddle 注册数据中间件
func (e *Rproxy) RegisterDataMiddle(middles ...DataMiddle) error {
	for _, middle := range middles {
		if e.dataMiddles == nil {
			e.dataMiddles = make(map[string]DataMiddle)
		}
		e.dataMiddles[middle.Name()] = middle
	}
	return nil
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
