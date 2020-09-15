package tools

import "net"

// GetIP 获取当前主机的IP
func GetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	thisIP := ""
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				thisIP = ipnet.IP.String()
			}
		}
	}
	return thisIP
}
