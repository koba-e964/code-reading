package char2

import (
	"hash/crc32"
	"testing"
)

func crc32WithGF232(b []byte) uint32 {
	var g, tmp GF232
	g.Mul(GF232(0xffff_ffff), *tmp.Weight32(len(b)))
	for i := 0; i < len(b); i++ {
		tmp.Mul(*tmp.Weight32(len(b) - i), GF232(uint32(b[i])))
		g.Add(g, tmp)
	}
	return uint32(g) ^ 0xffff_ffff
}

func TestWeight32(t *testing.T) {
	tests := []struct {
		pos  int
		want GF232
	}{
		{0, 0x8000_0000},
		{1, 0x80_0000},
		{2, 0x8000},
		{3, 0x80},
		{4, crc32.IEEE},
	}
	var got GF232
	for _, tt := range tests {
		if got.Weight32(tt.pos); got != tt.want {
			t.Errorf("weight32(%d) = %#v; want %#v", tt.pos, got, tt.want)
		}
	}
}

func TestCrc32(t *testing.T) {
	tests := [][]byte{
		{},
		{'a'},
		[]byte("hello, world"),
		[]byte("hello, world!"),
		[]byte("test"),
	}
	for _, b := range tests {
		left := crc32WithGF232(b)
		right := crc32.ChecksumIEEE(b)
		if left != right {
			t.Errorf("crc32WithGF232(%q) = %#v != crc32.ChecksumIEEE(%q) = %#v", b, left, b, right)
		}
	}
}
