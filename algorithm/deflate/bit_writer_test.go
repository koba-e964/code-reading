package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitWriter0(t *testing.T) {
	// Test case from https://github.com/Frommi/miniz_oxide/blob/0.7.0/miniz_oxide/src/deflate/core.rs#L2456
	b := new(BitWriter)

	b.Bits("1")        // BFINAL=1
	b.Bits("10")       // BTYPE=01
	b.Bits("00110001") // Lit 0x01 = Direct 0x01
	b.Bits("00110010") // Lit 0x02 = Direct 0x02
	b.Bits("00110011") // Lit 0x03 = Direct 0x03
	b.Bits("00110100") // Lit 0x04 = Direct 0x04
	b.Bits("0000001")  // Lit 0x101 = <length=3>
	b.Bits("00011")    // Dist 0x03 = Distance 0x04
	b.Bits("0000011")  // Lit 0x103 = <length=5>
	b.Bits("00010")    // Dist 0x02 = Distance 0x03
	b.Bits("00110110") // Lit 0x06 = Direct 0x06
	b.Bits("0000011")  // Lit 0x103 = <length=5>
	b.Bits("00100")    // Dist 0x04 = Distance [0x05, 0x07), Extra = 1
	b.Bits("1")        // DExtra: Distance = 0x06
	b.Bits("00110011") // Lit 0x03 = Direct 0x03
	b.Bits("00110010") // Lit 0x02 = Direct 0x02
	b.Bits("00110011") // Lit 0x03 = Direct 0x03
	b.Bits("00110001") // Lit 0x01 = Direct 0x01
	b.Bits("00110010") // Lit 0x02 = Direct 0x02
	b.Bits("00110011") // Lit 0x03 = Direct 0x03
	b.Bits("0000000")  // Lit 0x100 = end

	bs := b.Emit()
	expected := []byte{99, 100, 98, 102, 1, 98, 48, 98, 3, 147, 204, 76, 204, 140, 76, 204, 0}
	assert.Equal(t, expected, bs)
}

func TestBitWriter1(t *testing.T) {
	// Test case from README.md
	// See 2.3 Member Format of <https://datatracker.ietf.org/doc/html/rfc1952>
	b := new(BitWriter)

	b.Bytes([]byte{0x1f, 0x8b})
	b.Bytes([]byte{0x08})                   // CM = 8 deflate
	b.Bytes([]byte{0x00})                   // FLG = 0
	b.Bytes([]byte{0xfc, 0x59, 0x96, 0x65}) // MTIME
	b.Bytes([]byte{0x02})                   // XFL = 2 max compression
	b.Bytes([]byte{0x03})                   // OS = 3 Unix
	// deflate stream starts here
	b.Bits("1")            // BFINAL=1
	b.Bits("10")           // BTYPE=01
	b.Bits("01100001")     // Lit 0x31 = Direct 0x31
	b.Bits("01100010")     // Lit 0x32 = Direct 0x32
	b.Bits("01100011")     // Lit 0x33 = Direct 0x33
	b.Bits("01100001")     // Lit 0x31 = Direct 0x31
	b.Bits("0001100")      // Lit 0x10c = Length [0x11, 0x12)
	b.Bits("0")            // LExtra: Length = 0x11
	b.Bits("00010")        // Dist 0x02 = Distance 0x03
	b.Bits("00111010")     // Lit 0x0a = Direct 0x0a
	b.Bits("0000000")      // Lit 0x100 = end
	b.SkipToByteBoundary() // deflate stream ends here, so pad with 0s
	// deflate stream ends here
	b.Bytes([]byte{0x5e, 0x96, 0xa9, 0x24}) // CRC32
	b.Int(0x16, 32)                         // ISIZE

	bs := b.Emit()
	expected := []byte{
		0x1f, 0x8b, 0x08, 0x00, 0xfc, 0x59, 0x96, 0x65, 0x02, 0x03, 0x33, 0x34, 0x32, 0x36, 0xc4, 0x40,
		0x5c, 0x00, 0x5e, 0x96, 0xa9, 0x24, 0x16, 0x00, 0x00, 0x00,
	}
	assert.Equal(t, expected, bs)
}
