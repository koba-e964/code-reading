package main

import (
	"encoding/binary"
	"hash/crc32"
	"io"
	"os"
	"time"
)

const blockSize = 65535

// Creates a new gzip file without compression.
func main() {
	file := os.Stdin
	if len(os.Args) > 1 {
		currentFile, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		file = currentFile
	}
	defer file.Close()
	stream, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	header := []byte{
		0x1f, 0x8b,
		0x08,                   // CM = 8 deflate
		0x00,                   // FLG = 0
		0xff, 0xff, 0xff, 0xff, // MTIME
		0x00, // XFL = 0
		0x03, // OS = Unix
	}
	now := uint32(time.Now().Unix())
	binary.LittleEndian.PutUint32(header[4:8], now)
	if _, err := os.Stdout.Write(header); err != nil {
		panic(err)
	}
	for i := 0; i < (len(stream)+blockSize-1)/blockSize; i++ {
		block := stream[i*blockSize:]
		if len(block) > blockSize {
			block = block[:blockSize]
		}
		var final byte
		if i == (len(stream)+blockSize-1)/blockSize-1 {
			final = 1
		}
		length := uint16(len(block))
		nlength := uint16(^length)
		if _, err := os.Stdout.Write([]byte{
			final, // BFINAL = final, BTYPE = 00
			byte(length), byte(length >> 8),
			byte(nlength), byte(nlength >> 8),
		}); err != nil {
			panic(err)
		}
		if _, err := os.Stdout.Write(block); err != nil {
			panic(err)
		}
	}

	footer := []byte{
		0x00, 0x00, 0x00, 0x00, // CRC32
		0x00, 0x00, 0x00, 0x00, // ISIZE
	}
	binary.LittleEndian.PutUint32(footer[0:4], crc32.ChecksumIEEE(stream))
	binary.LittleEndian.PutUint32(footer[4:8], uint32(len(stream)))
	if _, err := os.Stdout.Write(footer); err != nil {
		panic(err)
	}
}
