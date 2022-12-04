// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package core

// BoundedBuffer keeps [size] entries of type [K] in a buffer and calls
// [callback] on any item that is evicted. This is typically used for
// dereferencing old roots during block processing.
type BoundedBuffer[K any] struct {
	lastPos  int
	size     int
	callback func(K)
	buffer   []K

	cycled bool
}

// NewBoundedBuffer creates a new [BoundedBuffer].
func NewBoundedBuffer[K any](size int, callback func(K)) *BoundedBuffer[K] {
	return &BoundedBuffer[K]{
		size:     size,
		callback: callback,
		buffer:   make([]K, size),
	}
}

// Insert adds a new value to the buffer. If the buffer is full, the
// oldest value will be evicted and [callback] will be invoked.
func (b *BoundedBuffer[K]) Insert(h K) {
	nextPos := (b.lastPos + 1) % b.size // the first item added to the buffer will be at position 1
	if b.cycled {
		// We ensure we have cycled through the buffer once before invoking the
		// [callback] to ensure we don't call it with unset values.
		b.callback(b.buffer[nextPos])
	}
	b.buffer[nextPos] = h
	b.lastPos = nextPos
	if !b.cycled && nextPos == 0 {
		// Set [cycled] once we are back the 0th element (recall we start at index
		// 1, so we need to wait until the 0th element has been processed at least
		// once.
		b.cycled = true
	}
}

// Last retrieves the last item added to the buffer.
// If no items have been added to the buffer, Last returns an empty hash.
func (b *BoundedBuffer[K]) Last() K {
	return b.buffer[b.lastPos]
}
