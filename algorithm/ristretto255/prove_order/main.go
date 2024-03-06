package main

import (
	"encoding/hex"
	"fmt"
	"slices"

	"filippo.io/edwards25519"
)

func toString(e *edwards25519.Point) string {
	x, y, z, _ := e.ExtendedCoordinates()
	zInv := z.Invert(z)
	x.Multiply(x, zInv)
	y.Multiply(y, zInv)
	xHex := hex.EncodeToString(x.Bytes())
	yHex := hex.EncodeToString(y.Bytes())
	return fmt.Sprintf("(%s, %s)", xHex, yHex)
}

// * 8
func times8(e *edwards25519.Point) *edwards25519.Point {
	var eight [32]byte
	eight[0] = 8
	eightScalar, err := new(edwards25519.Scalar).SetCanonicalBytes(eight[:])
	if err != nil {
		panic(err)
	}
	return new(edwards25519.Point).ScalarMult(eightScalar, e)
}

// * l
func timesL(e *edwards25519.Point) *edwards25519.Point {
	// https://safecurves.cr.yp.to/base.html
	order, err := hex.DecodeString("1000000000000000000000000000000014def9dea2f79cd65812631a5cf5d3ed")
	if err != nil {
		panic(err)
	}
	slices.Reverse(order)
	order[0] -= 1
	orderScalar, err := new(edwards25519.Scalar).SetCanonicalBytes(order)
	if err != nil {
		panic(err)
	}
	value := new(edwards25519.Point).ScalarMult(orderScalar, e)
	return value.Add(value, e)
}

func testPoint(e *edwards25519.Point) {
	fmt.Println("P =", toString(e))
	torsion := timesL(e)
	fmt.Println("  l * P =", toString(torsion))
	var eight [32]byte
	eight[0] = 8
	tmp := edwards25519.NewIdentityPoint()
	var fourLE edwards25519.Point
	for i := 0; i <= 8; i++ {
		if i == 4 {
			fourLE.Set(torsion)
		}
		fmt.Printf("  %d * l * P = %s\n", i, toString(tmp))
		tmp.Add(tmp, torsion)
	}
	eightE := times8(e)
	if (eightE.Equal(edwards25519.NewIdentityPoint()) | fourLE.Equal(edwards25519.NewIdentityPoint())) == 0 {
		fmt.Println("  8 * P =", toString(eightE))
		fmt.Println("  4 * l * P =", toString(&fourLE))
		fmt.Println("  P is a generator")
	}
}

func main() {
	base := edwards25519.NewGeneratorPoint()
	testPoint(base)
	for i := 2; i < 4; i++ {
		var buf [32]byte
		buf[0] = byte(i)
		tmp, err := new(edwards25519.Point).SetBytes(buf[:])
		if err != nil {
			continue
		}
		testPoint(tmp)
	}
}
