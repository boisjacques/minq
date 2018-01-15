package minq

import (
	"testing"
	"net"
	"strconv"
	"fmt"
	"sync"
	"time"
)


func TestSchedulingAlgorithm(t *testing.T) {
	results := make(map[*Path]int)
	pathIds := make([]uint32,0)
	ah := NewAddressHelper()
	/*
	local,err := net.ResolveUDPAddr("udp", "127.0.0.1:44333")
	if err != nil{
		log.Println(err)
		return
	}
	remote,err := net.ResolveUDPAddr("udp", "127.0.0.2:44333")
	if err != nil{
		log.Println(err)
		return
	}
	usock, err := ah.openSocket(local)
	if err != nil {
		log.Println(err)
		return
	}
	initTrans := NewUdpTransport(usock, remote)

	pathZero := Path{
		nil,
		true,
		initTrans,
		0,
		1000,
		0,
		nil,
		nil,
	}
	*/

	conn := Connection{

	}
	conn.log = newConnectionLogger(&conn)

	scheduler := Scheduler{
		make(map[uint32]*Path),
		&conn,
		0,
		nil,
		pathIds,
		0,
		ah,
		make(chan string),
		nil,
		nil,
		make(map[string]*net.UDPAddr),
		sync.RWMutex{},
		sync.RWMutex{},
		false,
		1000,
	}

	for i := 0; i < 10; i++ {
		j :=  i + 1
		local,err := net.ResolveUDPAddr("udp", "127.0.0.1:443" + strconv.Itoa(j))
		if err != nil{
			fmt.Println(err)
			return
		}
		remote,err := net.ResolveUDPAddr("udp", "127.0.0.1:442" + strconv.Itoa(j))
		if err != nil{
			fmt.Println(err)
			return
		}
		scheduler.newPath(local, remote)
	}
	fmt.Println("Created Paths")

	for k := 0; k < 100000 ; k++ {
		path,err := scheduler.weightedSelect()
		if err != nil {
			fmt.Println(err)
			return
		}
		results[path]++
	}
	fmt.Println("Selected Paths")
	time.Sleep(100 * time.Millisecond)
	for path,amount := range results{
		fmt.Println("Path: ", path.GetPathID())
		fmt.Println("\tWeight: ", path.GetWeight())
		fmt.Println("\tCount: ", amount)
	}

	for _,socket := range ah.sockets {
		socket.Close()
	}
	fmt.Println("Test finished")
}