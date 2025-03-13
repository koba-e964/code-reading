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
	coef[0] = new(big.Int).Set(x)
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
	for i := range prec {
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

func (l *Laurent) MulScalar(x *Laurent, y *big.Int) *Laurent {
	coef := make([]*big.Int, x.prec)
	for i := range x.prec {
		coef[i] = new(big.Int).Mul(x.coef[i], y)
	}
	l.val = x.val
	l.prec = x.prec
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
	for i := range x.prec {
		l.coef[i] = new(big.Int).Set(x.coef[i])
	}
	return l
}

// Val returns the valuation of the Laurent series.
func (l *Laurent) Val() (val int, isZero bool) {
	allZero := true
	minNonZero := -1
	for i, v := range l.coef {
		if v.Cmp(big.NewInt(0)) != 0 {
			allZero = false
			minNonZero = i
			break
		}
	}
	if allZero {
		return 0, true
	}
	return l.val + minNonZero, false
}

// Add returns the sum of two Laurent series.
func (l *Laurent) Add(x, y *Laurent) *Laurent {
	xVal, xIsZero := x.Val()
	if xIsZero {
		return l.Set(y)
	}
	yVal, yIsZero := y.Val()
	if yIsZero {
		return l.Set(x)
	}
	if xVal > yVal {
		x, y = y, x
		xVal, yVal = yVal, xVal
	}
	prec := min(x.prec, y.prec+yVal-xVal)
	coef := make([]*big.Int, prec)
	for i := 0; i < prec; i++ {
		coef[i] = new(big.Int).Add(x.Coef(i+xVal), y.Coef(i+xVal))
	}
	l.val = xVal
	l.prec = prec
	l.coef = coef
	return l
}

// Sub returns the difference of two Laurent series.
func (l *Laurent) Sub(x, y *Laurent) *Laurent {
	tmp := new(Laurent).MulScalar(y, big.NewInt(-1))
	return l.Add(x, tmp)
}

func (l *Laurent) Prec() int {
	allZero := true
	minNonZero := -1
	for i, v := range l.coef {
		if v.Cmp(big.NewInt(0)) != 0 {
			allZero = false
			minNonZero = i
			break
		}
	}
	if allZero {
		return 0
	}
	return l.prec - minNonZero
}

func (l *Laurent) Shrink(x *Laurent) *Laurent {
	val, isZero := x.Val()
	if isZero {
		prec := x.prec + x.val
		l.val = 0
		l.prec = prec
		l.coef = make([]*big.Int, prec)
		for i := range prec {
			l.coef[i] = new(big.Int)
		}
		return l
	}
	prec := x.prec - (val - x.val)
	coef := make([]*big.Int, prec)
	for i := range prec {
		coef[i] = new(big.Int).Set(x.coef[i+val-x.val])
	}
	l.val = val
	l.prec = prec
	l.coef = coef
	return l
}

func (l *Laurent) Coef(i int) *big.Int {
	if i < l.val || i >= l.prec+l.val {
		return new(big.Int)
	}
	return new(big.Int).Set(l.coef[i-l.val])
}

func (l *Laurent) V(x *Laurent, v int) *Laurent {
	newCoef := make([]*big.Int, x.prec*v)
	newVal := x.val * v
	for i := range newCoef {
		newCoef[i] = new(big.Int)
	}
	for i := range x.prec {
		newCoef[i*v].Set(x.coef[i])
	}
	l.coef = newCoef
	l.val = newVal
	l.prec = len(newCoef)
	return l
}
