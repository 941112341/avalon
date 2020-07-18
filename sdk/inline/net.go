package inline

import (
	"github.com/pkg/errors"
	"net"
)

func InetAddress() (string, error) {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addresses {
		// check ip address equal with loopback address
		if ip, ok := address.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String(), nil
			}

		}
	}
	return "", errors.New("Can not find the client ip address!")
}

func GetIP() string {
	addr, _ := InetAddress()
	return addr
}
