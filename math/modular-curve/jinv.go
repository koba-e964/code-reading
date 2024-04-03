package main

import (
	"math/big"
)

type Laurent struct {
	val  int
	prec int
	coef []*big.Int
}

// https://mathoverflow.net/questions/71704/computing-the-q-series-of-the-j-invariant
func (l *Laurent) JInv(prec int) *Laurent {
	if prec < 0 {
		panic("negative precision")
	}
	var g2 Laurent
	g2.G2(prec)
	g2.Pow(&g2, 3)
	var deltaInv Laurent
	deltaInv.Delta(prec)
	deltaInv.Inv(&deltaInv)
	l.Mul(&g2, &deltaInv)
	return l
}

func (l *Laurent) G2(prec int) *Laurent {
	if prec < 0 {
		panic("negative precision")
	}
	coef := make([]*big.Int, prec)
	coef[0] = big.NewInt(1)
	for i := 1; i < prec; i++ {
		coef[i] = new(big.Int)
		var tmp big.Int
		for j := 1; j <= i; j++ {
			if i%j == 0 {
				tmp.SetInt64(int64(j))
				tmp.Exp(&tmp, big.NewInt(3), nil)
				coef[i].Add(coef[i], &tmp)
			}
		}
		coef[i].Mul(coef[i], big.NewInt(240))
	}
	l.val = 0
	l.prec = prec
	l.coef = coef
	return l
}

func (l *Laurent) Delta(prec int) *Laurent {
	if prec < 0 {
		panic("negative precision")
	}
	if prec > 0 {
		var tmp, prod Laurent
		prod.SetInt(big.NewInt(1), prec)
		for i := 1; i < prec; i++ {
			tmp.SetInt(big.NewInt(1), prec)
			tmp.SetCoef(i, big.NewInt(-1))
			tmp.Pow(&tmp, 24)
			prod.Mul(&prod, &tmp)
		}
		l.Shl(&prod, 1)
	} else {
		l.val = 0
		l.prec = 0
		l.coef = make([]*big.Int, 0)
	}
	return l
}

func (l *Laurent) SetInt(x *big.Int, prec int) *Laurent {
	if prec < 0 {
		panic("negative precision")
	}
	coef := make([]*big.Int, prec)
	coef[0] = x
	for i := 1; i < prec; i++ {
		coef[i] = new(big.Int)
	}
	l.val = 0
	l.prec = prec
	l.coef = coef
	return l
}

func (l *Laurent) SetCoef(i int, x *big.Int) *Laurent {
	if i < 0 || i >= l.prec {
		panic("index out of range")
	}
	l.coef[i] = x
	return l
}

func (l *Laurent) Shl(x *Laurent, n int) *Laurent {
	coef := make([]*big.Int, x.prec)
	for i := 0; i < x.prec; i++ {
		coef[i] = x.coef[i]
	}
	l.val = x.val + n
	l.prec = x.prec
	l.coef = coef
	return l
}

func (l *Laurent) Mul(x, y *Laurent) *Laurent {
	prec := min(x.prec, y.prec)
	coef := make([]*big.Int, prec)
	for i := 0; i < prec; i++ {
		coef[i] = new(big.Int)
		for j := max(0, i-y.prec+1); j <= min(i, x.prec-1); j++ {
			var tmp big.Int
			tmp.Mul(x.coef[j], y.coef[i-j])
			coef[i].Add(coef[i], &tmp)
		}
	}
	l.val = x.val + y.val
	l.prec = prec
	l.coef = coef
	return l
}

func (l *Laurent) Pow(x *Laurent, n int) *Laurent {
	if n == 0 {
		return l.SetInt(big.NewInt(1), x.prec)
	}
	if n == 1 {
		return l.Shl(x, 0)
	}
	var tmp Laurent
	tmp.Pow(x, n/2)
	tmp.Mul(&tmp, &tmp)
	if n%2 == 1 {
		tmp.Mul(&tmp, x)
	}
	l.Set(&tmp)
	return l
}

func (l *Laurent) Inv(x *Laurent) *Laurent {
	if x.coef[0].Cmp(big.NewInt(1)) != 0 {
		panic("constant term must be 1")
	}
	tmp := make([]*big.Int, x.prec)
	for i := 1; i < x.prec; i++ {
		tmp[i] = new(big.Int).Neg(x.coef[i])
	}
	coef := make([]*big.Int, x.prec)
	coef[0] = big.NewInt(1)
	for i := 1; i < x.prec; i++ {
		coef[i] = new(big.Int).Set(tmp[i])
		for j := i + 1; j < x.prec; j++ {
			tmp[j].Sub(tmp[j], new(big.Int).Mul(coef[i], x.coef[j-i]))
		}
	}

	l.val = -x.val
	l.prec = x.prec
	l.coef = coef

	return l
}

func (l *Laurent) Set(x *Laurent) *Laurent {
	l.val = x.val
	l.prec = x.prec
	l.coef = make([]*big.Int, x.prec)
	copy(l.coef, x.coef)
	return l
}

// Val returns the valuation of the Laurent series.
func (l *Laurent) Val() int {
	return l.val
}

func (l *Laurent) Prec() int {
	return l.prec
}

func (l *Laurent) Coef(i int) *big.Int {
	if i < l.val || i >= l.prec+l.val {
		return new(big.Int)
	}
	return new(big.Int).Set(l.coef[i-l.val])
}
