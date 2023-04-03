package rproxy

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"
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

// GetAppDatadir 获取当前系统的 app 数据目录
func GetAppDatadir() string {
	if runtime.GOOS == "windows" {
		return path.Join(os.Getenv("APPDATA"), AppDatapath)
	} else {
		return path.Join(os.Getenv("HOME"), "/Library/Containers", AppDatapath)
	}
}

// FileExists 文件是否存在
func FileExists(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

// GenerateCA 生成根证书
func GenerateCA() (ca []byte, key []byte, err error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), crand.Reader)
	if err != nil {
		return
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(rand.Int63()),
		Subject: pkix.Name{
			CommonName:   "rockrabbit",
			Country:      []string{"China"},
			Organization: []string{"Rproxy"},
			Province:     []string{"Shandong"},
			Locality:     []string{"Jinan"},
		},
		NotBefore:             time.Now().AddDate(0, -1, 0),
		NotAfter:              time.Now().AddDate(20, 0, 0),
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            2,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		EmailAddresses:        []string{"2896865355@qq.com"},
	}

	derBytes, err := x509.CreateCertificate(crand.Reader, tmpl, tmpl, &priv.PublicKey, priv)
	if err != nil {
		return
	}
	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	}
	ca = pem.EncodeToMemory(certBlock)

	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return
	}
	keyBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privBytes,
	}
	key = pem.EncodeToMemory(keyBlock)
	return
}
