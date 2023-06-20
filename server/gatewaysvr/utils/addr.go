package utils

import (
	"net"
)

// 获取空闲端口
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	defer l.Close()
	if err != nil {
		return 0, err
	}

	return l.Addr().(*net.TCPAddr).Port, nil
}
