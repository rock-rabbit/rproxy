package rproxy

import (
	"fmt"
	"os/exec"
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

// SetProxy 设置代理
func SetProxy(enable bool, server string) bool {
	var err error
	if !enable {
		cmdStr1 := fmt.Sprintf(`networksetup -setwebproxystate "%s" off`, Networkservice)
		_, err = GetCommandStdout("bash", "-c", cmdStr1)
		if err != nil {
			return false
		}
		cmdStr2 := fmt.Sprintf(`networksetup -setsecurewebproxystate "%s" off`, Networkservice)
		_, err = GetCommandStdout("bash", "-c", cmdStr2)
		return err == nil
	}
	addr := strings.Split(server, ":")
	if len(addr) < 2 {
		return false
	}
	cmdStr1 := fmt.Sprintf(`networksetup -setwebproxy "%s" %s %s`, Networkservice, addr[0], addr[1])
	cmdStr2 := fmt.Sprintf(`networksetup -setsecurewebproxy "%s" %s %s`, Networkservice, addr[0], addr[1])
	_, err = GetCommandStdout("bash", "-c", cmdStr1)
	if err != nil {
		return false
	}
	_, err = GetCommandStdout("bash", "-c", cmdStr2)
	return err == nil
}
