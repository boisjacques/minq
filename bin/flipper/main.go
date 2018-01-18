package main

import (
	"os"
	"fmt"
	"io/ioutil"
)

func main() {
	filename := os.Args[1]
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print(err)
	}
	for i, _ := range b {
		if b[i] > 0x40 {
			b[i] ^= 0x20
		}
	}
	err = ioutil.WriteFile("flipped-" + filename, b, 0644)
	if err != nil {
		fmt.Println(err)
	}
}