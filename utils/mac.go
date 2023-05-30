package utils

import (
	"net"
)

func GetMac() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, inter := range interfaces {
		if mac := inter.HardwareAddr.String(); mac != "" {
			return mac
		}
	}
	return ""
}
