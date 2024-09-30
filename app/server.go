package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

const SEVEN int32 = 7

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}
	defer l.Close()
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close() // Ensure the connection is closed properly

	fmt.Println("Successfully set up connection")

	// Read APIVersions requests
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection:", err.Error())
		return
	}
	fmt.Printf("Received %d bytes: %s\n", n, string(buf[:n]))

	// first 4 bytes are 0, last 4 bytes should represent 7 in big endian byte
	response := make([]byte, 8)
	binary.BigEndian.PutUint32(response, uint32(SEVEN))
	conn.Write(response)

	fmt.Println("Connection handled successfully")
}
