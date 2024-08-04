package util

import "bytes"

type RingBuffer struct {
	buffer *bytes.Buffer
	size   int
}

func NewLimitedBuffer(size int) *RingBuffer {
	return &RingBuffer{
		buffer: bytes.NewBuffer(make([]byte, size)),
		size:   size,
	}
}

func (lb *RingBuffer) Write(p []byte) (int, error) {
	if len(p) > lb.size {
		// If the input data is larger than the buffer size, write only the last part of it
		p = p[len(p)-lb.size:]
	}

	if lb.buffer.Len()+len(p) > lb.size {
		// Remove the oldest data to make space
		overflow := lb.buffer.Len() + len(p) - lb.size
		lb.buffer.Next(overflow)
	}

	return lb.buffer.Write(p)
}

func (lb *RingBuffer) Read(p []byte) (int, error) {
	return lb.buffer.Read(p)
}
