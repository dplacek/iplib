package iplib

import (
	"net"
	"sort"
	"strconv"
	"testing"
)

var Network4Tests = []struct {
	inaddrStr  string
	ipaddr     net.IP
	inaddrMask int
	network    net.IP
	netmask    net.IPMask
	wildcard   net.IPMask
	broadcast  net.IP
	firstaddr  net.IP
	lastaddr   net.IP
	version    int
	count      string // might overflow uint64
}{
	{
		"10.1.2.3/8",
		net.IP{10, 1, 2, 3},
		8,
		net.IP{10, 0, 0, 0},
		net.IPMask{255, 0, 0, 0},
		net.IPMask{0, 255, 255, 255},
		net.IP{10, 255, 255, 255},
		net.IP{10, 0, 0, 1},
		net.IP{10, 255, 255, 254},
		4,
		"16777214",
	},
	{
		"192.168.1.1/23",
		net.IP{192, 168, 1, 1},
		23,
		net.IP{192, 168, 0, 0},
		net.IPMask{255, 255, 254, 0},
		net.IPMask{0, 0, 1, 255},
		net.IP{192, 168, 1, 255},
		net.IP{192, 168, 0, 1},
		net.IP{192, 168, 1, 254},
		4,
		"510",
	},
	{
		"192.168.1.61/26",
		net.IP{192, 168, 1, 61},
		26,
		net.IP{192, 168, 1, 0},
		net.IPMask{255, 255, 255, 192},
		net.IPMask{0, 0, 0, 63},
		net.IP{192, 168, 1, 63},
		net.IP{192, 168, 1, 1},
		net.IP{192, 168, 1, 62},
		4,
		"62",
	},
	{
		"192.168.1.66/26",
		net.IP{192, 168, 1, 66},
		26,
		net.IP{192, 168, 1, 64},
		net.IPMask{255, 255, 255, 192},
		net.IPMask{0, 0, 0, 63},
		net.IP{192, 168, 1, 127},
		net.IP{192, 168, 1, 65},
		net.IP{192, 168, 1, 126},
		4,
		"62",
	},
	{
		"192.168.1.1/30",
		net.IP{192, 168, 1, 1},
		30,
		net.IP{192, 168, 1, 0},
		net.IPMask{255, 255, 255, 252},
		net.IPMask{0, 0, 0, 3},
		net.IP{192, 168, 1, 3},
		net.IP{192, 168, 1, 1},
		net.IP{192, 168, 1, 2},
		4,
		"2",
	},
	{
		"192.168.1.1/31",
		net.IP{192, 168, 1, 1},
		31,
		net.IP{192, 168, 1, 0},
		net.IPMask{255, 255, 255, 254},
		net.IPMask{0, 0, 0, 1},
		net.IP{192, 168, 1, 1},
		net.IP{192, 168, 1, 0},
		net.IP{192, 168, 1, 1},
		4,
		"0",
	},
	{
		"192.168.1.15/32",
		net.IP{192, 168, 1, 15},
		32,
		net.IP{192, 168, 1, 15},
		net.IPMask{255, 255, 255, 255},
		net.IPMask{0, 0, 0, 0},
		net.IP{192, 168, 1, 15},
		net.IP{192, 168, 1, 15},
		net.IP{192, 168, 1, 15},
		4,
		"1",
	},
}

func TestNet4_BroadcastAddress(t *testing.T) {
	for _, tt := range Network4Tests {
		if tt.version == 6 {
			continue
		}
		_, ipn, _ := ParseCIDR(tt.inaddrStr)
		ipn4 := ipn.(Net4)
		if addr := ipn4.BroadcastAddress(); !tt.broadcast.Equal(addr) {
			t.Errorf("On %s got Network.Broadcast == %v, want %v", tt.inaddrStr, addr, tt.broadcast)
		}
	}
}

