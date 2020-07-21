package iputil

import (
	"log"
	"net"
	"testing"
	"keywea.com/cloud/pblib/pbconverter"
)

func assertEqual(t *testing.T, expected interface{}, got interface{}) {
	if got != expected {
		log.Printf("\nExpected %#v\nbut got  %#v", expected, got)
		t.FailNow()
	}
}

func TestIpv4ToInt(t *testing.T) {
	tests := map[string]uint32{
		"0.0.0.0":   0,
		"0.0.0.255": 255,
		"0.0.1.0":   256,
		"0.0.1.100": 356,
		"0.1.0.0":   65536,
		"0.1.1.1":   65793,
		"1.0.0.0":   16777216,
	}

	for test, expected := range tests {
		i, _ := Ipv4ToInt(net.ParseIP(test))
		assertEqual(t, i, expected)
	}
}

func TestIpToUint32Errors(t *testing.T) {
	_, err := Ipv4ToInt(net.ParseIP("2001:0db8:85a3:0042:1000:8a2e:0370:7334"))
	if err == nil {
		t.FailNow()
	}
}

func TestUint32ToIp(t *testing.T) {
	tests := map[uint32]string{
		0:        "0.0.0.0",
		255:      "0.0.0.255",
		256:      "0.0.1.0",
		356:      "0.0.1.100",
		65536:    "0.1.0.0",
		65793:    "0.1.1.1",
		16777216: "1.0.0.0",
	}
	for test, expected := range tests {
		ip := IntToIPv4(test).String()
		assertEqual(t, ip, expected)
	}
}

func TestIPv6ToBigInt(t *testing.T) {
	tests := map[string]string{
		"2001:470:0:76::2": "42540578165168461141553663388954918914",
	}

	for test, expected := range tests {
		ret := IPv6ToBigInt(net.ParseIP(test))
		assertEqual(t, ret.String(), expected)
	}
}

func TestBigIntToIPv6(t *testing.T) {
	tests := map[string]string{
		"42540578165168461141553663388954918914": "2001:470:0:76::2",
	}

	for test, expected := range tests {
		ret := BigIntToIPv6(pbconverter.ToBigInt(test, 10))
		assertEqual(t, ret.String(), expected)
	}
}

type hl map[string]uint64

func TestIPv6ToHLint(t *testing.T) {
	tests := map[string]hl{
		"2001:470:0:76::2": {
			"high": 2306129363273252982,
			"low":  2,
		},
	}

	for test, expected := range tests {
		high, low := IPv6ToHLint(net.ParseIP(test))
		assertEqual(t, high, expected["high"])
		assertEqual(t, low, expected["low"])
	}
}

func TestHLintToIPv6(t *testing.T) {
	high := uint64(2306129363273252982)
	low := uint64(2)
	expected := "2001:470:0:76::2"

	ip := HLintToIPv6(high, low)
	assertEqual(t, ip.String(), expected)
}

func TestLastAddress(t *testing.T) {
	_, netw, _ := net.ParseCIDR("1.2.3.4/24")
	last := LastAddress(netw)
	assertEqual(t, "1.2.3.255", last.String())

	_, net6, _ := net.ParseCIDR("2001:658:22a:cafe::/64")
	last = LastAddress(net6)
	assertEqual(t, "2001:658:22a:cafe:ffff:ffff:ffff:ffff", last.String())
}

func TestGetHostIPs(t *testing.T) {
	ipnet, error := GetHostIPs()
	if error != nil {

	}
	t.Log(ipnet)
}