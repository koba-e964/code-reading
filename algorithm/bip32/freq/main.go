package main

import (
	"log"

	bip39 "github.com/tyler-smith/go-bip39"
)

func main() {
	list := bip39.GetWordList()
	freq := [10]int{}
	for _, word := range list {
		freq[len(word)]++
	}
	for i, f := range freq {
		if f > 0 {
			log.Printf("%d: %d", i, f)
		}
	}
}
