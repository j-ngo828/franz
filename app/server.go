package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"slices"
)

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

func writeApiKeysArray(res []byte) {
	// Write number of API keys (1)
	binary.BigEndian.PutUint16(res[10:12], 1)

	// Write APIVersions key (18)
	binary.BigEndian.PutUint16(res[12:14], 18)
	// Min version
	binary.BigEndian.PutUint16(res[14:16], 0)
	// Max version
	binary.BigEndian.PutUint16(res[16:18], 4)
}

func writeTaggedFields(res []byte) {
	// Write varint 0 to indicate no tagged fields
	res[18] = 0
}

func writeThrottleTime(res []byte) {
	binary.BigEndian.PutUint32(res[19:23], 0)
}

func writeFinalTaggedFields(res []byte) {
	// Write varint 0 to indicate no tagged fields
	res[23] = 0
}

func createResponse(len int) []byte {
	return make([]byte, len)
}

func main() {
	validAPIVersions := []uint16{0, 1, 2, 3, 4}
	fmt.Println("Logs from your program will appear here!")

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

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection:", err.Error())
		return
	}
	fmt.Printf("Received %d bytes: %s\n", n, string(buf[:n]))

	response := createResponse(24) // Total size including message length field
	messageLen := uint32(20)       // Size of everything after message length field
	writeMessageLen(response, messageLen)

	correlationId := getCorrelationId(buf)
	writeCorrelationId(response, correlationId)

	requestApiVersion := getRequestApiVersion(buf)
	if !slices.Contains(validAPIVersions, requestApiVersion) {
		writeErrorCode(response, 35)
		conn.Write(response)
		os.Exit(1)
	}

	writeErrorCode(response, 0)
	writeApiKeysArray(response)
	writeTaggedFields(response) // Tagged fields before throttle time
	writeThrottleTime(response)
	writeFinalTaggedFields(response) // Tagged fields at the end

	conn.Write(response)
	l.Close()
	conn.Close()
}
