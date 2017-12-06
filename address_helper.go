package minq

import (
	"github.com/boisjacques/golang-utils"
	"log"
	"net"
	"strings"
	"sync"
)

type AddressHelper struct {
	ipAddresses       map[string]*net.UDPAddr
	ipAddressesBool   map[string]bool
	listeners         []chan *net.UDPAddr
	lockAddresses     sync.RWMutex
	lockAddressesBool sync.RWMutex
}

func NewAddressHelper() *AddressHelper {
	ah := AddressHelper{
		make(map[string]*net.UDPAddr),
		make(map[string]bool),
		make([]chan *net.UDPAddr, 0),
		sync.RWMutex{},
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
	a.lockAddresses.Lock()
	a.lockAddressesBool.Lock()
	log.Println("locked: ", util.Tracer())
	defer a.lockAddresses.Unlock()
	defer a.lockAddressesBool.Unlock()
	defer log.Println("unlocked: ", util.Tracer())
	for key, value := range a.ipAddressesBool {
		if value == false {
			delete(a.ipAddresses, key)
			a.Publish(a.ipAddresses[key])
		}
	}
}

func (a *AddressHelper) GetAddresses() *map[string]*net.UDPAddr {
	a.lockAddresses.RLock()
	log.Println("locked: ", util.Tracer())
	defer a.lockAddresses.RUnlock()
	defer log.Println("unlocked: ", util.Tracer())
	return &a.ipAddresses
}

func (a *AddressHelper) write(addr *net.UDPAddr, bool bool) {
	a.lockAddresses.Lock()
	a.lockAddressesBool.Lock()
	log.Println("locked: ", util.Tracer())
	defer a.lockAddressesBool.Unlock()
	defer a.lockAddresses.Unlock()
	defer log.Println("unlocked: ", util.Tracer())
	a.ipAddresses[addr.String()] = addr
	a.ipAddressesBool[addr.String()] = bool
}

func (a *AddressHelper) containsBlocking(addr *net.UDPAddr) bool {
	a.lockAddresses.RLock()
	log.Println("locked: ", util.Tracer())
	defer a.lockAddresses.RUnlock()
	defer log.Println("unlocked: ", util.Tracer())
	_, contains := a.ipAddresses[addr.String()]
	return contains
}

func (a *AddressHelper) falsifyAddresses() {
	a.lockAddresses.Lock()
	a.lockAddressesBool.Lock()
	log.Println("locked: ", util.Tracer())
	defer a.lockAddresses.Unlock()
	defer a.lockAddressesBool.Unlock()
	defer log.Println("unlocked: ", util.Tracer())
	for address := range a.ipAddresses {
		a.ipAddressesBool[address] = false
	}
}
