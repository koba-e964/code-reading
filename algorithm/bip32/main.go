package main

import (
	"encoding/hex"
	"log"

	"github.com/tyler-smith/go-bip32"
)

// https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki#test-vector-3
func testVector3() {
	seedString := "4b381541583be4423346c643850da4b320e46a87ae3d2a4e6da11eba819cd4acba45d239319ac14f863b8d5ab5a0d0c64d2e8a1e7d1457df2e5a3c51c73235be"
	seed, _ := hex.DecodeString(seedString)

	master, _ := bip32.NewMasterKey(seed)
	log.Print(hex.EncodeToString(master.Key))
	if master.Key[0] != 0 {
		panic("master.Key[0] != 0")
	}
}

// https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki#test-vector-4
func testVector4() {
	seedString := "3ddd5602285899a946114506157c7997e5444528f3003f6134712147db19b678"
	seed, _ := hex.DecodeString(seedString)

	master, _ := bip32.NewMasterKey(seed)
	child0, _ := master.NewChildKey(bip32.FirstHardenedChild + 0)
	log.Print(hex.EncodeToString(master.Key))
	log.Print(hex.EncodeToString(child0.Key))
	if child0.Key[0] != 0 {
		panic("child0.Key[0] != 0")
	}
}

func main() {
	testVector3()
	testVector4()
}
