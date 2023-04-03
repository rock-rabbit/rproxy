package rproxy

import (
	"crypto/tls"
	"sync"
)

type CretCache struct {
	m sync.Map
}

func NewCretCache() *CretCache {
	return &CretCache{}
}

// Set 设置证书
func (c *CretCache) Set(host string, cert *tls.Certificate) {
	c.m.Store(host, cert)
}

// Get 获取证书
func (c *CretCache) Get(host string) *tls.Certificate {
	v, ok := c.m.Load(host)
	if !ok {
		return nil
	}
	return v.(*tls.Certificate)
}
