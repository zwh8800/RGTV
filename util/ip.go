package util

import (
	"fmt"
	"net"
)

func GetIP() string {
	interfaces, _ := net.Interfaces()

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, addr := range addrs {
			// 处理地址类型
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// 排除回环地址
			if ip == nil || ip.IsLoopback() {
				continue
			}

			// 仅处理IPv4地址
			ip = ip.To4()
			if ip == nil {
				continue
			}

			return ip.String()
		}
	}
	return ""
}
