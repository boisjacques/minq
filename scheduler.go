package minq

import (
	"fmt"
	"github.com/boisjacques/golang-utils"
	"hash/adler32"
	"net"
	"sync"
	"time"
)

type direcionAddr uint8

const (
	local  = direcionAddr(0)
	remote = direcionAddr(1)
)

type Scheduler struct {
	paths         map[uint32]Path
	connection    *Connection
	referenceRTT  uint16
	pathZero      *Path
	pathIds       []uint32
	lastPath      uint32
	addressHelper *AddressHelper
	addrChan      chan *net.UDPAddr
	localAddrs    map[*net.UDPAddr]bool
	remoteAddrs   map[*net.UDPAddr]struct{}
	lockRemote    sync.RWMutex
	lockPaths     sync.RWMutex
	isInitialized bool
}

func NewScheduler(initTrans Transport, connection *Connection, ah *AddressHelper) Scheduler {
	connection.log(logTypeMultipath, "New scheduler built for connection %v", connection.clientConnId)
	pathZero := &Path{
		connection,
		true,
		initTrans,
		0,
		100,
		0,
		nil,
		nil,
	}
	paths := make(map[uint32]Path)
	paths[pathZero.pathID] = *pathZero
	pathIds := make([]uint32, 0)
	pathIds = append(pathIds, pathZero.pathID)
	return Scheduler{
		paths,
		connection,
		0,
		pathZero,
		pathIds,
		0,
		ah,
		make(chan *net.UDPAddr),
		ah.ipAddressPtr,
		make(map[*net.UDPAddr]struct{}),
		sync.RWMutex{},
		sync.RWMutex{},
		false,
	}
}

// TODO: Implement proper scheduling, simple round robin right now
func (s *Scheduler) Send(p []byte) error {
	s.lastPath = s.lastPath + 1%uint32(len(s.pathIds))
	err := s.paths[s.pathIds[s.lastPath]].transport.Send(p)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if s.lastPath == 0 {
		s.connection.log(logTypeMultipath, "Packet sent. Used path zero")
	} else {
		s.connection.log(logTypeMultipath, "Packet sent. local: %v \n remote: %x", *s.paths[s.pathIds[int(s.lastPath)]].local, *s.paths[s.pathIds[int(s.lastPath)]].remote)
	}
	return nil
}

// TODO: Consider using CRC32 instead of adler32
func (s *Scheduler) newPath(local, remote *net.UDPAddr) {
	usock, err := net.ListenUDP("udp", local)
	if err != nil {
		s.connection.log(logTypeMultipath, "Error while creating path local IP: %x remote IP %v", *local, *remote)
	}
	transport := NewUdpTransport(usock, remote)
	checksum := adler32.Checksum(xor([]byte(local.String()), []byte(remote.String())))
	p := NewPath(s.connection, transport, checksum, local, remote)
	s.connection.log(logTypeMultipath, "Path successfully created. Endpoints: local %v remote %x", local, remote)
	//p.updateMetric(s.referenceRTT)
	s.paths[p.pathID] = *p
	s.pathIds = append(s.pathIds, p.pathID)
}

func (s *Scheduler) addLocalAddress(local net.UDPAddr) {
	s.connection.log(logTypeMultipath, "Adding local address %v", local)
	for remote := range s.remoteAddrs {
		if isSameVersion(&local, remote) {
			s.newPath(&local, remote)
		}
	}
}

func (s *Scheduler) addRemoteAddress(remote *net.UDPAddr) {
	s.connection.log(logTypeMultipath, "Adding remote address %v", *remote)
	s.remoteAddrs[remote] = struct{}{}
	s.addressHelper.lock.RLock()
	s.connection.log(logTypeMutex, "locked: ", util.Tracer())
	defer s.addressHelper.lock.RUnlock()
	defer s.connection.log(logTypeMutex, "unlocked: ", util.Tracer())
	for local := range s.localAddrs {
		if isSameVersion(local, remote) {
			s.newPath(local, remote)
		}
	}
}

func (s *Scheduler) removeAddress(address *net.UDPAddr) {
	if s.containsBlocking(address, remote) {
		s.delete(address, remote)
		s.connection.log(logTypeMultipath, "Deleted remote address %v", *address)
	}
	if s.containsBlocking(address, local) {
		s.delete(address, local)
		s.connection.log(logTypeMultipath, "Deleted local address %v", *address)
	}
	for k, v := range s.paths {
		if v.contains(address.String()) {
			s.removePath(k)
		}
	}
}

func (s *Scheduler) initializePaths() {
	s.addressHelper.lock.RLock()
	s.lockRemote.RLock()
	s.connection.log(logTypeMutex, "locked: ", util.Tracer())
	defer s.addressHelper.lock.RUnlock()
	defer s.lockRemote.RUnlock()
	defer s.connection.log(logTypeMutex, "unlocked: ")
	for local := range s.localAddrs {
		for remote := range s.remoteAddrs {
			if isSameVersion(local, remote) {
				s.newPath(local, remote)
			}
		}
	}
	s.connection.log(logTypeMultipath, "First flight paths initialized")
	s.isInitialized = true
}

