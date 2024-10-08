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

const SEVEN uint32 = 7

func getCorrelationId(buf []byte) []byte {
	return buf[8:12]
}

func setCorrelationId(res []byte, correlationId uint32) {
	binary.BigEndian.PutUint32(res[4:], correlationId)
}

func main() {
	validAPIVersions := []int{0, 1, 2, 3}
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("Successfully set up connection")

	// Read APIVersions requests
	// NOTE: we must read in request sent to this server
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection:", err.Error())
		return
	}
	fmt.Printf("Received %d bytes: %s\n", n, string(buf[:n]))

	// first 4 bytes are 0, last 4 bytes should represent correlation_id field from request
	response := make([]byte, 8)

	// NOTE:
	// read 4 bytes, starting from the 8th byte, since correlation_id is 32 bits.
	// first 4 bytes from the buffer are for the message length!!
	correlationId := binary.BigEndian.Uint32(getCorrelationId(buf))
	setCorrelationId(response, correlationId)
	conn.Write(response)

	fmt.Println("Connection handled successfully")
	l.Close()
	conn.Close() // Ensure the connection is closed properly
}
