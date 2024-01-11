package main

import (
	"encoding/binary"
	"fmt"
	"time"
)

func emitRow(startpos uint64, bytes []byte, explanation string, a ...any) {
	BytesWidth := 4
	fmt.Printf("%04x: ", startpos)
	for index, b := range bytes {
		fmt.Printf("%02x ", b)
		if index%BytesWidth == BytesWidth-1 && index != len(bytes)-1 {
			fmt.Println()
			fmt.Print("      ")
		}
	}
	for i := 0; i < (BytesWidth-len(bytes)%BytesWidth)%BytesWidth; i++ {
		fmt.Print("   ")
	}
	formatted := fmt.Sprintf(explanation, a...)
	fmt.Printf("%-20s", formatted)
	fmt.Println()
}

func main() {
	stream := []byte{
		0x1f, 0x8b, 0x08, 0x00, 0xfc, 0x59, 0x96, 0x65, 0x02, 0x03, 0x33, 0x34, 0x32, 0x36, 0xc4, 0x40,
		0x5c, 0x00, 0x5e, 0x96, 0xa9, 0x24, 0x16, 0x00, 0x00, 0x00,
	}
	entry, err := ParseGzip(stream)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", entry)
	emitRow(0, entry.Header.ID[:], "ID (0x1f 0x8b)")
	emitRow(2, []byte{entry.Header.CompressionMethod}, "CM")
	emitRow(3, []byte{entry.Header.Flags}, "FLG")
	emitRow(4, stream[4:8], "MTIME %s", time.Unix(int64(entry.Header.Mtime), 0).UTC().Format(time.RFC3339))
	emitRow(8, []byte{entry.Header.XFL}, "XFL")
	emitRow(9, []byte{entry.Header.OS}, "OS")
	emitRow(entry.DataPos, entry.Data, "Data")
	emitRow(entry.FooterPos, binary.LittleEndian.AppendUint32(nil, entry.Footer.CRC32), "CRC32 0x%08x", entry.Footer.CRC32)
	emitRow(entry.FooterPos+4, binary.LittleEndian.AppendUint32(nil, entry.Footer.Isize), "isize %d", entry.Footer.Isize)
}
