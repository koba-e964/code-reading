package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"

	"github.com/koba-e964/code-reading/algorithm/crc-tampering/char2"
)

func tamper(arrayLen int, modifiableIndices []int, targetCRC32 uint32) ([]byte, error) {
	n := len(modifiableIndices)
	rows := make([]char2.GF232, 8*n)
	for i, pos := range modifiableIndices {
		rows[8*i+7].Weight32(arrayLen - pos + 3)
		for j := 6; j >= 0; j-- {
			rows[8*i+j].MulX(rows[8*i+j+1])
		}
	}
	// Is targetCRC32 a linear combination of rows?
	rowsBytes := make([][]byte, 8*n)
	for i := range rows {
		rowsBytes[i] = make([]byte, 4)
		binary.LittleEndian.PutUint32(rowsBytes[i], uint32(rows[i]))
	}
	targetBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(targetBytes, targetCRC32)
	sol, ok := char2.SolveLinear(rowsBytes, targetBytes)
	if !ok {
		return nil, errors.New("no solution")
	}

	// Construct the return value.
	ret := make([]byte, n)
	for i, v := range sol {
		if v {
			ret[i/8] |= 1 << (i % 8)
		}
	}
	return ret, nil
}

func main() {
	data := []byte("hello, world!")
	modifiableIndices := []int{1, 2, 3, 4}
	targetCRC32 := uint32(0x1111_1111) ^ crc32.ChecksumIEEE(data)
	sol, err := tamper(len(data), modifiableIndices, targetCRC32)
	if err != nil {
		panic(err)
	}
	modified := make([]byte, len(data))
	copy(modified, data)
	for i, pos := range modifiableIndices {
		modified[pos] ^= sol[i]
	}
	fmt.Printf("crc32: %#v\n", crc32.ChecksumIEEE(modified))
	fmt.Printf("crc32 diff: %#v\n", crc32.ChecksumIEEE(modified)^crc32.ChecksumIEEE(data))
}