func (s *Scheduler) removePath(pathId uint32) {
	delete(s.paths, pathId)
	s.connection.log(logTypeMultipath, "Removed path %v", pathId)
}

func (s *Scheduler) ListenOnChannel() {
	s.addressHelper.Subscribe(s.addrChan)
	s.connection.log(logTypeMultipath, "Subscribes to Address Helper")
	go func() {
<<<<<<< HEAD
=======
		index := 0
>>>>>>> master
		oldTime := time.Now().Second()
		for {
			if s.connection.state == StateEstablished {
				addr := <-s.addrChan
				if !s.containsBlocking(addr, local) {
					s.write(addr)
					s.connection.sendFramesInPacket(packetType1RTTProtectedPhase1, s.assembleAddrModFrame(kAddAddress, *addr))
				} else {
					s.delete(addr, local)
					s.connection.sendFramesInPacket(packetType1RTTProtectedPhase1, s.assembleAddrModFrame(kDeleteAddress, *addr))
				}
			} else {
<<<<<<< HEAD
				if time.Now().Second()-oldTime == 10 {
=======
				if time.Now().Second()-oldTime == 2 {
>>>>>>> master
					s.connection.log(logTypeMultipath, "Waiting for connection establishment", util.Tracer())
					oldTime = time.Now().Second()
				}
			}
		}
	}()
}

func (s *Scheduler) assebleAddrArrayFrame() []frame {
	arr := make([]net.UDPAddr, 0)
	s.addressHelper.lock.RLock()
	s.connection.log(logTypeMutex, "locked: ", util.Tracer())
	defer s.addressHelper.lock.RUnlock()
	defer s.connection.log(logTypeMutex, "unlocked: ", util.Tracer())
	for k, _ := range s.localAddrs {
		arr = append(arr, *k)
	}
	frames := make([]frame, 0)
	frame := newAddrArrayFrame(0, arr)
	frames = append(frames, frame)
	s.connection.log(logTypeMultipath, "Assembled frame", frame)
	return frames
}

func (s *Scheduler) assembleAddrModFrame(delete operation, addr net.UDPAddr) []frame {
	frames := make([]frame, 0)
	frame := newAddrModFrame(0, delete, addr)
	frames = append(frames, frame)
	s.connection.log(logTypeMultipath, "Assembled frame", frame)
	return frames
}

func xor(local, remote []byte) []byte {
	rval := make([]byte, 0)
	for i := 0; i < len(local); i++ {
		rval[i] = local[i] ^ remote[i]
	}

	return rval
}

func isSameVersion(local, remote *net.UDPAddr) bool {
	if local.IP.To4() == nil && remote.IP.To4() == nil {
		return true
	}

	if local.IP.To4() != nil && remote.IP.To4() != nil {
		return true
	}
	return false
}

func (s *Scheduler) containsBlocking(addr *net.UDPAddr, direcion direcionAddr) bool {
	var contains bool
	if direcion == local {
		s.addressHelper.lock.RLock()
		s.connection.log(logTypeMutex, "locked: ", util.Tracer())
		defer s.addressHelper.lock.RUnlock()
		defer s.connection.log(logTypeMutex, "unlocked: ", util.Tracer())
		_, contains = s.localAddrs[addr]
	} else if direcion == remote {
		s.lockRemote.Lock()
		s.connection.log(logTypeMutex, "locked: ", util.Tracer())
		defer s.lockRemote.Unlock()
		defer s.connection.log(logTypeMutex, "unlocked: ", util.Tracer())
		_, contains = s.remoteAddrs[addr]
	}
	return contains
}

func (s *Scheduler) delete(addr *net.UDPAddr, direction direcionAddr) {
	if direction == local {
		s.addressHelper.lock.Lock()
		s.connection.log(logTypeMutex, "locked: ", util.Tracer())
		defer s.addressHelper.lock.Unlock()
		defer s.connection.log(logTypeMutex, "unlocked: ", util.Tracer())
		delete(s.localAddrs, addr)
	}
	if direction == remote {
		s.lockRemote.Lock()
		s.connection.log(logTypeMutex, "locked: ", util.Tracer())
		defer s.lockRemote.Unlock()
		defer s.connection.log(logTypeMutex, "unlocked: ", util.Tracer())
		delete(s.remoteAddrs, addr)
	}
}

func (s *Scheduler) deletePath(pathId uint32) {
	s.lockPaths.Lock()
	s.connection.log(logTypeMutex, "locked: ", util.Tracer())
	defer s.lockPaths.Unlock()
	defer s.connection.log(logTypeMutex, "unlocked: ", util.Tracer())
	delete(s.paths, pathId)
}

func (s *Scheduler) write(addr *net.UDPAddr) {
	s.addressHelper.lock.Lock()
	s.connection.log(logTypeMutex, "locked: ", util.Tracer())
	defer s.addressHelper.lock.Unlock()
	defer s.connection.log(logTypeMutex, "unlocked: ", util.Tracer())
	s.localAddrs[addr] = false
}
