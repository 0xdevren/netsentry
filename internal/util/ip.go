package util

import (
	"fmt"
	"net"
)

// ParseCIDR parses a CIDR notation string and returns the network and host
// address components.
func ParseCIDR(cidr string) (*net.IPNet, net.IP, error) {
	ip, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, nil, fmt.Errorf("ip: parse CIDR %q: %w", cidr, err)
	}
	return network, ip, nil
}

// NetworkOverlaps reports whether two CIDR networks overlap.
func NetworkOverlaps(a, b string) (bool, error) {
	ipA, netA, err := net.ParseCIDR(a)
	if err != nil {
		return false, fmt.Errorf("ip: parse network %q: %w", a, err)
	}
	ipB, netB, err := net.ParseCIDR(b)
	if err != nil {
		return false, fmt.Errorf("ip: parse network %q: %w", b, err)
	}
	return netA.Contains(ipB) || netB.Contains(ipA), nil
}

// IsPrivateIP reports whether the given IP string is in an RFC-1918 or
// RFC-4193 private range.
func IsPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	private := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"fc00::/7",
	}
	for _, cidr := range private {
		_, block, _ := net.ParseCIDR(cidr)
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// MaskToPrefix converts a dotted-decimal subnet mask to its CIDR prefix length.
// Returns an error if the mask is not a valid contiguous subnet mask.
func MaskToPrefix(mask string) (int, error) {
	ip := net.ParseIP(mask)
	if ip == nil {
		return 0, fmt.Errorf("ip: invalid mask %q", mask)
	}
	ip = ip.To4()
	if ip == nil {
		return 0, fmt.Errorf("ip: mask %q is not IPv4", mask)
	}
	ones, _ := net.IPMask(ip).Size()
	return ones, nil
}
