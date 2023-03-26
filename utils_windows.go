package rproxy

import (
	"errors"

	"golang.org/x/sys/windows/registry"
)

var (
	// ErrProxyNotEnable 代理未开启错误
	ErrProxyNotEnable = errors.New("proxy not enable")
	// RegPath 注册表路径
	RegPath = `SOFTWARE\Microsoft\Windows\CurrentVersion\Internet Settings`
)

// GetProxyEnable 获取代理是否开启
func GetProxyEnable() error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, RegPath, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer key.Close()
	val, _, err := key.GetIntegerValue("ProxyEnable")
	if err != nil {
		return err
	}
	if val != 0 {
		return nil
	}
	return ErrProxyNotEnable
}

// SetProxy 设置代理
func SetProxy(enable bool, server string) error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, RegPath, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer key.Close()
	if enable {
		err = key.SetDWordValue("ProxyEnable", uint32(1))
	} else {
		err = key.SetDWordValue("ProxyEnable", uint32(0))
	}
	if err != nil {
		return err
	}
	// 开启时设置代理服务器
	if enable {
		err = key.SetStringValue("ProxyServer", server)
		if err != nil {
			return err
		}
	}
	return nil
}
