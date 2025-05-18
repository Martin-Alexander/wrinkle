package main

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	cert, err := tls.LoadX509KeyPair("/app/.ssl/server.crt", "/app/.ssl/server.key")
	if err != nil {
		log.Fatalf("Failed to load cert: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	fmt.Println("Starting server on port 54321...")

	listener, err := net.Listen("tcp4", ":54321")

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
		}

		backendConn, err := connectToBackend()
		if err != nil {
			log.Printf("Failed to connect to backend: %v", err)
			return
		}

		log.Printf("Accepted connection from %s", conn.RemoteAddr())

		go handleConnection(conn, backendConn, tlsConfig)
	}
}

func handleConnection(frontendConn net.Conn, backendConn net.Conn, tlsConfig *tls.Config) {
	// Set read/write timeout for both connections
	timeout := time.Second * 30
	if err := setReadWriteTimeout(frontendConn, timeout); err != nil {
		log.Printf("Failed to set timeout: %v", err)
		return
	}
	if err := setReadWriteTimeout(backendConn, timeout); err != nil {
		log.Printf("Failed to set timeout: %v", err)
		return
	}

	// Read the TLS negotiation packet
	pgTlsNegotiationPacket := getThePgTlsNegotiationPacket()
	firstEightBytes, err := readExactlyNBytes(frontendConn, 8)
	if err != nil {
		log.Printf("Failed to read first 8 bytes: %v", err)
		return
	}
	if string(firstEightBytes) != string(pgTlsNegotiationPacket) {
		log.Printf("Invalid first 8 bytes: %v", firstEightBytes)
		return
	}

	logPacket("TLS negotiation packet", firstEightBytes)

	// Send the TLS negotiation packet to the backend
	if _, err := backendConn.Write(pgTlsNegotiationPacket); err != nil {
		log.Printf("Failed to write to backend: %v", err)
		return
	}

	// Read the response bytes from the backend
	responseByte, err := readExactlyNBytes(backendConn, 1)
	if err != nil {
		log.Printf("Failed to read response from backend: %v", err)
		return
	}
	if len(responseByte) < 1 {
		log.Printf("Invalid response length: %d", len(responseByte))
		return
	}

	logPacket("TLS status byte", responseByte)

	// Send the response byte to the frontend
	if _, err := frontendConn.Write(responseByte); err != nil {
		log.Printf("Failed to write response: %v", err)
		return
	}

	// Close if response was not 'S'
	if responseByte[0] != 'S' {
		log.Printf("Response was not 'S': %v", responseByte)
		return
	}

	frontendTlsConn, err := frontendTlsHandshake(frontendConn, tlsConfig)
	if err != nil {
		log.Printf("Failed to upgrade frontend connection to TLS: %v", err)
		return
	}

	backendTlsConn, err := backendTlsHandshake(backendConn)
	if err != nil {
		log.Printf("Failed to upgrade backend connection to TLS: %v", err)
		return
	}

	go relay(frontendTlsConn, backendTlsConn, "PostgreSQL client")
	go relay(backendTlsConn, frontendTlsConn, "PostgreSQL server")
}

func readExactlyNBytes(conn net.Conn, n int) ([]byte, error) {
	buffer := make([]byte, n)
	if _, err := io.ReadFull(conn, buffer); err != nil {
		return nil, err
	}

	return buffer, nil
}

func readUpTo(conn net.Conn, maxBytes int) ([]byte, error) {
    buf := make([]byte, maxBytes)
    n, err := conn.Read(buf)
    if err != nil {
        return nil, err
    }
    return buf[:n], nil
}

func getThePgTlsNegotiationPacket() []byte {
	return []byte{0x00, 0x00, 0x00, 0x08, 0x04, 0xd2, 0x16, 0x2f} 
}

func frontendTlsHandshake(conn net.Conn, config *tls.Config) (*tls.Conn, error) {
	tlsConn := tls.Server(conn, config)
	err := tlsConn.Handshake()
	if err != nil {
		return nil, err
	}
	return tlsConn, nil
}

func backendTlsHandshake(conn net.Conn) (*tls.Conn, error) {
	tlsConn := tls.Client(conn, &tls.Config{
		InsecureSkipVerify: true,
	})
	err := tlsConn.Handshake()
	if err != nil {
		return nil, err
	}
	return tlsConn, nil
}

func connectToBackend() (net.Conn, error) {
	conn, err := net.Dial("tcp4", "postgres:5432")
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func setReadWriteTimeout(conn net.Conn, timeout time.Duration) error {
	if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}

	if err := conn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}

	return nil
}

func relay(sourceConn net.Conn, destinationConn net.Conn, logLabel string) {
	for {
		data, err := readUpTo(sourceConn, 4096)
		if err != nil {
			log.Printf("Failed to read from connection: %v", err)
			return
		}

		if _, err := destinationConn.Write(data); err != nil {
			log.Printf("Error writing to destination: %v", err)
			return
		}

		logPacket(logLabel, data)
	}
}

func logPacket(label string, data []byte) {
	log.Printf(" -- [%s]:\n%s\n", label, hex.Dump(data))
}