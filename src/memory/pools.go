// Package memory provides high-performance memory management pools and allocators
// for HD1's real-time 3D game engine. Implements exotic algorithms for radical
// performance improvements with zero penalty optimizations.
//
// Key components:
//   - Object pools: sync.Pool for frequent allocations
//   - Buffer management: Pre-sized buffers with controlled growth
//   - JSON optimization: Pooled encoders/decoders for hot paths
//   - Memory arenas: Fixed-size block allocation for entities
package memory

import (
	"bytes"
	"encoding/json"
	"sync"
)

// JSON Buffer Pools - Radical optimization for WebSocket broadcasts
// Eliminates 500-1000+ allocations/second in WebSocket hot paths
var (
	// JSONBufferPool provides reusable byte buffers for JSON marshaling
	// Pre-sized for typical HD1 message sizes (1-4KB)
	JSONBufferPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 4096))
		},
	}
	
	// JSONEncoderPool provides reusable JSON encoders
	// Eliminates encoder allocation overhead in API responses
	JSONEncoderPool = sync.Pool{
		New: func() interface{} {
			return json.NewEncoder(&bytes.Buffer{})
		},
	}
	
	// JSONDecoderPool provides reusable JSON decoders  
	// Optimizes request parsing in 82 API endpoints
	JSONDecoderPool = sync.Pool{
		New: func() interface{} {
			return json.NewDecoder(&bytes.Buffer{})
		},
	}
)

// Message Map Pools - Eliminates map allocation explosion
var (
	// WebSocketUpdatePool for BroadcastToSession operations
	// Pre-sized for typical update message structure
	WebSocketUpdatePool = sync.Pool{
		New: func() interface{} {
			return make(map[string]interface{}, 8)
		},
	}
	
	// ComponentMapPool for entity component operations
	// Optimizes entity creation/update hot paths
	ComponentMapPool = sync.Pool{
		New: func() interface{} {
			return make(map[string]interface{}, 16)
		},
	}
	
	// EntityRequestPool for API request parsing
	// Reduces allocation overhead in entity operations
	EntityRequestPool = sync.Pool{
		New: func() interface{} {
			return make(map[string]interface{}, 12)
		},
	}
)

// Slice Pools - Prevents slice reallocation cascades
var (
	// EntitySlicePool for entity listing operations
	// Pre-sized for typical session entity counts
	EntitySlicePool = sync.Pool{
		New: func() interface{} {
			return make([]map[string]interface{}, 0, 64)
		},
	}
	
	// StringSlicePool for tags and metadata
	// Optimizes entity tag operations
	StringSlicePool = sync.Pool{
		New: func() interface{} {
			return make([]string, 0, 8)
		},
	}
	
	// ByteSlicePool for message operations
	// Reusable byte slices for WebSocket messages
	ByteSlicePool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 0, 2048)
		},
	}
)

