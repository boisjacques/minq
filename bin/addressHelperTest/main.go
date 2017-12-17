package main

import (
	"fmt"
	"github.com/boisjacques/minq"
)

func main() {
	ah := minq.NewAddressHelper()
	for i := 0; i < 10; i++ {
		ah.GatherAddresses()
	}
	addresses := ah.GetAddresses()
	fmt.Println(addresses)
}
