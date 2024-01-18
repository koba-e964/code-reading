package gzip

// BitReader is a struct for reading bits from a byte slice.
//
// (8i+j)-th bit is stored in i-th byte at j-th bit.
// Specifically, the least significant bit is stored in the first bit of the byte.
// It can be retrieved with (slice[i] >> j) & 1.
type BitReader struct {
	// The byte slice to read from.
	bytes []byte
	// The number of bits read from the current byte.
	numBitsRead int
	// The current index in the byte slice.
	index int
}

// NewBitReader creates a new BitReader.
func NewBitReader(bytes []byte) *BitReader {
	return &BitReader{bytes: bytes}
}

// Bit reads a bit.
func (b *BitReader) Bit() byte {
	result := (b.bytes[b.index] >> b.numBitsRead) & 1
	b.numBitsRead++
	if b.numBitsRead == 8 {
		b.numBitsRead = 0
		b.index++
	}
	return result
}

// Bytes reads the given number of bytes.
func (b *BitReader) Bytes(length int) []byte {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = byte(b.Int(8))
	}
	return result
}

// Int reads an integer of the given length.
func (b *BitReader) Int(length int) uint64 {
	var result uint64
	for i := 0; i < length; i++ {
		result |= uint64(b.Bit()) << i
	}
	return result
}

// IntBE reads an integer of the given length, regarded as big endian.
func (b *BitReader) IntBE(length int) uint64 {
	var result uint64
	for i := length - 1; i >= 0; i-- {
		result |= uint64(b.Bit()) << i
	}
	return result
}

func (b *BitReader) SkipToByteBoundary() {
	if b.numBitsRead > 0 {
		b.numBitsRead = 0
		b.index++
	}
}

// Remaining returns the number of bits remaining to be read.
func (b *BitReader) Remaining() int {
	return (len(b.bytes)-b.index)*8 - b.numBitsRead
}

// Position returns the current position in bits.
func (b *BitReader) Position() uint64 {
	return uint64(b.index)*8 + uint64(b.numBitsRead)
}
