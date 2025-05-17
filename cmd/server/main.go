package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp4", ":8080")

	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			panic(err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Received %d bytes: %s\n", n, string(buffer[:n]))

	_, err = conn.Write(buffer[:n])
	if err != nil {
		panic(err)
	}
}
