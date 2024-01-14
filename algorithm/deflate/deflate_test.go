package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDeflateBtype00(t *testing.T) {
	// $ echo $"0: 1f 8b 08 00 fc 59 96 65 02 03 01 01 00 fe ff 00 \n10: 8def02d2 01 00 00 00" | xxd -r >data.bin
	// $ gunzip data.bin -c | hexdump -C
	// 00000000  00                                                |.|
	// 00000001
	// The value 0xd202ef8d is the CRC32 of the uncompressed data:
	// $ crc32 <(echo "0:00" | xxd -r)
	// d202ef8d
	stream := []byte{0x01, 0x01, 0x00, 0xfe, 0xff, 0x00}
	deflates, err := ParseDeflate(stream)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(deflates))
		expected := []RangeDescription{
			{0, 3, "deflate header"},
			{0, 1, "BFINAL 1"},
			{1, 2, "BTYPE 0"},
			// Anything up to the next byte boundary is ignored.
			{8, 40, "deflate content"},
			{8, 16, "Length"},
			{24, 16, "NLength (inverted)"},
			{40, 8, "literal 1bytes"},
		}
		assert.Equal(t, expected, deflates[0].Show())
	}
}

func TestParseDeflateBtype01(t *testing.T) {
	stream := []byte{0x33, 0x34, 0x32, 0x36, 0xc4, 0x40, 0x5c, 0x00}
	deflates, err := ParseDeflate(stream)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(deflates))
		actual := deflates[0].Show()
		expected := []RangeDescription{
			{0, 3, "deflate header"},
			{0, 1, "BFINAL 1"},
			{1, 2, "BTYPE 1"},
			{3, 60, "deflate content"},
			{3, 8, `Literal "01100001" -> 49`},
			{11, 8, `Literal "01100010" -> 50`},
			{19, 8, `Literal "01100011" -> 51`},
			{27, 8, `Literal "01100001" -> 49`},
			{35, 7, `Length "0001100" -> 268 (calculated: 17)`},
			{42, 1, `LExtra "0" -> 0 (calculated: 17)`},
			{43, 5, `Distance "00010" -> 2 (calculated: 3)`},
			{48, 8, `Literal "00111010" -> 10`},
			{56, 7, `End "0000000" -> 256`},
		}
		assert.Equal(t, expected, actual)
	}
}
