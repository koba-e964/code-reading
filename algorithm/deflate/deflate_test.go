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
			{42, 2, `LExtra "10" -> 1 (calculated: 24)`},
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

func TestParseDeflateBtype10_0(t *testing.T) {
	stream := []byte{0x4d, 0x8a, 0xc1, 0x0d,
		0x00, 0x20, 0x0c, 0x02,
		0xff, 0x4e, 0x03, 0xfb,
		0x2f, 0xa7, 0x70, 0xd5,
		0x48, 0x48, 0xe8, 0x41,
		0x6d, 0x2b, 0x72, 0xd3,
		0x1b,
	}
	deflates, err := ParseDeflate(stream)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(deflates))
		actual := deflates[0].Show()
		expected := []RangeDescription{
			{StartPos: 0, Length: 3, Description: "deflate header"},
			{StartPos: 0, Length: 1, Description: "BFINAL 1"},
			{StartPos: 1, Length: 2, Description: "BTYPE 2"},
			{StartPos: 3, Length: 0xc2, Description: "deflate content"},
			{StartPos: 3, Length: 5, Description: "NumLitLen 9"},
			{StartPos: 0x8, Length: 5, Description: "NumDist 10"},
			{StartPos: 0xd, Length: 4, Description: "NumCode 12"},
			{StartPos: 0x11, Length: 3, Description: "CodeLen[16] 0"},
			{StartPos: 0x14, Length: 3, Description: "CodeLen[17] 4"},
			{StartPos: 0x17, Length: 3, Description: "CodeLen[18] 3"},
			{StartPos: 0x1a, Length: 3, Description: "CodeLen[0] 3"},
			{StartPos: 0x1d, Length: 3, Description: "CodeLen[8] 0"},
			{StartPos: 0x20, Length: 3, Description: "CodeLen[7] 0"},
			{StartPos: 0x23, Length: 3, Description: "CodeLen[9] 0"},
			{StartPos: 0x26, Length: 3, Description: "CodeLen[6] 0"},
			{StartPos: 0x29, Length: 3, Description: "CodeLen[10] 0"},
			{StartPos: 0x2c, Length: 3, Description: "CodeLen[5] 2"},
			{StartPos: 0x2f, Length: 3, Description: "CodeLen[11] 0"},
			{StartPos: 0x32, Length: 3, Description: "CodeLen[4] 3"},
			{StartPos: 0x35, Length: 3, Description: "CodeLen[12] 0"},
			{StartPos: 0x38, Length: 3, Description: "CodeLen[3] 2"},
			{StartPos: 0x3b, Length: 3, Description: "CodeLen[13] 0"},
			{StartPos: 0x3e, Length: 3, Description: "CodeLen[2] 4"},
			{StartPos: 0x41, Length: 4, Description: `Code[0] "1111" -> 17`},
			{StartPos: 0x45, Length: 3, Description: `Repeat 0 "111" -> 7 (calculated: 10)`},
			{StartPos: 0x48, Length: 2, Description: `Code[10] "01" -> 5`},
			{StartPos: 0x4a, Length: 3, Description: `Code[11] "110" -> 18`},
			{StartPos: 0x4d, Length: 7, Description: `Repeat 0 "0101100" -> 26 (calculated: 37)`},
			{StartPos: 0x54, Length: 2, Description: `Code[48] "00" -> 3`},
			{StartPos: 0x56, Length: 2, Description: `Code[49] "00" -> 3`},
			{StartPos: 0x58, Length: 3, Description: `Code[50] "110" -> 18`},
			{StartPos: 0x5b, Length: 7, Description: `Repeat 0 "1111111" -> 127 (calculated: 138)`},
			{StartPos: 0x62, Length: 3, Description: `Code[188] "110" -> 18`},
			{StartPos: 0x65, Length: 7, Description: `Repeat 0 "1001110" -> 57 (calculated: 68)`},
			{StartPos: 0x6c, Length: 2, Description: `Code[256] "01" -> 5`},
			{StartPos: 0x6e, Length: 2, Description: `Code[257] "01" -> 5`},
			{StartPos: 0x70, Length: 2, Description: `Code[258] "00" -> 3`},
			{StartPos: 0x72, Length: 2, Description: `Code[259] "00" -> 3`},
			{StartPos: 0x74, Length: 4, Description: `Code[260] "1110" -> 2`},
			{StartPos: 0x78, Length: 3, Description: `Code[261] "101" -> 4`},
			{StartPos: 0x7b, Length: 2, Description: `Code[262] "01" -> 5`},
			{StartPos: 0x7d, Length: 2, Description: `Code[263] "01" -> 5`},
			{StartPos: 0x7f, Length: 3, Description: `Code[264] "100" -> 0`},
			{StartPos: 0x82, Length: 2, Description: `Code[265] "01" -> 5`},
			{StartPos: 0x84, Length: 2, Description: `Code[266] "00" -> 3`},
			{StartPos: 0x86, Length: 3, Description: `Code[267] "100" -> 0`},
			{StartPos: 0x89, Length: 2, Description: `Code[268] "00" -> 3`},
			{StartPos: 0x8b, Length: 3, Description: `Code[269] "100" -> 0`},
			{StartPos: 0x8e, Length: 3, Description: `Code[270] "100" -> 0`},
			{StartPos: 0x91, Length: 2, Description: `Code[271] "00" -> 3`},
			{StartPos: 0x93, Length: 3, Description: `Code[272] "101" -> 4`},
			{StartPos: 0x96, Length: 4, Description: `Code[273] "1110" -> 2`},
			{StartPos: 0x9a, Length: 2, Description: `Code[274] "00" -> 3`},
			{StartPos: 0x9c, Length: 2, Description: `Code[275] "00" -> 3`},
			{StartPos: 0x9e, Length: 3, Description: `Code[276] "101" -> 4`},
			{StartPos: 0xa1, Length: 3, Description: `Literal "011" -> 49`},
			{StartPos: 0xa4, Length: 3, Description: `Literal "011" -> 49`},
			{StartPos: 0xa7, Length: 3, Description: `Literal "011" -> 49`},
			{StartPos: 0xaa, Length: 3, Description: `Literal "010" -> 48`},
			{StartPos: 0xad, Length: 3, Description: `Length "100" -> 258 (calculated: 4)`},
			{StartPos: 0xb0, Length: 3, Description: `Distance "010" -> 0 (calculated: 1)`},
			{StartPos: 0xb3, Length: 3, Description: `Literal "011" -> 49`},
			{StartPos: 0xb6, Length: 3, Description: `Length "101" -> 259 (calculated: 5)`},
			{StartPos: 0xb9, Length: 3, Description: `Distance "100" -> 5 (calculated: 7)`},
			{StartPos: 0xbc, Length: 1, Description: `DExtra "1" -> 1 (calculated: 8)`},
			{StartPos: 0xbd, Length: 3, Description: `Literal "011" -> 49`},
			{StartPos: 0xc0, Length: 5, Description: `End "11011" -> 256`},
		}
		assert.Equal(t, expected, actual)
	}
}
