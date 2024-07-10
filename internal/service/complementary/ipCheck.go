package complementary

import "net"

// IPRangeCotains gets an ip address and a match IP range.
// If the given ip range contains the give ip address, returns true
func IPRangeCotains(ipToCheck string, ReferenceIPRange string) bool {

	// If the given ReferenceIPRange is not an absolute IP
	if _, cidr, err := net.ParseCIDR(ReferenceIPRange); err == nil {
		return cidr.Contains(net.ParseIP(ipToCheck))
	}

	// If the given ReferenceIPRange is an absolute IP
	return ipToCheck == ReferenceIPRange

}
