package utils

import (
	"log"
	"net"
	"strings"
)

// 获取本机ip
func LocalIpv4() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		address, err := iface.Addrs()
		if err != nil {
			panic(err)
		}
		for _, addr := range address {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String()
		}
	}
	panic("not find available ip")
}

// 获取公网ip
func PublicNetworkIp() string {
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := conn.Close(); e != nil {
			log.Println(err)
		}
	}()
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}
