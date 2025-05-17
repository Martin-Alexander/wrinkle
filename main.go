package main

import (
	"net"
)

func main() {
	listener, err := net.Listen("tcp4", "8080")

	if err != nil {
		panic(err)
	}

	listener.Accept()
}