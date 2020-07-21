package iputil

import (
	"encoding/binary"
	"math/big"
	"net"

	"strings"
	"errors"
	"keywea.com/cloud/pblib/pbconverter"
)

func isIPv4(ip net.IP) bool {
	return ip != nil && ip.To4() != nil
}

func IsIPv4(str string) bool {
	addr := net.ParseIP(str)
	return isIPv4(addr)
}

func IsIPv6(str string) bool {
	addr := net.ParseIP(str)
	return addr != nil && !isIPv4(addr) && addr.To16() != nil
}

func IsCIDR(s string) bool {
	_, _, err := net.ParseCIDR(s)
	return err == nil
}

func IsMAC(str string) bool {
	_, err := net.ParseMAC(str)
	return err == nil
}

// checks if a net.IP is in a list of net.IPNet
func IpInMasks(ip net.IP, masks []net.IPNet) bool {
	for _, mask := range masks {
		if mask.Contains(ip) {
			return true
		}
	}
	return false
}

// converts a list of subnets' string to a list of net.IPNet.
func ToMasks(lancidrs []string) (masks []net.IPNet, err error) {
	for _, cidr := range lancidrs {
		var cidrnet *net.IPNet
		_, cidrnet, err = net.ParseCIDR(cidr)
		if err != nil {
			return
		}
		masks = append(masks, *cidrnet)
	}
	return
}

// IsPublicIP returns true if the given IP can be routed on the Internet.
func IsPublicIP(ip net.IP, privateMasks []net.IPNet) bool {
	if !ip.IsGlobalUnicast() {
		return false
	}
	return !IpInMasks(ip, privateMasks)
}

// Convert an ipv4 to a uint32
func Ipv4ToInt(ip net.IP) (uint32, error) {
	ip = ip.To4()
	if ip == nil {
		return 0, errors.New("not able to convert ip to ipv4")
	}

	return binary.BigEndian.Uint32([]byte(ip)), nil
}

//convert a uint32 to an ipv4
func IntToIPv4(ip uint32) net.IP {
	addr := net.IP{0, 0, 0, 0}
	binary.BigEndian.PutUint32(addr, ip)
	return addr
}

func IPv6ToBigInt(ip net.IP) *big.Int {
	IPv6Int := big.NewInt(0)
	IPv6Int.SetBytes(ip)

	return IPv6Int
}

func IPv6ToHLint(ip net.IP) (uint64, uint64) {
	b := IPv6ToBigInt(ip).Bytes()
	high := binary.BigEndian.Uint64(b[:8])
	low := binary.BigEndian.Uint64(b[8:])

	return high, low
}

func BigIntToIPv6(b *big.Int) net.IP {
	return (net.IP)(b.Bytes())
}

func HLintToIPv6(high uint64, low uint64) net.IP {
	hbyte := pbconverter.Int64ToBytes(int64(high))
	lbyte := pbconverter.Int64ToBytes(int64(low))
	b := append(hbyte, lbyte...)
	return (net.IP)(b)
}

// Return the final address of a net range. Convert to IPv4 if possible,
// otherwise return an ipv6
func LastAddress(n *net.IPNet) net.IP {
	ip := n.IP.To4()
	if ip == nil {
		ip = n.IP
		return net.IP{
			ip[0] | ^n.Mask[0], ip[1] | ^n.Mask[1], ip[2] | ^n.Mask[2],
			ip[3] | ^n.Mask[3], ip[4] | ^n.Mask[4], ip[5] | ^n.Mask[5],
			ip[6] | ^n.Mask[6], ip[7] | ^n.Mask[7], ip[8] | ^n.Mask[8],
			ip[9] | ^n.Mask[9], ip[10] | ^n.Mask[10], ip[11] | ^n.Mask[11],
			ip[12] | ^n.Mask[12], ip[13] | ^n.Mask[13], ip[14] | ^n.Mask[14],
			ip[15] | ^n.Mask[15]}
	}

	return net.IPv4(
		ip[0] | ^n.Mask[0],
		ip[1] | ^n.Mask[1],
		ip[2] | ^n.Mask[2],
		ip[3] | ^n.Mask[3])
}

// GetHostIPs returns a list of IP addresses of all host's interfaces.
func GetHostIPs() ([]net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var ips []net.IP
	for _, iface := range ifaces {
		if strings.HasPrefix(iface.Name, "docker") || (iface.Flags&net.FlagUp == 0) {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				ips = append(ips, ipnet.IP)
			}
		}
	}

	return ips, nil
}
