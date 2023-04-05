### 网络代理服务

### 计划

* 适配 macos 系统
* 自动安装证书

### 使用方式
纯数据中间人
``` go
rproxy.RegisterDataMiddle(rproxy.NewMiniDataMiddle("GET", "sspai.com", "/api/v1/recommend/page/get", func(res *http.Response, body []byte) {
	fmt.Println(res.Request.URL.String())
}))
```

### 证书位置
```
Windows: %APPDATA%\rockrabbit\rproxy\rproxy-ca-cert.crt
Macos:   %HOME%/Library/Containers/rockrabbit/rproxy/rproxy-ca-cert.crt
```


### Windows 安装证书
``` bash
certutil.exe -addstore root rproxy-ca-cert.crt

```

### Macos 安装证书
``` bash
sudo security add-trusted-cert -d -p ssl -p basic -k /Library/Keychains/System.keychain rproxy-ca-cert.crt
```