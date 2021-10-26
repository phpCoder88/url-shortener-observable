package helpers

import (
	"fmt"
	"net"
	"net/http"
)

func GetIP(req *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	return "", fmt.Errorf("no valid ip found")
}
