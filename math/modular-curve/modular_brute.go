package main

import (
	"fmt"
	"math/big"
	"slices"
)

type Entry struct {
	XDeg int
	YDeg int
	Coef *big.Int
}

func Psi(n int) int {
	p := 2
	psi := 1
	for n > 1 {
		c := 0
		for n%p == 0 {
			n /= p
			c++
		}
		if c >= 1 {
			psi *= p + 1
			for range c - 1 {
				psi *= p
			}
		}
		p++
	}
	return psi
}

func polyMulXPlusB(a []*Laurent, b *Laurent) []*Laurent {
	if a[0].Degree() != b.Degree() {
		panic(fmt.Sprint("polyMul: different degrees", a[0].Degree(), b.Degree()))
	}
	nTerms := len(a) + 1
	prod := make([]*Laurent, nTerms)
	for i := range a {
		prod[i+1] = new(Laurent).Set(a[i])
	}
	prod[0] = new(Laurent).SetInt(new(XInt).SetInt(big.NewInt(0), b.Degree()), b.Prec())
	for i := range a {
		prod[i].Add(prod[i], new(Laurent).Mul(a[i], b))
	}
	return prod
}

func recoverAsJsPolynomial(l, j *Laurent, lowerBound int) []*big.Int {
	val, isZero := l.Val()
	if isZero {
		return []*big.Int{big.NewInt(0)}
	}
	prec := l.Prec()
	if val+prec <= 0 {
		panic(fmt.Sprintf("val + prec > 0 should hold: %d + %d", val, prec))
	}
	tmp := new(Laurent).Set(l)
	coefs := make([]*big.Int, 0, -val)
	for {
		val, isZero := tmp.Val()
		if val > -lowerBound || isZero {
			break
		}
		coef := tmp.Coef(val)
		for len(coefs) <= -val {
			coefs = append(coefs, big.NewInt(0))
		}
		coefs[-val].Add(coefs[-val], coef.Coef(0))
		negCoef := new(XInt).Neg(coef)
		jPow := new(Laurent).Pow(j, -val)
		tmp2 := new(Laurent).MulScalar(jPow, negCoef)
		tmp.Add(tmp, tmp2)
	}
	return coefs
}

// works only for prime n
// TODO: support non-prime n
func ModularBrute(n int) []Entry {
	if n <= 1 {
		panic("n > 1 should hold")
	}
	w := new(XInt).RootOfUnity(n)
	wInv := new(XInt).Pow(w, n-1)
	degree := w.Degree()
	one := new(XInt).SetInt(big.NewInt(1), degree)
	minusOne := new(XInt).Neg(one)

	entries := make([]Entry, 0)
	psi := Psi(n)
	globalPrec := psi*n + 1
	jRaw := new(Laurent).JInv(globalPrec)

	// j(q)
	j := new(Laurent).V(jRaw, n)
	j.Resize(j, degree)
	// j(q^{1/n}w^i)
	jfrac := make([]*Laurent, n)
	wPow := new(XInt).Set(one)
	wInvPow := new(XInt).Set(one)
	for i := range jfrac {
		jfrac[i] = new(Laurent).Resize(jRaw, degree)
		jfrac[i] = new(Laurent).MulShift(jfrac[i], wPow, wInvPow)
		wPow.Mul(wPow, w)
		wInvPow.Mul(wInvPow, wInv)
	}
	mul := []*Laurent{new(Laurent).SetInt(one, globalPrec)}
	for i := range jfrac {
		mul = polyMulXPlusB(mul, jfrac[i].MulScalar(jfrac[i], minusOne))
	}
	for _, r := range mul {
		r.Resize(r, 1)
		r.InvV(r, n)
	}
	jRaw.TruncatePrec(jRaw, psi+1)
	// j(q^n)
	jn := new(Laurent).V(jRaw, n)
	mul = polyMulXPlusB(mul, jn.MulScalar(jn, new(XInt).SetInt(big.NewInt(-1), 1)))
	for i, r := range mul {
		coefs := recoverAsJsPolynomial(r, jRaw, i)
		for j := range coefs {
			if i <= j && coefs[j].Sign() != 0 {
				entries = append(entries, Entry{
					XDeg: j,
					YDeg: i,
					Coef: coefs[j],
				})
			}
		}
	}
	slices.SortFunc(entries, func(a, b Entry) int {
		if a.XDeg < b.XDeg {
			return -1
		}
		if a.XDeg > b.XDeg {
			return 1
		}
		if a.YDeg < b.YDeg {
			return -1
		}
		if a.YDeg > b.YDeg {
			return 1
		}
		return 0
	})
	return entries
}