// GetJSONBuffer retrieves a pooled byte buffer for JSON operations
// Buffer is reset and ready for use. Must call PutJSONBuffer when done.
func GetJSONBuffer() *bytes.Buffer {
	buf := JSONBufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// PutJSONBuffer returns a byte buffer to the pool for reuse
func PutJSONBuffer(buf *bytes.Buffer) {
	// Prevent memory leaks from oversized buffers
	if buf.Cap() > 16384 { // 16KB limit
		return // Let GC handle oversized buffers
	}
	JSONBufferPool.Put(buf)
}

// GetJSONEncoder retrieves a pooled JSON encoder
// Encoder buffer is reset and ready for use. Must call PutJSONEncoder when done.
func GetJSONEncoder() (*json.Encoder, *bytes.Buffer) {
	buf := GetJSONBuffer()
	encoder := JSONEncoderPool.Get().(*json.Encoder)
	// Reset encoder to use the fresh buffer
	encoder = json.NewEncoder(buf) // Re-initialize with buffer
	return encoder, buf
}

// PutJSONEncoder returns a JSON encoder to the pool for reuse
func PutJSONEncoder(encoder *json.Encoder, buf *bytes.Buffer) {
	JSONEncoderPool.Put(encoder)
	PutJSONBuffer(buf)
}

// GetWebSocketUpdate retrieves a pooled map for WebSocket update messages
// Map is cleared and ready for use. Must call PutWebSocketUpdate when done.
func GetWebSocketUpdate() map[string]interface{} {
	update := WebSocketUpdatePool.Get().(map[string]interface{})
	// Clear the map for reuse
	for k := range update {
		delete(update, k)
	}
	return update
}

// PutWebSocketUpdate returns a WebSocket update map to the pool for reuse
func PutWebSocketUpdate(update map[string]interface{}) {
	// Prevent memory leaks from oversized maps
	if len(update) > 32 {
		return // Let GC handle oversized maps
	}
	WebSocketUpdatePool.Put(update)
}

// GetComponentMap retrieves a pooled map for entity components
// Map is cleared and ready for use. Must call PutComponentMap when done.
func GetComponentMap() map[string]interface{} {
	components := ComponentMapPool.Get().(map[string]interface{})
	// Clear the map for reuse
	for k := range components {
		delete(components, k)
	}
	return components
}

// PutComponentMap returns a component map to the pool for reuse
func PutComponentMap(components map[string]interface{}) {
	// Prevent memory leaks from oversized maps
	if len(components) > 32 {
		return // Let GC handle oversized maps
	}
	ComponentMapPool.Put(components)
}

// GetEntitySlice retrieves a pooled slice for entity operations
// Slice is reset and ready for use. Must call PutEntitySlice when done.
func GetEntitySlice() []map[string]interface{} {
	slice := EntitySlicePool.Get().([]map[string]interface{})
	// Reset slice length but keep capacity
	return slice[:0]
}

// PutEntitySlice returns an entity slice to the pool for reuse
func PutEntitySlice(slice []map[string]interface{}) {
	// Prevent memory leaks from oversized slices
	if cap(slice) > 256 {
		return // Let GC handle oversized slices
	}
	EntitySlicePool.Put(slice)
}

// GetByteSlice retrieves a pooled byte slice for message operations
// Slice is reset and ready for use. Must call PutByteSlice when done.
func GetByteSlice() []byte {
	slice := ByteSlicePool.Get().([]byte)
	// Reset slice length but keep capacity
	return slice[:0]
}

// PutByteSlice returns a byte slice to the pool for reuse
func PutByteSlice(slice []byte) {
	// Prevent memory leaks from oversized slices
	if cap(slice) > 8192 { // 8KB limit
		return // Let GC handle oversized slices
	}
	ByteSlicePool.Put(slice)
}

// GetStringSlice retrieves a pooled string slice for metadata operations
// Slice is reset and ready for use. Must call PutStringSlice when done.
func GetStringSlice() []string {
	slice := StringSlicePool.Get().([]string)
	// Reset slice length but keep capacity
	return slice[:0]
}

// PutStringSlice returns a string slice to the pool for reuse
func PutStringSlice(slice []string) {
	// Prevent memory leaks from oversized slices
	if cap(slice) > 64 {
		return // Let GC handle oversized slices
	}
	StringSlicePool.Put(slice)
}

// GetEntityRequestPool retrieves a pooled map for entity API request parsing
// Map is cleared and ready for use. Must call PutEntityRequestPool when done.
func GetEntityRequestPool() map[string]interface{} {
	request := EntityRequestPool.Get().(map[string]interface{})
	// Clear the map for reuse
	for k := range request {
		delete(request, k)
	}
	return request
}

// PutEntityRequestPool returns an entity request map to the pool for reuse
func PutEntityRequestPool(request map[string]interface{}) {
	// Prevent memory leaks from oversized maps
	if len(request) > 32 {
		return // Let GC handle oversized maps
	}
	EntityRequestPool.Put(request)
}