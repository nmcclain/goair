package commands

import (
	"fmt"

	"net"
)

//GetIPRange returns a list of IPs
func GetIPRange(startAddress, endAddress string) ([]string, error) {
	//ip, ipnet, err := net.ParseCIDR("62.76.47.12/28")
	startIP, _, err := net.ParseCIDR(startAddress + "/32")
	if err != nil {
		return []string{}, fmt.Errorf("problem parsing IP: %v", err)
	}
	endIP, _, err := net.ParseCIDR(endAddress + "/32")
	if err != nil {
		return []string{}, fmt.Errorf("problem parsing IP: %v", err)
	}

	ipRange := make([]string, 0)
	ip := net.IP{}
	for ip = startIP; string(ip) != string(endIP); incrementIP(ip) {
		ipRange = append(ipRange, ip.String())
	}
	ipRange = append(ipRange, ip.String())

	return ipRange, nil
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
