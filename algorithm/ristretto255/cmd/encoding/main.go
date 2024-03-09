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

	s := new(edwards25519.ExtendedPoint).SetTorsion1()
	fmt.Println("s1 =", s.Ristretto(), s.String())

	s = new(edwards25519.ExtendedPoint).SetTorsion2()
	fmt.Println("s2 =", s.Ristretto(), s.String())

	s = new(edwards25519.ExtendedPoint).SetTorsion3()
	fmt.Println("s3 =", s.Ristretto())

	s = new(edwards25519.ExtendedPoint).SetZero()
	fmt.Println("s0 =", s.Ristretto())
}
