package benchmark

import (
	"testing"

	"github.com/jang3435/franz/pkg/common"
)

// BenchmarkMessageSize measures the performance of message size calculation
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

// BenchmarkMessageClone measures the performance of message cloning
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

// BenchmarkMessageSizeWithDifferentSizes tests message size calculation
// with different message sizes
func BenchmarkMessageSizeWithDifferentSizes(b *testing.B) {
	sizes := []struct {
		name    string
		key     int
		value   int
		headers int
	}{
		{"Small", 10, 100, 1},
		{"Medium", 100, 1024, 5},
		{"Large", 256, 10240, 10},
		{"XLarge", 512, 102400, 20},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			msg := &common.Message{
				Key:     make([]byte, size.key),
				Value:   make([]byte, size.value),
				Headers: make(map[string][]byte),
			}

			// Add headers
			for i := 0; i < size.headers; i++ {
				keyStr := "header" + string(rune(i))
				msg.Headers[keyStr] = make([]byte, 16)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = msg.Size()
			}
		})
	}
}

// BenchmarkCompressionTypeString measures performance of compression type string conversion
func BenchmarkCompressionTypeString(b *testing.B) {
	compressionTypes := []common.CompressionType{
		common.CompressionNone,
		common.CompressionGzip,
		common.CompressionSnappy,
		common.CompressionLZ4,
		common.CompressionZstd,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, ct := range compressionTypes {
			_ = ct.String()
		}
	}
}

// BenchmarkParseCompressionType measures performance of parsing compression types
func BenchmarkParseCompressionType(b *testing.B) {
	inputs := []string{
		"none",
		"gzip",
		"snappy",
		"lz4",
		"zstd",
		"invalid",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, input := range inputs {
			_ = common.ParseCompressionType(input)
		}
	}
}

// BenchmarkMessageBatchCreation measures performance of creating message batches
func BenchmarkMessageBatchCreation(b *testing.B) {
	// Pre-allocate messages for batch
	messages := make([]*common.Message, 100)
	for i := range messages {
		messages[i] = &common.Message{
			Key:   []byte("key" + string(rune(i))),
			Value: []byte("value" + string(rune(i))),
			Headers: map[string][]byte{
				"id": []byte(string(rune(i))),
			},
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		batch := &common.MessageBatch{
			Messages:    messages,
			FirstOffset: int64(i * 100),
			LastOffset:  int64((i+1)*100 - 1),
		}
		_ = batch
	}
}

// BenchmarkMemoryAllocation tests memory allocation patterns
func BenchmarkMemoryAllocation(b *testing.B) {
	b.Run("NewMessage", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			msg := &common.Message{
				Key:     make([]byte, 32),
				Value:   make([]byte, 256),
				Headers: make(map[string][]byte),
			}
			_ = msg
		}
	})

	b.Run("MessageWithHeaders", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			msg := &common.Message{
				Key:   make([]byte, 32),
				Value: make([]byte, 256),
				Headers: map[string][]byte{
					"content-type": []byte("application/json"),
					"user-id":      []byte("12345"),
					"trace-id":     []byte("abc123"),
				},
			}
			_ = msg
		}
	})
}
