package main

import (
	"fmt"

	"github.com/bwesterb/go-ristretto/edwards25519"
)

func main() {
	i := new(edwards25519.FieldElement).SetI()
	fmt.Println("i =", i.Bytes())
	fmt.Printf("i = 0x%s\n", i.BigInt().Text(16))

	base := new(edwards25519.ExtendedPoint).SetBase()
	fmt.Println("base =", base.Ristretto())
	fmt.Printf("base.y = 0x%s\n", base.Y.BigInt().Text(16))
}
