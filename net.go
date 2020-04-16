package iplib

import (
	"net"
	"strings"
)

type Net interface {
	Contains(ip net.IP) bool
	ContainsNet(network Net) bool
	FirstAddress() net.IP
	IP() net.IP
	LastAddress() net.IP
	Mask() net.IPMask
	String() string
	Version() int
}

// NewNet returns a new Net object containing ip at the specified masklen. In
// the Net6 case the hostbits value will be set to 0. If the masklen is set
// to an insane value (greater than 32 for IPv4 or 128 for IPv6) an empty Net
// will be returned
func NewNet(ip net.IP, masklen int) Net {
	version := Version(ip)
	if version == 6 {
		return NewNet6(ip, masklen, 0)
	}
	return NewNet4(ip, masklen)
}

// NewNetBetween takes two net.IP's as input and will return the largest
// netblock that can fit between them (exclusive of the IP's themselves).
// If there is an exact fit it will set a boolean to true, otherwise the bool
// will be false. If no fit can be found (probably because a >= b) an
// ErrNoValidRange will be returned.
func NewNetBetween(a, b net.IP) (Net, bool, error) {
	var exact = false
	v := CompareIPs(a, b)
	if v != -1 {
		return nil, exact, ErrNoValidRange
	}

	if Version(a) != Version(b) {
		return nil, exact, ErrNoValidRange
	}

	maskMax := 128
	if EffectiveVersion(a) == 4 {
		maskMax = 32
	}

	ipa := NextIP(a)
	ipb := PreviousIP(b)
	for i := 1; i <= maskMax; i++ {
		xnet := NewNet(ipa, i)

		va := CompareIPs(xnet.FirstAddress(), ipa)
		vb := CompareIPs(xnet.LastAddress(), ipb)
		if va >= 0 && vb <= 0 {
			if va == 0 && vb == 0 {
				exact = true
			}
			return xnet, exact, nil
		}
	}
	return nil, exact, ErrNoValidRange
}

// ByNet implements sort.Interface for iplib.Net based on the
// starting address of the netblock, with the netmask as a tie breaker. So if
// two Networks are submitted and one is a subset of the other, the enclosing
// network will be returned first.
type ByNet []Net

// Len implements sort.interface Len(), returning the length of the
// ByNetwork array
func (bn ByNet) Len() int {
	return len(bn)
}

// Swap implements sort.interface Swap(), swapping two elements in our array
func (bn ByNet) Swap(a, b int) {
	bn[a], bn[b] = bn[b], bn[a]
}

// Less implements sort.interface Less(), given two elements in the array it
// returns true if the LHS should sort before the RHS. For details on the
// implementation, see CompareNets()
func (bn ByNet) Less(a, b int) bool {
	val := CompareNets(bn[a], bn[b])
	if val == -1 {
		return true
	}
	return false
}

// ParseCIDR returns a new Net object. It is a passthrough to net.ParseCIDR
// and will return any error it generates to the caller. There is one major
// difference between how net.IPNet manages addresses and how ipnet.Net does,
// and this function exposes it: net.ParseCIDR *always* returns an IPv6
// address; if given a v4 address it returns the RFC4291 IPv4-mapped IPv6
// address internally, but treats it like v4 in practice. In contrast
// iplib.ParseCIDR will re-encode it as a v4
func ParseCIDR(s string) (net.IP, Net, error) {
	ip, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return ip, nil, err
	}
	masklen, _ := ipnet.Mask.Size()

	if strings.Contains(s, ".") {
		n := NewNet4(ForceIP4(ip), masklen)
		return ForceIP4(ip), n, err
	}

	n := NewNet6(ip, masklen, 0)
	return ip, n, err
}
