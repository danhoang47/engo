package protocol

import (
	"math"
	"unsafe"
)

const InitialSize = 256 * 1024 // 1MB of buffer with uint32

type CommandBuffer struct {
	Data []uint32
}

func NewCommandBuffer() *CommandBuffer {
	return &CommandBuffer{
		Data: make([]uint32, 0, InitialSize),
	}
}

func NewCommandBufferWithSize(size int) *CommandBuffer {
	return &CommandBuffer{
		Data: make([]uint32, 0, size),
	}
}

func (cb *CommandBuffer) Reset() {
	cb.Data = cb.Data[:0]
}

func (cb *CommandBuffer) GetPtr() unsafe.Pointer {
	if len(cb.Data) == 0 {
		return nil
	}
	return unsafe.Pointer(&cb.Data[0])
}

func (cb *CommandBuffer) GetSize() int {
	return len(cb.Data)
}

func (cb *CommandBuffer) writeHeader(op OpCode, length uint32) {
	header := (length << 8) | uint32(op)
	cb.Data = append(cb.Data, header)
}

func (cb *CommandBuffer) WriteFloat(v float32) {
	cb.Data = append(cb.Data, math.Float32bits(v))
}

func (cb *CommandBuffer) WriteUint(v uint32) {
	cb.Data = append(cb.Data, v)
}

func (cb *CommandBuffer) WriteEof() {
	cb.writeHeader(OpEof, 0)
}

func Color(r, g, b, a uint8) uint32 {
	return (uint32(a) << 24) | (uint32(b) << 16) | (uint32(g) << 8) | uint32(r)
}
