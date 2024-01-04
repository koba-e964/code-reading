package main

// BitWriter is a struct for writing bits to a byte slice.
//
// (8i+j)-th bit is stored in i-th byte at j-th bit.
// Specifically, the least significant bit is stored in the first bit of the byte.
// It can be retrieved with (slice[i] >> j) & 1.
type BitWriter struct {
	// The byte slice to write to.
	bytes []byte
	// The current byte being written to.
	currentByte byte
	// The number of bits written to the current byte.
	numBitsWritten int
}

// Bits appends the given bits represented as a string.
func (b *BitWriter) Bits(s string) {
	for _, c := range s {
		b.Bit(c == '1')
	}
}

// Bit appends the given bit.
func (b *BitWriter) Bit(bit bool) {
	if bit {
		b.currentByte |= 1 << b.numBitsWritten
	}
	b.numBitsWritten++
	if b.numBitsWritten == 8 {
		b.bytes = append(b.bytes, b.currentByte)
		b.currentByte = 0
		b.numBitsWritten = 0
	}
}

// Bytes appends the given bytes.
func (b *BitWriter) Bytes(data []byte) {
	for _, c := range data {
		for i := 0; i < 8; i++ {
			b.Bit((c>>i)&1 == 1)
		}
	}
}

// Int appends the given integer as a length-bit bit string.
func (b *BitWriter) Int(i uint64, length int) {
	for j := 0; j < length; j++ {
		b.Bit((i>>j)&1 == 1)
	}
}

// SkipToByteBoundary skips to the next byte boundary, filling with zeros.
func (b *BitWriter) SkipToByteBoundary() {
	if b.numBitsWritten > 0 {
		b.bytes = append(b.bytes, b.currentByte)
		b.currentByte = 0
		b.numBitsWritten = 0
	}
}

// Emit emits the bit sequence as a byte slice.
func (b *BitWriter) Emit() []byte {
	b.SkipToByteBoundary()
	return b.bytes
}
