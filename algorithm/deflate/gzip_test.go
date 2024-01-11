package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGzip0(t *testing.T) {
	stream := []byte{
		0x1f, 0x8b, 0x08, 0x00, 0xfc, 0x59, 0x96, 0x65, 0x02, 0x03, 0x33, 0x34, 0x32, 0x36, 0xc4, 0x40,
		0x5c, 0x00, 0x5e, 0x96, 0xa9, 0x24, 0x16, 0x00, 0x00, 0x00,
	}
	entry, err := ParseGzip(stream)
	if assert.NoError(t, err) {
		assert.Equal(t, GzipEntry{
			Header: GzipHeader{
				ID:                [2]byte{0x1f, 0x8b},
				CompressionMethod: 8,
				Flags:             0,
				Mtime:             0x659659fc,
				XFL:               2,
				OS:                3,
			},
			DataPos: 10,
			Data: []byte{
				0x33, 0x34, 0x32, 0x36, 0xc4, 0x40, 0x5c, 0x00,
			},
			FooterPos: 18,
			Footer: GzipFooter{
				CRC32: 0x24a9965e,
				Isize: 0x16,
			},
		}, *entry)
	}
}
