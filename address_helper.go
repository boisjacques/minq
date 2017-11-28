package minq

import (
	"fmt"
	"net"
	"strings"
)

type AddressHelper struct {
	ipAddressPtr map[*net.UDPAddr]bool
	ipAddresses  []net.UDPAddr
	listeners    []chan *net.UDPAddr
}

func NewAddressHelper() *AddressHelper {
	ah := AddressHelper{
		make(map[*net.UDPAddr]bool),
		make([]net.UDPAddr, 0),
		make([]chan *net.UDPAddr, 0),
	}
	ah.GatherAddresses()
	return &ah
}

func (a *AddressHelper) Subscribe(c chan *net.UDPAddr) {
	a.listeners = append(a.listeners, c)
}

func (a *AddressHelper) Publish(msg *net.UDPAddr) {
	if len(a.listeners) > 0 {
		for _, c := range a.listeners {
			c <- msg
		}
	}
}

func (a *AddressHelper) GatherAddresses() {
	for address := range a.ipAddressPtr {
		a.ipAddressPtr[address] = false
	}
	interfaces, _ := net.Interfaces()
	for _, iface := range interfaces {
		flags := iface.Flags.String()
		if !strings.Contains(flags, "loopback") {
			addrs, _ := iface.Addrs()
			ipv6 := false
			for _, addr := range addrs {
				if !(strings.Contains(addr.String(), "127") ||
					strings.Contains(addr.String(), "fe80")) {
					if strings.Contains(addr.String(), ":") {

					}
					arr := strings.Split(addr.String(), "/")
					udpAddr, err := net.ResolveUDPAddr("udp", arr[0]+":4433")
					if err != nil {
						fmt.Println("Error parsing IP address: ", addr)
						fmt.Println(err)
					} else {
						_, containsAddr := a.ipAddressPtr[udpAddr]
						if containsAddr {
							a.ipAddressPtr[udpAddr] = true
						}
						if !containsAddr {
							if (udpAddr.IP.To4 != nil) || (ipv6 == false && udpAddr.IP.To4() == nil) {
								if udpAddr.IP.To4() == nil {
									ipv6 = true
								}
								a.ipAddressPtr[udpAddr] = true
								a.Publish(udpAddr)
							}
						}
					}
				}
			}
		}
	}
	a.cleanUp()
}

func (a *AddressHelper) cleanUp() {
	for address, value := range a.ipAddressPtr {
		if value == false {
			delete(a.ipAddressPtr, address)
			a.Publish(address)
		}
	}
}

func (a *AddressHelper) GetAddresses() *map[*net.UDPAddr]bool {
	return &a.ipAddressPtr
}
