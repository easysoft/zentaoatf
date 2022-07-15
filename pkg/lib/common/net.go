package commonUtils

import (
	"net"
)

func GetIp() (ip net.IP, mac net.HardwareAddr) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return
		}
		for _, addr := range addrs {
			ip = getIpFromAddr(addr)
			if ip == nil {
				continue
			}

			mac = iface.HardwareAddr
			return
		}
	}
	return
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

func GetValidPort(start int, end int, existPortMap *map[int]bool) (ret int) {
	newPort := 0

	for i := 0; i < 99; i++ {
		port := start + i
		if port > end {
			break
		}

		if _, ok := (*existPortMap)[port]; !ok {
			newPort = port
			break
		}
	}

	if newPort > 0 {
		ret = newPort
		(*existPortMap)[newPort] = true
	}
	return
}
