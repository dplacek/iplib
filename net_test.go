package iplib

import (
	"net"
	"testing"
)

var NewNetTests = []struct {
	ip      net.IP
	masklen int
	out     string
}{
	{
		net.IP{192, 168, 0, 0},
		32,
		"192.168.0.0/32",
	},
	{
		net.IP{192, 168, 0, 0},
		24,
		"192.168.0.0/24",
	},
	{
		net.IP{192, 168, 0, 7},
		32,
		"192.168.0.7/32",
	},
	{
		net.IP{192, 168, 0, 7},
		24,
		"192.168.0.0/24",
	},
	{
		net.IP{32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		64,
		"2001:db8::/64",
	},
}

func TestNewNet(t *testing.T) {
	for _, tt := range NewNetTests {
		xnet, _ := NewNet(tt.ip, tt.masklen)
		_, pnet, _ := net.ParseCIDR(tt.out)
		if xnet.String() != pnet.String() {
			t.Errorf("On NewNet(%s, %d) expected %s got %s", tt.ip.String(), tt.masklen, pnet.String(), xnet.String())
		}
	}
}

var NewNetBetweenTests = []struct {
	start net.IP
	end   net.IP
	xnet  string
	exact bool
	err   error
}{
	{
		net.IP{192, 168, 0, 255},
		net.IP{192, 168, 2, 0},
		"192.168.1.0/24",
		true,
		nil,
	},
	{
		net.IP{192, 168, 0, 255},
		net.IP{10, 0, 0, 0},
		"",
		false,
		ErrNoValidRange,
	},
	{
		net.IP{192, 168, 0, 255},
		net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52},
		"",
		false,
		ErrNoValidRange,
	},
	{
		net.IP{192, 168, 0, 255},
		net.IP{192, 168, 0, 255},
		"",
		false,
		ErrNoValidRange,
	},
	{
		net.IP{192, 168, 0, 255},
		net.IP{192, 168, 1, 1},
		"192.168.1.0/32",
		true,
		nil,
	},
	{
		net.IP{192, 168, 1, 0},
		net.IP{192, 168, 1, 2},
		"192.168.1.1/32",
		true,
		nil,
	},
	{
		net.IP{192, 168, 0, 255},
		net.IP{192, 168, 1, 2},
		"192.168.1.0/31",
		true,
		nil,
	},
	{
		net.IP{192, 168, 0, 255},
		net.IP{192, 168, 1, 3},
		"192.168.1.0/31",
		false,
		nil,
	},
	{
		net.IP{192, 168, 1, 0},
		net.IP{192, 168, 1, 3},
		"192.168.1.1/32",
		false,
		nil,
	},
	{
		net.IP{192, 168, 0, 255},
		net.IP{192, 168, 1, 4},
		"192.168.1.0/30",
		true,
		nil,
	},
	{
		net.IP{192, 168, 0, 255},
		net.IP{192, 168, 1, 5},
		"192.168.1.0/30",
		false,
		nil,
	},
	{
		net.IP{192, 168, 0, 254},
		net.IP{192, 168, 2, 0},
		"192.168.0.255/32",
		false,
		nil,
	},
	{
		net.IP{192, 168, 0, 255},
		net.IP{192, 168, 2, 0},
		"192.168.1.0/24",
		true,
		nil,
	},
	{
		net.IP{32, 1, 13, 183, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		net.IP{32, 1, 13, 184, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		"2001:db8::/64",
		true,
		nil,
	},
}

func TestNewNetBetween(t *testing.T) {
	for _, tt := range NewNetBetweenTests {
		xnet, exact, err := NewNetBetween(tt.start, tt.end)
		if tt.err != nil {
			if tt.err != err {
				t.Errorf("On NewNetBetween(%s, %s) expected error '%v', got '%v'", tt.start, tt.end, tt.err, err)
			}
		} else {
			if xnet.String() != tt.xnet {
				t.Errorf("On NewNetBetween(%s, %s) expected '%s', got '%s'", tt.start, tt.end, tt.xnet, xnet.String())
			}
			if exact != tt.exact {
				t.Errorf("On NewNetBetween(%s, %s) expected '%t', got '%t'", tt.start, tt.end, tt.exact, exact)
			}
		}
	}
}

// ParseCIDR wraps net.ParseCIDR so it's redundant to test it except to make sure the wildcard is correct
func TestParseCIDR(t *testing.T) {
	for _, tt := range Network4Tests {
		if tt.version == 6 {
			continue
		}
		_, ipn, _ := ParseCIDR(tt.inaddrStr)

		if ipn.LastAddress().String() != tt.wildcard.String() {
			t.Errorf("On %s got Network.Wildcard == %v, want %v", tt.inaddrStr, ipn.LastAddress(), tt.wildcard)
		}
	}
}