func TestNet4_Version(t *testing.T) {
	for _, tt := range Network4Tests {
		_, ipnp, _ := ParseCIDR(tt.inaddrStr)
		ipnn, _ := NewNet(tt.ipaddr, tt.inaddrMask)
		if ipnp.Version() != tt.version {
			t.Errorf("From ParseCIDR %s got Network.Version == %d, expect %d", tt.inaddrStr, ipnp.Version(), tt.version)
		}
		if ipnn.Version() != tt.version {
			t.Errorf("From NewNet %s got Network.Version == %d, want %d", tt.inaddrStr, ipnn.Version(), tt.version)
		}
	}
}

func TestNet4_Count(t *testing.T) {
	for _, tt := range Network4Tests {
		_, ipn, _ := ParseCIDR(tt.inaddrStr)
		ipn4 := ipn.(Net4)
		count, _ := strconv.Atoi(tt.count)
		if ipn4.Count() != uint32(count) {
			t.Errorf("On %s got Network.Count == %d, want %d", tt.inaddrStr, ipn4.Count(), count)
		}
	}
}

func TestNet4_Count4(t *testing.T) {
	for _, tt := range Network4Tests {
		if tt.version == 6 {
			continue
		}
		_, ipn, _ := ParseCIDR(tt.inaddrStr)
		ipn4 := ipn.(Net4)
		count, _ := strconv.Atoi(tt.count)
		if ipn4.Count() != uint32(count) {
			t.Errorf("On %s got Network.Count4 == %d, want %d", tt.inaddrStr, ipn4.Count(), count)
		}
	}
}

func TestNet4_FirstAddress(t *testing.T) {
	for _, tt := range Network4Tests {
		_, ipn, _ := ParseCIDR(tt.inaddrStr)
		if addr := ipn.FirstAddress(); !tt.firstaddr.Equal(addr) {
			t.Errorf("On %s got Network.FirstAddress == %v, want %v", tt.inaddrStr, addr, tt.firstaddr)
		}
	}
}

func TestNet4_finalAddress(t *testing.T) {
	for _, tt := range Network4Tests {
		_, ipn, _ := ParseCIDR(tt.inaddrStr)
		ipn4 := ipn.(Net4)
		if addr, ones := ipn4.finalAddress(); !tt.broadcast.Equal(addr) {
			t.Errorf("On %s got Network.finalAddress == %v, want %v mask length %d)", tt.inaddrStr, addr, tt.broadcast, ones)
		}
	}
}

func TestNet4_LastAddress(t *testing.T) {
	for _, tt := range Network4Tests {
		_, ipn, _ := ParseCIDR(tt.inaddrStr)
		if addr := ipn.LastAddress(); !tt.lastaddr.Equal(addr) {
			t.Errorf("On %s got Network.LastAddress == %v, want %v", tt.inaddrStr, addr, tt.lastaddr)
		}
	}
}

func TestNet4_NetworkAddress(t *testing.T) {
	for _, tt := range Network4Tests {
		if tt.version == 6 {
			continue
		}
		_, ipn, _ := ParseCIDR(tt.inaddrStr)
		if addr := ipn.IP(); !tt.network.Equal(addr) {
			t.Errorf("On %s got Network.IP == %v, want %v", tt.inaddrStr, addr, tt.network)
		}
	}
}

var enumerate4Tests = []struct {
	inaddr string
	total  int
	last   net.IP
}{
	{"192.168.0.0/22", 1022, net.IP{192, 168, 3, 254}},
	{"192.168.0.0/23", 510, net.IP{192, 168, 1, 254}},
	{"192.168.0.0/24", 254, net.IP{192, 168, 0, 254}},
	{"192.168.0.0/25", 126, net.IP{192, 168, 0, 126}},
	{"192.168.0.0/26", 62, net.IP{192, 168, 0, 62}},
	{"192.168.0.0/27", 30, net.IP{192, 168, 0, 30}},
	{"192.168.0.0/28", 14, net.IP{192, 168, 0, 14}},
	{"192.168.0.0/29", 6, net.IP{192, 168, 0, 6}},
	{"192.168.0.0/30", 2, net.IP{192, 168, 0, 2}},
	{"192.168.0.0/31", 2, net.IP{192, 168, 0, 1}},
	{"192.168.0.0/32", 1, net.IP{192, 168, 0, 0}},
}

