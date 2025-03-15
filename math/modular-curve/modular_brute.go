package main

import (
	"fmt"
	"math/big"
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

func divide(deg, n int) (int, int) {
	x, y := deg/n, deg%n
	if x < y {
		panic(fmt.Sprintf("x >= y should hold:  %d/%d", deg, n))
	}
	return x, y
}

func ModularBrute(n int) []Entry {
	if n <= 1 {
		panic("n > 1 should hold")
	}
	one := new(XInt).SetInt(big.NewInt(1), 1)

	entries := make([]Entry, 0)
	psi := Psi(n)
	globalPrec := psi*n + 1
	j := new(Laurent).JInv(globalPrec)
	jn := new(Laurent).V(j, n)
	jjn := new(Laurent).Mul(j, jn)
	jPow := new(Laurent).Pow(j, psi)
	jnPow := new(Laurent).Pow(jn, psi)
	dat := new(Laurent).Add(jnPow, jPow)
	entries = append(entries, Entry{XDeg: psi, YDeg: 0, Coef: big.NewInt(1)})
	jPow.Pow(j, psi-1-n)
	jnPow.Pow(jn, psi-1-n)
	jjnPow := new(Laurent).Pow(jjn, n)
	tmp := new(Laurent).Add(jPow, jnPow)
	if psi == 1+n {
		tmp.SetInt(one, globalPrec)
	}
	tmp.Mul(tmp, jjnPow)
	dat.Sub(dat, tmp)
	entries = append(entries, Entry{XDeg: psi - 1, YDeg: n, Coef: big.NewInt(-1)})

	for {
		for _, r := range entries {
			fmt.Printf("%d %d %v\n", r.XDeg, r.YDeg, r.Coef)
		}
		val, isZero := dat.Val()
		if val > 0 || isZero {
			break
		}
		a, b := divide(-val, n)
		jjnPow.Pow(jjn, b)
		jPow.Pow(j, a-b)
		jnPow.Pow(jn, a-b)
		tmp.Add(jPow, jnPow)
		if a == b {
			tmp.SetInt(one, globalPrec)
		}
		tmp.Mul(tmp, jjnPow)
		coef := dat.Coef(val)
		coef.Neg(coef)
		entries = append(entries, Entry{XDeg: a, YDeg: b, Coef: coef.Coef(0)})
		tmp.MulScalar(tmp, coef)
		dat.Add(dat, tmp)
	}
	return entries
}
