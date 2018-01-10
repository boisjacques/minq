package minq

import "net"

type Path struct {
	connection *Connection
	isPathZero bool
	transport  Transport
	pathID     uint32
	weight     int

	owd        uint64
	local      *net.UDPAddr
	remote     *net.UDPAddr
}

func NewPath(connection *Connection, transport Transport, pathId uint32, local, remote *net.UDPAddr, weight int) *Path {
	return &Path{
		connection,
		false,
		transport,
		pathId,
		weight,
		0,
		local,
		local,
	}
}

// Send a ping Frame, evaluate the RTT in relation to pathZero
/*
func (p *Path) updateMetric(referenceRTT uint16) uint16 {
	frames := make([]frame, 1)
	f := newPingFrame()
	frames[0] = f
	p.connection.sendFramesInPacket(0, frames)
	p.connection.send
	b := make([]byte, 1024)
	p.connection.streams[0].Read(b)
	return owd
}
*/

func (p *Path) GetWeight() int {
	return p.weight
}

func (p *Path) GetPathID() uint32 {
	return p.pathID
}

func (p *Path) setOwd(owd int64) {
	p.owd = uint64(owd)
}

func (p *Path) String() string {
	return "local: " + p.local.String() + " " + "remote: " + p.remote.String()
}

func (p *Path) contains(address string) bool {
	return (p.local.String() == address || p.remote.String() == address)
}