func TestNet4_Enumerate(t *testing.T) {
	for _, tt := range enumerate4Tests {
		_, ipn, _ := ParseCIDR(tt.inaddr)
		ipn4 := ipn.(Net4)
		addrlist := ipn4.Enumerate(0, 0)
		if len(addrlist) != tt.total {
			t.Errorf("On %s Network.Enumerate(0,0) got size %d, want %d", tt.inaddr, len(addrlist), tt.total)
		}
		x := CompareIPs(tt.last, addrlist[tt.total-1])
		if x != 0 {
			t.Errorf("On %s Network.Enumerate(0,0) got last member %+v, want %+v", tt.inaddr, addrlist[tt.total-1], tt.last)
		}

	}
}

var enumerate4VariableTests = []struct {
	offset int
	size   int
	total  int
	first  net.IP
	last   net.IP
}{
	{0, 0, 1022, net.IP{192, 168, 0, 1}, net.IP{192, 168, 3, 254}},
	{1, 0, 1021, net.IP{192, 168, 0, 2}, net.IP{192, 168, 3, 254}},
	{256, 0, 766, net.IP{192, 168, 1, 1}, net.IP{192, 168, 3, 254}},
	{0, 128, 128, net.IP{192, 168, 0, 1}, net.IP{192, 168, 0, 128}},
	{20, 128, 128, net.IP{192, 168, 0, 21}, net.IP{192, 168, 0, 148}},
	{1000, 100, 22, net.IP{192, 168, 3, 233}, net.IP{192, 168, 3, 254}},
}

func TestNet4_EnumerateWithVariables(t *testing.T) {
	_, ipn, _ := ParseCIDR("192.168.0.0/22")
	ipn4 := ipn.(Net4)
	for _, tt := range enumerate4VariableTests {
		addrlist := ipn4.Enumerate(tt.size, tt.offset)
		if len(addrlist) != tt.total {
			t.Errorf("On Network.Enumerate(%d,%d) got size %d, want %d", tt.size, tt.offset, len(addrlist), tt.total)
		}
		x := CompareIPs(tt.first, addrlist[0])
		if x != 0 {
			t.Errorf("On Network.Enumerate(%d,%d) got first member %+v, want %+v", tt.size, tt.offset, addrlist[0], tt.first)
		}
		y := CompareIPs(tt.last, addrlist[len(addrlist)-1])
		if y != 0 {
			t.Errorf("On Network.Enumerate(%d,%d) got last member %+v, want %+v", tt.size, tt.offset, addrlist[len(addrlist)-1], tt.last)
		}

	}
}

var incr4Tests = []struct {
	inaddr   string
	ipaddr   net.IP
	nextaddr net.IP
	nexterr  error
}{
	{
		"192.168.1.0/23",
		net.IP{192, 168, 1, 0},
		net.IP{192, 168, 1, 1},
		nil,
	},
	{
		"192.168.1.0/24",
		net.IP{192, 168, 1, 254},
		net.IP{192, 168, 1, 255},
		ErrBroadcastAddress,
	},
	{
		"192.168.2.0/24",
		net.IP{192, 168, 2, 1},
		net.IP{192, 168, 2, 2},
		nil,
	},
	{
		"192.168.3.0/24",
		net.IP{192, 168, 3, 0},
		net.IP{192, 168, 3, 1},
		nil,
	},
	{
		"192.168.4.0/24",
		net.IP{192, 168, 5, 1},
		net.IP{},
		ErrAddressOutOfRange,
	},
	{
		"192.168.1.0/31",
		net.IP{192, 168, 1, 0},
		net.IP{192, 168, 1, 1},
		ErrBroadcastAddress,
	},
	{
		"192.168.1.0/32",
		net.IP{192, 168, 1, 0},
		net.IP{},
		ErrAddressAtEndOfRange,
	},
}

