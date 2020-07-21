package iputil

import (
	"net"
	"sync"
)

type Options struct {
	AllowedIPs []string
	BlockedIPs []string

	// explicity allowed country ISO codes
	AllowedAreas []string

	// explicity blocked country ISO codes
	BlockedAreas []string

	// block by default (defaults to allow)
	DefaultBlocked bool
}

type IPFilter struct {
	opts           Options
	mut            sync.RWMutex
	defaultAllowed bool

	ips map[string]bool

	// countries code allowed/blocked
	codes map[string]bool

	subnets []*subnet
}

type subnet struct {
	str     string
	ipnet   *net.IPNet
	allowed bool
}

func New(opts Options) *IPFilter {
	f := &IPFilter{
		opts:           opts,
		ips:            map[string]bool{},
		codes:          map[string]bool{},
		defaultAllowed: !opts.DefaultBlocked,
	}

	for _, ip := range opts.BlockedIPs {
		f.BlockIP(ip)
	}
	for _, ip := range opts.AllowedIPs {
		f.AllowIP(ip)
	}
	for _, code := range opts.BlockedAreas {
		f.BlockArea(code)
	}
	for _, code := range opts.AllowedAreas {
		f.AllowArea(code)
	}

	return f
}

func Default() *IPFilter {
	return New(Options{})
}

func (f *IPFilter) AllowIP(ip string) bool {
	return f.toggleIP(ip, true)
}

func (f *IPFilter) BlockIP(ip string) bool {
	return f.toggleIP(ip, false)
}

func (f *IPFilter) toggleIP(str string, allowed bool) bool {
	// check if subnet
	if ip, ipnet, err := net.ParseCIDR(str); err == nil {
		// only one ip
		if n, total := ipnet.Mask.Size(); n == total {
			f.mut.Lock()
			f.ips[ip.String()] = allowed
			f.mut.Unlock()
			return true
		}

		//check for existing
		f.mut.Lock()
		found := false
		for _, subnet := range f.subnets {
			if subnet.str == str {
				found = true
				subnet.allowed = allowed
				break
			}
		}
		if !found {
			f.subnets = append(f.subnets, &subnet{
				str:     str,
				ipnet:   ipnet,
				allowed: allowed,
			})
		}
		f.mut.Unlock()
		return true
	}

	//check if plain ip
	if ip := net.ParseIP(str); ip != nil {
		f.mut.Lock()
		f.ips[ip.String()] = allowed
		f.mut.Unlock()
		return true
	}
	return false
}

// TODO
func (f *IPFilter) RemoveIP(str string) bool {
	// check if subnet
	if ip, ipnet, err := net.ParseCIDR(str); err == nil {
		// only one ip
		if n, total := ipnet.Mask.Size(); n == total {
			f.mut.Lock()
			delete(f.ips, ip.String())
			f.mut.Unlock()
			return true
		}

		//check for existing
		f.mut.Lock()
		for _, subnet := range f.subnets {
			if subnet.str == str {
				break
			}
		}
		f.mut.Unlock()
		return true
	}

	//check if plain ip
	if ip := net.ParseIP(str); ip != nil {
		f.mut.Lock()
		delete(f.ips, ip.String())
		f.mut.Unlock()
		return true
	}
	return false
}

func (f *IPFilter) AllowArea(code string) {
	f.toggleArea(code, true)
}

func (f *IPFilter) BlockArea(code string) {
	f.toggleArea(code, false)
}

func (f *IPFilter) toggleArea(code string, allowed bool) {
	f.mut.Lock()
	f.codes[code] = allowed
	f.mut.Unlock()
}

func (f *IPFilter) RemoveCountry(code string) {
	f.mut.Lock()
	delete(f.codes, code)
	f.mut.Unlock()
}

func (f *IPFilter) ToggleDefault(allowed bool) {
	f.mut.Lock()
	f.defaultAllowed = allowed
	f.mut.Unlock()
}

func (f *IPFilter) Allowed(ip string) bool {
	return f.NetAllowed(net.ParseIP(ip))
}

func (f *IPFilter) Blocked(ip string) bool {
	return !f.Allowed(ip)
}

func (f *IPFilter) NetAllowed(ip net.IP) bool {
	if ip == nil {
		return false
	}

	// check single ip list
	allowed, ok := f.ips[ip.String()]
	if ok {
		return allowed
	}

	// scan subnets
	for _, subnet := range f.subnets {
		if subnet.ipnet.Contains(ip) {
			return subnet.allowed
		}
	}

	return f.defaultAllowed
}

func (f *IPFilter) NetBlocked(ip net.IP) bool {
	return !f.NetAllowed(ip)
}

func (f *IPFilter) AreaAllowed(code string) bool {
	// check code
	if code != "" {
		if allowed, ok := f.codes[code]; ok {
			return allowed
		}
	}

	return f.defaultAllowed
}

func (f *IPFilter) AreaBlocked(code string) bool {
	return !f.AreaAllowed(code)
}