package main

import (
	"crypto/rand"
	"fmt"
	"log"

	bip32 "github.com/koba-e964/bip32-typesafe"
)

func main() {
	// Generate random 32 bytes
	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		panic(err)
	}

	master := bip32.NewMasterKey(seed)
	log.Println(master.PrivateKey())
	child0, err := master.NewChildKey(0) // master/0
	if err != nil {
		panic(err)
	}
	fmt.Println("master/0 =", child0.B58Serialize())
	childH0, err := master.NewChildKey(bip32.FirstHardenedChildIndex + 0) // master/0_H
	if err != nil {
		panic(err)
	}
	fmt.Println("master/0_H =", childH0.B58Serialize())
}
