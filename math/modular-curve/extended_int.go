package main

import (
	"fmt"
	"math/big"
)

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func getPolyLength(n int) int {
	if n == 1 {
		return 1
	}
	return n - 1
}

func getMinPoly(n int) []*big.Int {
	if n == 1 {
		return []*big.Int{new(big.Int).SetInt64(0)}
	}
	minPoly := make([]*big.Int, n-1)
	for i := range n - 1 {
		if i%2 == 0 {
			minPoly[i] = new(big.Int).SetInt64(-1)
		} else {
			minPoly[i] = new(big.Int).SetInt64(1)
		}
	}
	return minPoly
}

// coef[0] + coef[1]*w + coef[2]*w^2 + ... + coef[degree-1]*w^(degree) where w = exp(pi*i/degree)
// w^degree = -1
// note that zero value is invalid
type XInt struct {
	degree  int
	coefs   []*big.Int
	minPoly []*big.Int
}

// degree must be 1 or an odd prime
// TODO: relax this condition
func (x *XInt) SetInt(val *big.Int, degree int) *XInt {
	if !isPrime(degree) && degree != 1 {
		panic(fmt.Sprint("degree must be 1 or an odd prime: ", degree))
	}
	x.degree = degree
	if degree == 1 {
		x.coefs = []*big.Int{new(big.Int).Set(val)}
		minPoly := []*big.Int{new(big.Int).SetInt64(0)}
		x.minPoly = minPoly
		return x
	}
	x.coefs = make([]*big.Int, degree-1)
	for i := range degree - 1 {
		x.coefs[i] = new(big.Int)
	}
	x.coefs[0].Set(val)
	minPoly := getMinPoly(degree)
	x.minPoly = minPoly
	return x
}

func (x *XInt) Set(x2 *XInt) *XInt {
	x.degree = x2.degree
	x.coefs = make([]*big.Int, len(x2.coefs))
	for i := range x2.coefs {
		x.coefs[i] = new(big.Int).Set(x2.coefs[i])
	}
	x.minPoly = make([]*big.Int, len(x2.minPoly))
	for i := range x2.minPoly {
		x.minPoly[i] = new(big.Int).Set(x2.minPoly[i])
	}
	return x
}

func (x *XInt) SetFromRaw(degree int, coef []*big.Int, minPoly []*big.Int) *XInt {
	if len(coef) != len(minPoly) {
		panic("len(coef) != len(minPoly)")
	}
	x.degree = degree
	x.coefs = make([]*big.Int, len(coef))
	for i := range coef {
		x.coefs[i] = new(big.Int).Set(coef[i])
	}
	x.minPoly = make([]*big.Int, len(minPoly))
	for i := range minPoly {
		x.minPoly[i] = new(big.Int).Set(minPoly[i])
	}
	return x
}

func (x *XInt) RootOfUnity(n int) *XInt {
	if n <= 1 {
		panic("n > 1 should hold")
	}
	var coef []*big.Int
	var degree int
	if n == 2 {
		coef = make([]*big.Int, 1)
		degree = 1
		coef[0] = new(big.Int).SetInt64(-1)
	} else if n%2 == 0 {
		coef = make([]*big.Int, n/2)
		for i := range n / 2 {
			degree = n / 2
			coef[i] = new(big.Int)
		}
		coef[1].SetInt64(1)
	} else {
		coef = make([]*big.Int, n-1)
		degree = n
		for i := range n - 1 {
			coef[i] = new(big.Int)
		}
		if n > 3 {
			coef[2].SetInt64(1)
		} else {
			coef[0].SetInt64(-1)
			coef[1].SetInt64(1)
		}
	}
	minPoly := getMinPoly(degree)
	return x.SetFromRaw(degree, coef, minPoly)
}

func (x *XInt) Coef(i int) *big.Int {
	return x.coefs[i]
}

func (x *XInt) Degree() int {
	return x.degree
}

func (x *XInt) Add(x1, x2 *XInt) *XInt {
	if x1.degree != x2.degree {
		panic(fmt.Sprint("add: degree mismatch: ", x1.degree, x2.degree))
	}
	if len(x1.coefs) != len(x2.coefs) {
		panic(fmt.Sprint("add: len(coefs) mismatch: ", len(x1.coefs), len(x2.coefs)))
	}
	x.degree = x1.degree
	coefs := make([]*big.Int, len(x1.coefs))
	for i := range x1.coefs {
		coefs[i] = new(big.Int).Add(x1.coefs[i], x2.coefs[i])
	}
	x.coefs = coefs
	x.minPoly = x1.minPoly
	return x
}

