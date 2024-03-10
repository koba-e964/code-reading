package main

import (
	"encoding/hex"
	"fmt"
	"slices"

	"filippo.io/edwards25519"
	rEdwards25519 "github.com/bwesterb/go-ristretto/edwards25519"
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

// * 4
func times4(e *edwards25519.Point) *edwards25519.Point {
	var eight [32]byte
	eight[0] = 4
	eightScalar, err := new(edwards25519.Scalar).SetCanonicalBytes(eight[:])
	if err != nil {
		panic(err)
	}
	return new(edwards25519.Point).ScalarMult(eightScalar, e)
}

// * 2^k l
func times2kL(e *edwards25519.Point, k int) *edwards25519.Point {
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
	value.Add(value, e)
	for i := 0; i < k; i++ {
		value.Add(value, value)
	}
	return value
}

func main() {
	for i := 6; i < 7; i++ {
		// Ristretto
		var buf [32]byte
		buf[0] = byte(i)
		s := new(rEdwards25519.ExtendedPoint)
		if ok := s.SetRistretto(&buf); !ok {
			continue
		}
		var zInv, x, y rEdwards25519.FieldElement
		zInv.Inverse(&s.Z)
		x.Mul(&s.X, &zInv)
		y.Mul(&s.Y, &zInv)
		yBytes := y.Bytes()
		yBytes[31] |= byte(x.IsNegativeI()) << 7
		point, err := new(edwards25519.Point).SetBytes(yBytes[:])
		if err != nil {
			panic(err)
		}
		fmt.Println("i =", i)
		fmt.Println("P =", s.String())
		fmt.Println("phi(P) =", toString(point))
		fourP := times4(point)
		twoLP := times2kL(point, 1)
		fourLP := times2kL(twoLP, 2)
		fmt.Println("4P =", toString(fourP))
		fmt.Println("2lP =", toString(twoLP))
		zero := edwards25519.NewIdentityPoint()
		if fourLP.Equal(zero) == 0 {
			panic("4lP != 0")
		}
		if fourP.Equal(zero) == 0 && twoLP.Equal(zero) == 0 {
			fmt.Println("phi(P) is a generator")
			break
		}
	}
}
