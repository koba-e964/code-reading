package main

import "github.com/koba-e964/code-reading/algorithm/crc-tampering/char2"

func tamper(arrayLen int, modifiableIndices []int, targetCRC32 uint32) []byte {
	n := len(modifiableIndices)
	rows := make([]char2.GF232, 8*n)
	for i, pos := range modifiableIndices {
		rows[8*i+7].Weight32(arrayLen - pos + 3)
		for j := 6; j >= 0; j-- {
			rows[8*i+j].MulX(rows[8*i+j+1])
		}
	}
	// Is targetCRC32 a linear combination of rows?
	return nil
}

func main() {

}