func TestNet4_NextIP(t *testing.T) {
	for _, tt := range incr4Tests {
		_, ipn, _ := ParseCIDR(tt.inaddr)
		ipn4 := ipn.(Net4)
		addr, err := ipn4.NextIP(tt.ipaddr)
		if !addr.Equal(tt.nextaddr) {
			t.Errorf("For %s expected %v, got %v", tt.inaddr, tt.nextaddr, addr)
		}
		if err != tt.nexterr {
			t.Errorf("For %s expected \"%v\", got \"%v\"", tt.inaddr, tt.nexterr, err)
		}
	}
}

var decr4Tests = []struct {
	inaddr   string
	ipaddr   net.IP
	prevaddr net.IP
	preverr  error
}{
	{
		"192.168.1.0/23",
		net.IP{192, 168, 1, 0},
		net.IP{192, 168, 0, 255},
		nil,
	},
	{
		"192.168.1.0/24",
		net.IP{192, 168, 1, 254},
		net.IP{192, 168, 1, 253},
		nil,
	},
	{
		"192.168.2.0/24",
		net.IP{192, 168, 2, 1},
		net.IP{192, 168, 2, 0},
		ErrNetworkAddress,
	},
	{
		"192.168.3.0/24",
		net.IP{192, 168, 3, 0},
		net.IP{},
		ErrAddressAtEndOfRange,
	},
	{
		"192.168.4.0/24",
		net.IP{192, 168, 5, 1},
		net.IP{},
		ErrAddressOutOfRange,
	},
	{
		"192.168.1.1/31",
		net.IP{192, 168, 1, 1},
		net.IP{192, 168, 1, 0},
		ErrNetworkAddress,
	},
	{
		"192.168.1.0/32",
		net.IP{192, 168, 1, 0},
		net.IP{},
		ErrAddressAtEndOfRange,
	},
}

func TestNet4_PreviousIP(t *testing.T) {
	for _, tt := range decr4Tests {
		_, ipn, _ := ParseCIDR(tt.inaddr)
		ipn4 := ipn.(Net4)
		addr, err := ipn4.PreviousIP(tt.ipaddr)
		if !addr.Equal(tt.prevaddr) {
			t.Errorf("For %s expected %v, got %v", tt.inaddr, tt.prevaddr, addr)
		}
		if err != tt.preverr {
			t.Errorf("For %s expected \"%v\", got \"%v\"", tt.inaddr, tt.preverr, err)
		}
	}
}

var subnet4Tests = []struct {
	in       string
	prevmask int
	prevnet  string
	nextmask int
	nextnet  string
	submask  int
	subnets  []string
}{
	{
		"192.168.0.0/24",
		24,
		"192.167.255.0/24",
		24,
		"192.168.1.0/24",
		26,
		[]string{"192.168.0.0/26", "192.168.0.64/26", "192.168.0.128/26", "192.168.0.192/26"},
	},
}

func TestNet4_Subnet(t *testing.T) {
	for _, tt := range subnet4Tests {
		_, ipn, _ := ParseCIDR(tt.in)
		ipn4 := ipn.(Net4)
		subnets, _ := ipn4.Subnet(tt.submask)
		v := compareNet4ArraysToStringRepresentation(subnets, tt.subnets)
		if v == false {
			t.Errorf("On Net{%s}.Subnet(%d) expected %v got %v", tt.in, tt.submask, tt.subnets, subnets)
		}
	}
}

func TestNet4_SubnetBadMasklen(t *testing.T) {
	_, ipn, _ := ParseCIDR("192.168.1.0/24")
	ipn4 := ipn.(Net4)
	_, err := ipn4.Subnet(23)
	if err == nil {
		t.Error("Net{192.168.1.0/24}.Subnet(23) expected error, but got none")
	}
}