func (x *XInt) MulScalar(x1 *XInt, y *big.Int) *XInt {
	x.degree = x1.degree
	coefs := make([]*big.Int, len(x1.coefs))
	for i := range x1.coefs {
		coefs[i] = new(big.Int).Mul(x1.coefs[i], y)
	}
	x.coefs = coefs
	x.minPoly = x1.minPoly
	return x
}

func (x *XInt) Neg(x1 *XInt) *XInt {
	return x.MulScalar(x1, big.NewInt(-1))
}

func (x *XInt) Sub(x1, x2 *XInt) *XInt {
	if x1.degree != x2.degree {
		panic(fmt.Sprint("sub: degree mismatch: ", x1.degree, x2.degree))
	}
	var tmp XInt
	tmp.Neg(x2)
	return x.Add(x1, &tmp)
}

func (x *XInt) Mul(x1, x2 *XInt) *XInt {
	if x1.degree != x2.degree {
		panic(fmt.Sprint("mul: degree mismatch: ", x1.degree, x2.degree))
	}
	if len(x1.coefs) != len(x2.coefs) {
		panic(fmt.Sprint("mul: len(coefs) mismatch: ", len(x1.coefs), len(x2.coefs)))
	}
	degree := x1.degree
	if degree == 1 {
		x.degree = 1
		x.coefs = []*big.Int{new(big.Int).Mul(x1.coefs[0], x2.coefs[0])}
		x.minPoly = x1.minPoly
		return x
	}
	coefs := make([]*big.Int, len(x1.coefs)*2-1)
	var tmp big.Int
	for i := range coefs {
		coefs[i] = new(big.Int)
	}
	for i := range x1.coefs {
		for j := range x1.coefs {
			tmp.Mul(x1.coefs[i], x2.coefs[j])
			coefs[i+j].Add(coefs[i+j], &tmp)
		}
	}
	for i := len(x1.coefs)*2 - 2; i >= len(x1.coefs); i-- {
		cur := new(big.Int).Set(coefs[i])
		for j := range x1.minPoly {
			if x1.minPoly[j].Cmp(big.NewInt(-1)) == 0 {
				tmp.Neg(cur)
			} else if x1.minPoly[j].Cmp(big.NewInt(1)) == 0 {
				tmp.Set(cur)
			} else {
				tmp.Mul(cur, x1.minPoly[j])
			}
			coefs[i-len(x1.coefs)+j].Add(coefs[i-len(x1.coefs)+j], &tmp)
		}
		coefs[i].SetInt64(0)
	}
	x.degree = degree
	x.coefs = coefs[:len(x1.coefs)]
	x.minPoly = x1.minPoly
	return x
}

func (x *XInt) Pow(x1 *XInt, n int) *XInt {
	if n == 0 {
		return x.SetInt(big.NewInt(1), x1.degree)
	}
	var tmp XInt
	tmp.Set(x1)
	x.Set(x1)
	for i := 1; i < n; i++ {
		x.Mul(x, &tmp)
	}
	return x
}

// TODO: let it work for all cases
func (x *XInt) Resize(x1 *XInt, newDegree int) *XInt {
	if newDegree == x1.degree {
		return x.Set(x1)
	}
	if newDegree == 1 {
		// shrink mode
		coefs := make([]*big.Int, newDegree)
		for i := range x1.coefs {
			if i == 0 {
				coefs[0] = new(big.Int).Set(x1.coefs[i])
			} else {
				if x1.coefs[i].Sign() != 0 {
					panic(fmt.Sprint("coefficient should be zero: ", x1))
				}
			}
		}
		x.degree = newDegree
		x.coefs = coefs
		x.minPoly = getMinPoly(newDegree)
		return x
	}
	if x1.degree == 1 {
		// expand mode
		coefs := make([]*big.Int, getPolyLength(newDegree))
		for i := range coefs {
			if i == 0 {
				coefs[i] = new(big.Int).Set(x1.coefs[0])
			} else {
				coefs[i] = new(big.Int)
			}
		}
		x.degree = newDegree
		x.coefs = coefs
		x.minPoly = getMinPoly(newDegree)
		return x
	}
	panic(fmt.Sprint("degree mismatch: ", x1.degree, " -> ", newDegree))
}

func (x *XInt) Eq(y *XInt) bool {
	if x.degree != y.degree {
		panic("degree mismatch")
	}
	for i := range x.coefs {
		if x.coefs[i].Cmp(y.coefs[i]) != 0 {
			return false
		}
	}
	return true
}

func (x *XInt) String() string {
	return fmt.Sprint(x.coefs)
}
