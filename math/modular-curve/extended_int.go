package main

import (
	"fmt"
	"math/big"
)

// coef[0] + coef[1]*w + coef[2]*w^2 + ... + coef[degree-1]*w^(degree) where w = exp(pi*i/degree)
// w^degree = -1
// note that zero value is invalid
type XInt struct {
	degree int
	coefs  []*big.Int
}

func (x *XInt) SetInt(val *big.Int, degree int) *XInt {
	x.degree = degree
	x.coefs = make([]*big.Int, degree)
	for i := range degree {
		x.coefs[i] = new(big.Int)
	}
	x.coefs[0].Set(val)
	return x
}

func (x *XInt) Set(x2 *XInt) *XInt {
	x.degree = x2.degree
	x.coefs = make([]*big.Int, x2.degree)
	for i := range x2.degree {
		x.coefs[i] = new(big.Int).Set(x2.coefs[i])
	}
	return x
}

func (x *XInt) SetCoef(coef []*big.Int) *XInt {
	degree := len(coef)
	x.degree = degree
	x.coefs = make([]*big.Int, degree)
	for i := range degree {
		x.coefs[i] = new(big.Int).Set(coef[i])
	}
	return x
}

func (x *XInt) RootOfUnity(n int) *XInt {
	if n <= 2 {
		panic("n > 2 should hold")
	}
	var coef []*big.Int
	if n%2 == 0 {
		coef = make([]*big.Int, n/2)
		for i := range n / 2 {
			coef[i] = new(big.Int)
		}
		coef[1] = new(big.Int).SetInt64(1)
	} else {
		coef = make([]*big.Int, n)
		for i := range n {
			coef[i] = new(big.Int)
		}
		coef[2] = new(big.Int).SetInt64(1)
	}
	return x.SetCoef(coef)
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
	x.degree = x1.degree
	coefs := make([]*big.Int, x.degree)
	for i := range x.degree {
		coefs[i] = new(big.Int).Add(x1.coefs[i], x2.coefs[i])
	}
	x.coefs = coefs
	return x
}

func (x *XInt) MulScalar(x1 *XInt, y *big.Int) *XInt {
	x.degree = x1.degree
	coefs := make([]*big.Int, x.degree)
	for i := range x.degree {
		coefs[i] = new(big.Int).Mul(x1.coefs[i], y)
	}
	x.coefs = coefs
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
	degree := x1.degree
	coefs := make([]*big.Int, degree)
	var tmp big.Int
	for i := range degree {
		coefs[i] = new(big.Int)
	}
	for i := range degree {
		for j := range degree {
			tmp.Mul(x1.coefs[j], x2.coefs[j])
			if i+j >= degree {
				tmp.Neg(&tmp)
			}
			coefs[(i+j)%degree].Add(coefs[(i+j)%degree], &tmp)
		}
	}
	x.degree = degree
	x.coefs = coefs
	return x
}

func (x *XInt) Resize(x1 *XInt, newDegree int) *XInt {
	if newDegree == x1.degree {
		return x.Set(x1)
	}
	if x1.degree%newDegree == 0 {
		// shrink mode
		r := x1.degree / newDegree
		coefs := make([]*big.Int, newDegree)
		for i := range x1.degree {
			if i%r == 0 {
				coefs[i/r] = new(big.Int).Set(x1.coefs[i])
			} else {
				if x1.coefs[i].Sign() != 0 {
					panic("coefficient should be zero")
				}
			}
		}
		x.degree = newDegree
		x.coefs = coefs
		return x
	}
	if newDegree%x1.degree == 0 {
		// expand mode
		r := newDegree / x1.degree
		coefs := make([]*big.Int, newDegree)
		for i := range newDegree {
			if i%r == 0 {
				coefs[i] = new(big.Int).Set(x1.coefs[i/r])
			} else {
				coefs[i] = new(big.Int)
			}
		}
		x.degree = newDegree
		x.coefs = coefs
		return x
	}
	panic(fmt.Sprint("degree mismatch: ", x1.degree, " -> ", newDegree))
}

func (x *XInt) Eq(y *XInt) bool {
	if x.degree != y.degree {
		panic("degree mismatch")
	}
	for i := range x.degree {
		if x.coefs[i].Cmp(y.coefs[i]) != 0 {
			return false
		}
	}
	return true
}
