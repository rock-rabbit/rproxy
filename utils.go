package rproxy

import (
	"bytes"
	"os/exec"
)

// GetCommandStdout 获取命令行输出
func GetCommandStdout(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	buf := bytes.NewBuffer([]byte{})
	cmd.Stdout = buf
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return buf.String(), err
}
