package unit

import (
	"strings"
	"testing"
	"time"

	"github.com/jang3435/franz/pkg/common"
)

func TestMessageCreation(t *testing.T) {
	key := []byte("test-key")
	value := []byte("test-value")
	headers := map[string][]byte{
		"content-type": []byte("application/json"),
		"user-id":      []byte("12345"),
	}

	msg := &common.Message{
		Key:       key,
		Value:     value,
		Timestamp: time.Now(),
		Headers:   headers,
		Partition: 1,
		Offset:    42,
	}

	// Test message fields
	if !strings.EqualFold(string(msg.Key), string(key)) {
		t.Errorf("Expected key %s, got %s", key, msg.Key)
	}

	if !strings.EqualFold(string(msg.Value), string(value)) {
		t.Errorf("Expected value %s, got %s", value, msg.Value)
	}

	if msg.Partition != 1 {
		t.Errorf("Expected partition 1, got %d", msg.Partition)
	}

	if msg.Offset != 42 {
		t.Errorf("Expected offset 42, got %d", msg.Offset)
	}

	// Test headers
	contentType, exists := msg.Headers["content-type"]
	if !exists {
		t.Error("Expected content-type header to exist")
	}
	if string(contentType) != "application/json" {
		t.Errorf("Expected content-type 'application/json', got %s", contentType)
	}

	userID, exists := msg.Headers["user-id"]
	if !exists {
		t.Error("Expected user-id header to exist")
	}
	if string(userID) != "12345" {
		t.Errorf("Expected user-id '12345', got %s", userID)
	}
}

func TestMessageSize(t *testing.T) {
	msg := &common.Message{
		Key:     []byte("key"),
		Value:   []byte("value"),
		Headers: map[string][]byte{"header": []byte("value")},
	}

	size := msg.Size()

	// Size should include key, value, timestamp (8), partition (4), offset (8), and headers
	expectedMinSize := len(msg.Key) + len(msg.Value) + 8 + 4 + 8 + len("header") + len("value")
	if size < int32(expectedMinSize) {
		t.Errorf("Expected size >= %d, got %d", expectedMinSize, size)
	}
}

func TestMessageClone(t *testing.T) {
	original := &common.Message{
		Key:       []byte("original-key"),
		Value:     []byte("original-value"),
		Headers:   map[string][]byte{"original": []byte("header")},
		Partition: 1,
		Offset:    100,
	}

	clone := original.Clone()

	// Test that values are equal
	if string(clone.Key) != string(original.Key) {
		t.Error("Cloned key doesn't match original")
	}
	if string(clone.Value) != string(original.Value) {
		t.Error("Cloned value doesn't match original")
	}
	if clone.Partition != original.Partition {
		t.Error("Cloned partition doesn't match original")
	}
	if clone.Offset != original.Offset {
		t.Error("Cloned offset doesn't match original")
	}

	// Test that headers are copied
	if string(clone.Headers["original"]) != "header" {
		t.Error("Cloned headers don't match original")
	}

	// Test that they're not the same slice (deep copy)
	original.Key[0] = 'X'
	if clone.Key[0] == 'X' {
		t.Error("Clone shares the same key slice as original")
	}

	original.Headers["original"][0] = 'Y'
	if clone.Headers["original"][0] == 'Y' {
		t.Error("Clone shares the same header slice as original")
	}
}

func TestCompressionTypeString(t *testing.T) {
	tests := []struct {
		compression common.CompressionType
		expected    string
	}{
		{common.CompressionNone, "none"},
		{common.CompressionGzip, "gzip"},
		{common.CompressionSnappy, "snappy"},
		{common.CompressionLZ4, "lz4"},
		{common.CompressionZstd, "zstd"},
		{common.CompressionType(99), "unknown"}, // Invalid type
	}

	for _, test := range tests {
		if test.compression.String() != test.expected {
			t.Errorf("Expected %s, got %s for compression type %d",
				test.expected, test.compression.String(), test.compression)
		}
	}
}

func TestParseCompressionType(t *testing.T) {
	tests := []struct {
		input    string
		expected common.CompressionType
	}{
		{"none", common.CompressionNone},
		{"gzip", common.CompressionGzip},
		{"snappy", common.CompressionSnappy},
		{"lz4", common.CompressionLZ4},
		{"zstd", common.CompressionZstd},
		{"invalid", common.CompressionNone}, // Default to none for invalid input
	}

	for _, test := range tests {
		if common.ParseCompressionType(test.input) != test.expected {
			t.Errorf("Expected %d for input %s, got %d",
				test.expected, test.input, common.ParseCompressionType(test.input))
		}
	}
}

func BenchmarkMessageSize(b *testing.B) {
	msg := &common.Message{
		Key:   make([]byte, 32),
		Value: make([]byte, 1024),
		Headers: map[string][]byte{
			"header1": make([]byte, 16),
			"header2": make([]byte, 16),
			"header3": make([]byte, 16),
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = msg.Size()
	}
}

func BenchmarkMessageClone(b *testing.B) {
	msg := &common.Message{
		Key:   make([]byte, 32),
		Value: make([]byte, 1024),
		Headers: map[string][]byte{
			"header1": make([]byte, 16),
			"header2": make([]byte, 16),
			"header3": make([]byte, 16),
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = msg.Clone()
	}
}