func TestNet4_PreviousNet(t *testing.T) {
	for _, tt := range subnet4Tests {
		_, ipn, _ := ParseCIDR(tt.in)
		_, pneta, _ := ParseCIDR(tt.prevnet)
		ipn4 := ipn.(Net4)

		pnetb, _ := ipn4.PreviousNet(tt.prevmask)

		if CompareNets(pneta, pnetb) != 0 {
			t.Errorf("On Net{%s}.PreviousNet(%d) expected %s got %s", tt.in, tt.prevmask, tt.prevnet, pneta.String())
		}
	}
}

func TestNet4_NextNet(t *testing.T) {
	for _, tt := range subnet4Tests {
		_, ipn, _ := ParseCIDR(tt.in)
		_, pneta, _ := ParseCIDR(tt.nextnet)
		ipn4 := ipn.(Net4)

		pnetb, _ := ipn4.NextNet(tt.nextmask)

		if CompareNets(pneta, pnetb) != 0 {
			t.Errorf("On Net{%s}.NextNet(%d) expected %s got %s", tt.in, tt.nextmask, tt.nextnet, pneta.String())
		}
	}
}

var supernetTests = []struct {
	in      string
	masklen int
	out     string
	err     error
}{
	{
		"192.168.1.0/24",
		23,
		"192.168.0.0/23",
		nil,
	},
	{
		"192.168.1.0/24",
		0,
		"192.168.0.0/23",
		nil,
	},
	{
		"192.168.1.0/24",
		22,
		"192.168.0.0/22",
		nil,
	},
	{
		"192.168.1.4/30",
		24,
		"192.168.1.0/24",
		nil,
	},
}

func TestNet4_Supernet(t *testing.T) {
	for _, tt := range supernetTests {
		_, ipn, _ := ParseCIDR(tt.in)
		ipn4 := ipn.(Net4)
		onet, _ := ipn4.Supernet(tt.masklen)
		if onet.String() != tt.out {
			t.Errorf("On Net{%s}.Supernet(%d) expected %s got %s", tt.in, tt.masklen, tt.out, onet.String())
		}
	}
}

var compareNetworks = map[int]string{
	0: "192.168.0.0/16",
	1: "192.168.0.0/23",
	2: "192.168.1.0/24",
	3: "192.168.1.0/24",
	4: "192.168.3.0/26",
	5: "192.168.3.64/26",
	6: "192.168.3.128/25",
	7: "192.168.4.0/24",
}

func TestCompareNets(t *testing.T) {
	a := ByNet{}
	for _, v := range compareNetworks {
		_, ipn, _ := ParseCIDR(v)
		a = append(a, ipn)
	}
	sort.Sort(ByNet(a))
	for k, v := range compareNetworks {
		if a[k].String() != v {
			t.Errorf("Subnet %s not at expected position %d. Got %s instead", v, k, a[k].String())
		}

	}

}

var compareCIDR = []struct {
	network string
	subnet  string
	result  bool
}{
	{"192.168.0.0/16", "192.168.45.0/24", true},
	{"192.168.45.0/24", "192.168.45.0/26", true},
	{"192.168.45.0/24", "192.168.46.0/26", false},
	{"10.1.1.1/24", "10.0.0.0/8", false},
}

func TestNet_ContainsNetwork(t *testing.T) {
	for _, cidr := range compareCIDR {
		_, ipn, _ := ParseCIDR(cidr.network)
		_, sub, _ := ParseCIDR(cidr.subnet)
		result := ipn.ContainsNet(sub)
		if result != cidr.result {
			t.Errorf("For \"%s contains %s\" expected %v got %v", cidr.network, cidr.subnet, cidr.result, result)
		}
	}
}

func compareNet4ArraysToStringRepresentation(a []Net4, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, n := range a {
		if n.String() != b[i] {
			return false
		}
	}

	return true
}
