package main

import (
	"encoding/hex"
	"fmt"
	"slices"

	"filippo.io/edwards25519"
	"filippo.io/edwards25519/field"
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
		tmp.Add(tmp, torsion)
	}
	eightE := times8(e)
	if (eightE.Equal(edwards25519.NewIdentityPoint()) | fourLE.Equal(edwards25519.NewIdentityPoint())) == 0 {
		fmt.Println("  8 * P =", toString(eightE))
		fmt.Println("  4 * l * P =", toString(&fourLE))
		fmt.Println("  P is a generator")
	}
}

// This function ignores the sign of the point.
// It will not affect the order.
func fromXToEdwards(u *field.Element) *edwards25519.Point {
	y := new(field.Element).Subtract(u, new(field.Element).One())
	yDen := new(field.Element).Add(u, new(field.Element).One())
	yDen.Invert(yDen)
	y.Multiply(y, yDen)
	yBytes := y.Bytes()
	returned, err := edwards25519.NewIdentityPoint().SetBytes(yBytes)
	if err != nil {
		panic("fromXToEdwards: " + err.Error())
	}
	return returned
}

func elligator2(r *field.Element) *edwards25519.Point {
	A := new(field.Element).Mult32(new(field.Element).One(), uint32(486662))
	v := new(field.Element).Square(r)
	v.Mult32(v, 2)
	v.Add(v, new(field.Element).One()) // 1 + 2r^2
	v.Negate(v)
	v.Invert(v)
	v.Multiply(v, A)
	other := new(field.Element).Set(v)
	other.Add(other, A)
	other.Negate(other) // -v - A
	multiplier := new(field.Element).Add(v, A)
	multiplier.Multiply(multiplier, v)
	multiplier.Add(multiplier, new(field.Element).One()) // v^2 + Av + 1
	_, wasSquare := new(field.Element).SqrtRatio(v, other)
	if wasSquare == 1 {
		panic("should not be square")
	}
	x := other
	if _, wasSquare := new(field.Element).SqrtRatio(v, multiplier); wasSquare == 1 {
		x = v
	}
	return fromXToEdwards(x)
}

func main() {
	p := fromXToEdwards(new(field.Element).Mult32(new(field.Element).One(), 9))
	fmt.Println("P =", toString(p))
	for i := 0; i < 10; i++ {
		r := new(field.Element).Mult32(new(field.Element).One(), uint32(i))
		p := elligator2(r)
		fmt.Println("i =", i, "P =", toString(p))
		testPoint(p)
	}
}
