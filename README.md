### 网络代理服务

### 使用方式
``` go
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
	fmt.Println(string(body))
}

var _ rproxy.Middle = &sspaiMiddle{}

func TestMiddle(t *testing.T) {
	rproxy.RegisterMiddle(&sspaiMiddle{})
	rproxy.Run(":8080")
}
```