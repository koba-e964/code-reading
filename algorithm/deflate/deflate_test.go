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

func TestParseDeflateBtype01_0(t *testing.T) {
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
func TestParseDeflateBtype01_1(t *testing.T) {
	// $ echo 1231231231231231231231231231 | gzip | ~/srcview/infgen/infgen -idd
	// ! infgen 3.2 output
	// !
	// time 1705306760         ! [UTC Mon Jan 15 08:19:20 2024]
	// gzip
	// !
	// last                    ! 1
	// fixed                   ! 01
	// literal '1              ! 10000110
	// literal '2              ! 01000110
	// literal '3              ! 11000110
	// literal '1              ! 10000110
	// match 24 3              ! 01000 01 0111000
	// literal 10              ! 01011100
	// end                     ! 0000000
	// !
	// crc
	// length
	// $ echo 1231231231231231231231231231 | gzip | hexdump -C
	// 00000000  1f 8b 08 00 8f ea a4 65  00 03 33 34 32 36 c4 85  |.......e..3426..|
	// 00000010  b8 00 6e e7 ad d4 1d 00  00 00                    |..n.......|
	// 0000001a
	stream := []byte{0x33, 0x34, 0x32, 0x36, 0xc4, 0x85, 0xb8, 0x00}
	deflates, err := ParseDeflate(stream)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(deflates))
		actual := deflates[0].Show()
		expected := []RangeDescription{
			{0, 3, "deflate header"},
			{0, 1, "BFINAL 1"},
			{1, 2, "BTYPE 1"},
			{3, 61, "deflate content"},
			{3, 8, `Literal "01100001" -> 49`},
			{11, 8, `Literal "01100010" -> 50`},
			{19, 8, `Literal "01100011" -> 51`},
			{27, 8, `Literal "01100001" -> 49`},
			{35, 7, `Length "0001110" -> 270 (calculated: 23)`},
			{42, 2, `LExtra "01" -> 1 (calculated: 24)`},
			{44, 5, `Distance "00010" -> 2 (calculated: 3)`},
			{49, 8, `Literal "00111010" -> 10`},
			{57, 7, `End "0000000" -> 256`},
		}
		assert.Equal(t, expected, actual)
	}
}

func TestParseDeflateBtype01_2(t *testing.T) {
	// $ echo "0: ff fe" | xxd -r | gzip | ~/srcview/infgen/infgen -idd
	// ! infgen 3.2 output
	// !
	// time 1705361531         ! [UTC Mon Jan 15 23:32:11 2024]
	// gzip
	// !
	// last                    ! 1
	// fixed                   ! 01
	// literal 255             ! 111111111
	// literal 254             ! 011111111
	// end                     ! 0000000
	//                         ! 0000
	// !
	// crc
	// length
	// $ echo "0: ff fe" | xxd -r | gzip | hexdump -C
	// 00000000  1f 8b 08 00 93 c0 a5 65  00 03 fb ff 0f 00 96 30  |.......e.......0|
	// 00000010  f8 88 02 00 00 00                                 |......|
	// 00000016
	stream := []byte{0xfb, 0xff, 0x0f, 0x00}
	deflates, err := ParseDeflate(stream)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(deflates))
		actual := deflates[0].Show()
		expected := []RangeDescription{
			{0, 3, "deflate header"},
			{0, 1, "BFINAL 1"},
			{1, 2, "BTYPE 1"},
			{3, 25, "deflate content"},
			{3, 9, `Literal "111111111" -> 255`},
			{12, 9, `Literal "111111110" -> 254`},
			{21, 7, `End "0000000" -> 256`},
		}
		assert.Equal(t, expected, actual)
	}
}
