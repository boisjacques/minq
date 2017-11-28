package main

import (
	"fmt"
	"github.com/boisjacques/minq"
)

func main() {
	ah := minq.NewAddressHelper()
	addresses := ah.GetAddresses()
	fmt.Println(addresses)
}
