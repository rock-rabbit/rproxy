package rproxy

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var (
	Networkservice string
)

func init() {
	Networkservice = GetNetworkservice()
}

func OpenURI(s string) {
	exec.Command("open", s).Start()
}

// GetNetworkservice 获取当前网络名称
func GetNetworkservice() (s string) {
	s, _ = GetCommandStdout("bash", "-c", "networksetup -listallnetworkservices | head -n 3 | tail -n 1")
	s = strings.Replace(s, "\n", "", 1)
	return
}

// GetProxyEnable 获取代理是否开启
func GetProxyEnable() bool {
	cmdStr := fmt.Sprintf(`networksetup -getwebproxy "%s" | head -n 1`, Networkservice)
	status, _ := GetCommandStdout("bash", "-c", cmdStr)
	return strings.Contains(status, "Yes")
}

// GetProxy 获取代理服务器
func GetProxy() (enable bool, server string, err error) {
	cmdStr := fmt.Sprintf(`networksetup -getwebproxy "%s"`, Networkservice)
	data, err := GetCommandStdout("bash", "-c", cmdStr)
	if err != nil {
		return
	}
	enable = strings.Contains(data, "Yes")
	ss := regexp.MustCompile(`(?m)Server: (.*?)\n.*?Port: (\d+)`).FindStringSubmatch(data)
	if len(ss) < 3 {
		err = errors.New("proxy info error")
		return
	}
	server = ss[1] + ":" + ss[2]
	return
}

// SetProxy 设置代理
func SetProxy(enable bool, server string) error {
	var err error
	if !enable {
		cmdStr1 := fmt.Sprintf(`networksetup -setwebproxystate "%s" off`, Networkservice)
		_, err = GetCommandStdout("bash", "-c", cmdStr1)
		if err != nil {
			return err
		}

		cmdStr2 := fmt.Sprintf(`networksetup -setsecurewebproxystate "%s" off`, Networkservice)
		_, err = GetCommandStdout("bash", "-c", cmdStr2)
		if err != nil {
			return err
		}

		cmdStr3 := fmt.Sprintf(`networksetup -setsocksfirewallproxystate "%s" off`, Networkservice)
		_, err = GetCommandStdout("bash", "-c", cmdStr3)
		return err
	}
	addr := strings.Split(server, ":")
	if len(addr) < 2 {
		return errors.New("addr error")
	}
	cmdStr1 := fmt.Sprintf(`networksetup -setwebproxy "%s" %s %s`, Networkservice, addr[0], addr[1])
	cmdStr2 := fmt.Sprintf(`networksetup -setsecurewebproxy "%s" %s %s`, Networkservice, addr[0], addr[1])
	cmdStr3 := fmt.Sprintf(`networksetup -setsocksfirewallproxy "%s" %s %s`, Networkservice, addr[0], addr[1])
	_, err = GetCommandStdout("bash", "-c", cmdStr1)
	if err != nil {
		return err
	}
	_, err = GetCommandStdout("bash", "-c", cmdStr2)
	if err != nil {
		return err
	}
	_, err = GetCommandStdout("bash", "-c", cmdStr3)
	return err
}
