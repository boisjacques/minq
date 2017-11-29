package minq

import (
	"log"
	"net"
	"strings"
	"sync"
)

type AddressHelper struct {
	ipAddressPtr map[*net.UDPAddr]bool
	ipAddresses  []net.UDPAddr
	listeners    []chan *net.UDPAddr
	lock         sync.RWMutex
}

func NewAddressHelper() *AddressHelper {
	ah := AddressHelper{
		make(map[*net.UDPAddr]bool),
		make([]net.UDPAddr, 0),
		make([]chan *net.UDPAddr, 0),
		sync.RWMutex{},
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
	a.falsifyAddresses()
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
					if strings.Contains(arr[0], ":") {
						arr[0] = "[" + arr[0] + "]"
					}
					udpAddr, err := net.ResolveUDPAddr("udp", arr[0]+":4433")
					if err != nil {
						log.Println(err)
					} else {
						if a.containsBlocking(udpAddr) {
							a.write(udpAddr, true)
						}
						if !a.containsBlocking(udpAddr) {
							if (udpAddr.IP.To4 != nil) || (ipv6 == false && udpAddr.IP.To4() == nil) {
								if udpAddr.IP.To4() == nil {
									ipv6 = true
								}
								a.write(udpAddr, true)
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
	a.lock.Lock()
	defer a.lock.Unlock()
	for address, value := range a.ipAddressPtr {
		if value == false {
			delete(a.ipAddressPtr, address)
			a.Publish(address)
		}
	}
}

func (a *AddressHelper) GetAddresses() *map[*net.UDPAddr]bool {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return &a.ipAddressPtr
}

func (a *AddressHelper) write(addr *net.UDPAddr, bool bool) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.ipAddressPtr[addr] = bool
}

func (a *AddressHelper) containsBlocking(addr *net.UDPAddr) bool {
	a.lock.RLock()
	defer a.lock.RUnlock()
	_, contains := a.ipAddressPtr[addr]
	return contains
}

func (a *AddressHelper) falsifyAddresses() {
	a.lock.Lock()
	defer a.lock.Unlock()
	for address := range a.ipAddressPtr {
		a.ipAddressPtr[address] = false
	}
}
