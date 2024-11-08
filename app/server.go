package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"slices"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func getMessageLen(buf []byte) uint32 {
	return binary.BigEndian.Uint32(buf[0:4])
}

func getRequestApiKey(buf []byte) uint16 {
	return binary.BigEndian.Uint16(buf[4:6])
}

func getRequestApiVersion(buf []byte) uint16 {
	return binary.BigEndian.Uint16(buf[6:8])
}

func getCorrelationId(buf []byte) uint32 {
	return binary.BigEndian.Uint32(buf[8:12])
}

func writeMessageLen(res []byte, messageLen uint32) {
	binary.BigEndian.PutUint32(res[0:4], messageLen)
}

func writeCorrelationId(res []byte, correlationId uint32) {
	binary.BigEndian.PutUint32(res[4:8], correlationId)
}

func writeErrorCode(res []byte, errCode uint16) {
	binary.BigEndian.PutUint16(res[8:10], errCode)
}

func writeRequestApiKey(res []byte, apiKey uint16) {
	binary.BigEndian.PutUint16(res[10:12], apiKey)
}

func createResponse(len int) []byte {
	return make([]byte, len, 1024)
}

func main() {
	validAPIVersions := []uint16{0, 1, 2, 3, 4}
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
	response := createResponse(n)
	fmt.Printf("response:\t %v \n", []byte(response))
	// NOTE:
	// read 4 bytes, starting from the 8th byte, since correlation_id is 32 bits.
	// first 4 bytes from the buffer are for the message length!!
	correlationId := getCorrelationId(buf)
	writeCorrelationId(response, correlationId)
	fmt.Printf("response:\t %v \n", []byte(response))

	requestApiVersion := getRequestApiVersion(buf)

	if !slices.Contains(validAPIVersions, requestApiVersion) {
		writeErrorCode(response, 35)
	} else {
		writeErrorCode(response, 0)
	}
	fmt.Printf("response:\t %v \n", []byte(response))

	messageLen := getMessageLen(buf)
	writeMessageLen(response, messageLen)
	fmt.Printf("response:\t %v \n", []byte(response))

	apiKey := getRequestApiKey(buf)
	writeRequestApiKey(response, apiKey)
	fmt.Printf("response:\t %v \n", []byte(response))

	conn.Write(response)
	fmt.Println("N: %d", n)
	fmt.Println("Connection handled successfully")
	l.Close()
	conn.Close() // Ensure the connection is closed properly
}
