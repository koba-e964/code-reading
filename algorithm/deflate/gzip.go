package gzip

import (
	"encoding/binary"
	"fmt"
)

type GzipEntry struct {
	Header    GzipHeader
	DataPos   uint64
	Data      []byte
	FooterPos uint64
	Footer    GzipFooter
}

type GzipHeader struct {
	ID                [2]byte
	CompressionMethod byte
	Flags             byte
	Mtime             uint32
	XFL               byte
	OS                byte

	Extra    []byte
	Filename []byte
}

type GzipFooter struct {
	CRC32 uint32
	Isize uint32
}

// ParseGzip parses a gzip stream that contains exactly one member.
func ParseGzip(stream []byte) (*GzipEntry, error) {
	var entry GzipEntry
	entry.Header.ID[0] = stream[0]
	entry.Header.ID[1] = stream[1]
	if entry.Header.ID != [2]byte{0x1f, 0x8b} {
		return nil, fmt.Errorf("invalid gzip header")
	}
	entry.Header.CompressionMethod = stream[2]
	entry.Header.Flags = stream[3]
	entry.Header.Mtime = binary.LittleEndian.Uint32(stream[4:8])
	entry.Header.XFL = stream[8]
	entry.Header.OS = stream[9]
	pos := 10
	if entry.Header.Flags&0x04 != 0 {
		// Extra field
		length := int(binary.LittleEndian.Uint16(stream[10:12]))
		entry.Header.Extra = stream[12 : 12+length]
		pos = 12 + length
	}
	if entry.Header.Flags&0x08 != 0 {
		// Filename
		for stream[pos] != 0 {
			entry.Header.Filename = append(entry.Header.Filename, stream[pos])
			pos++
		}
		pos++
	}
	if entry.Header.Flags&0x10 != 0 {
		// Comment
		for stream[pos] != 0 {
			pos++
		}
		pos++
	}

	if entry.Header.Flags&0x02 != 0 {
		// CRC16
		pos += 2
	}
	entry.DataPos = uint64(pos)

	// Take last 8 bytes as footer
	footerStart := len(stream) - 8
	entry.FooterPos = uint64(footerStart)
	entry.Footer.CRC32 = binary.LittleEndian.Uint32(stream[footerStart : footerStart+4])
	entry.Footer.Isize = binary.LittleEndian.Uint32(stream[footerStart+4 : footerStart+8])
	entry.Data = stream[pos:footerStart]

	return &entry, nil
}
