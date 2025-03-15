package main

import (
	"fmt"
	"math/big"
)

const NAIVE_MUL_THRESHOLD = 10

func naiveMul(a, b []*XInt, degree, prec int) []*XInt {
	n := len(a)
	m := len(b)
	prec = min(prec, n+m-1)
	z := make([]*XInt, prec)
	for i := range z {
		z[i] = new(XInt).SetInt(big.NewInt(0), degree)
	}
	for i := range a {
		for j := range min(len(b), prec-i) {
			z[i+j].Add(z[i+j], new(XInt).Mul(a[i], b[j]))
		}
	}
	return z
}

func karatsuba(a, b []*XInt, degree, prec int) []*XInt {
	if prec < 0 {
		return []*XInt{}
	}
	if len(a) <= NAIVE_MUL_THRESHOLD || len(b) <= NAIVE_MUL_THRESHOLD {
		return naiveMul(a, b, degree, prec)
	}
	n := min(prec, max(len(a), len(b)))
	n = (n + 1) / 2 * 2
	m := n / 2
	prec = min(prec, 2*n-1)
	aExtend := make([]*XInt, n)
	bExtend := make([]*XInt, n)
	copy(aExtend, a)
	for i := len(a); i < n; i++ {
		aExtend[i] = new(XInt).SetInt(big.NewInt(0), degree)
	}
	copy(bExtend, b)
	for i := len(b); i < n; i++ {
		bExtend[i] = new(XInt).SetInt(big.NewInt(0), degree)
	}

	a0, a1 := aExtend[:m], aExtend[m:]
	b0, b1 := bExtend[:m], bExtend[m:]
	z0 := karatsuba(a0, b0, degree, prec)
	z2 := karatsuba(a1, b1, degree, prec-m)
	t := make([]*XInt, m)
	for i := range t {
		t[i] = new(XInt).Add(a0[i], a1[i])
	}
	u := make([]*XInt, m)
	for i := range u {
		u[i] = new(XInt).Add(b0[i], b1[i])
	}
	z1 := karatsuba(t, u, degree, prec-m)
	for i := range min(prec-m, len(z1)) {
		z1[i].Sub(z1[i], z0[i])
		z1[i].Sub(z1[i], z2[i])
	}
	z := make([]*XInt, prec)
	for i := range z {
		z[i] = new(XInt).SetInt(big.NewInt(0), degree)
	}
	for i := range z0 {
		z[i].Add(z[i], z0[i])
	}
	for i := range z1 {
		if i+m >= len(z) {
			break
		}
		z[i+m].Add(z[i+m], z1[i])
	}
	for i := range z2 {
		if i+2*m >= len(z) {
			break
		}
		z[i+2*m].Add(z[i+2*m], z2[i])
	}
	return z
}

type Laurent struct {
	val    int
	prec   int
	coef   []*XInt
	degree int
}

// https://mathoverflow.net/questions/71704/computing-the-q-series-of-the-j-invariant
// degree = 1
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
	l.assertValid()
	return l
}

func (l *Laurent) G2(prec int) *Laurent {
	if prec < 0 {
		panic("negative precision")
	}
	coef := make([]*XInt, prec)
	degree := 1
	coef[0] = new(XInt).SetInt(big.NewInt(1), degree)
	for i := 1; i < prec; i++ {
		coef[i] = new(XInt).SetInt(big.NewInt(0), degree)
		var tmp big.Int
		for j := 1; j <= i; j++ {
			if i%j == 0 {
				tmp.SetInt64(int64(j))
				tmp.Exp(&tmp, big.NewInt(3), nil)
				coef[i].Add(coef[i], new(XInt).SetInt(&tmp, degree))
			}
		}
		coef[i].MulScalar(coef[i], big.NewInt(240))
	}
	l.val = 0
	l.prec = prec
	l.coef = coef
	l.degree = 1
	l.assertValid()
	return l
}

// degree = 1
func (l *Laurent) Delta(prec int) *Laurent {
	if prec < 0 {
		panic("negative precision")
	}
	if prec > 0 {
		var tmp, prod Laurent
		prod.SetInt(new(XInt).SetInt(big.NewInt(1), 1), prec)
		for i := 1; i < prec; i++ {
			tmp.SetInt(new(XInt).SetInt(big.NewInt(1), 1), prec)
			tmp.SetCoef(i, new(XInt).SetInt(big.NewInt(-1), 1))
			tmp.Pow(&tmp, 24)
			prod.Mul(&prod, &tmp)
		}
		l.Shl(&prod, 1)
	} else {
		l.val = 0
		l.prec = 0
		l.coef = make([]*XInt, 0)
	}
	l.degree = 1
	l.assertValid()
	return l
}

func (l *Laurent) SetInt(x *XInt, prec int) *Laurent {
	degree := x.degree
	if prec < 0 {
		panic("negative precision")
	}
	coef := make([]*XInt, prec)
	coef[0] = new(XInt).Set(x)
	for i := 1; i < prec; i++ {
		coef[i] = new(XInt).SetInt(big.NewInt(0), degree)
	}
	l.val = 0
	l.prec = prec
	l.coef = coef
	l.degree = degree
	l.assertValid()
	return l
}

func (l *Laurent) SetCoef(i int, x *XInt) *Laurent {
	if i < l.val || i >= l.prec+l.val {
		panic("index out of range")
	}
	l.coef[i-l.val] = x
	return l
}

func (l *Laurent) Shl(x *Laurent, n int) *Laurent {
	coef := make([]*XInt, x.prec)
	for i := range x.prec {
		coef[i] = x.coef[i]
	}
	l.val = x.val + n
	l.prec = x.prec
	l.coef = coef
	l.degree = x.degree
	l.assertValid()
	return l
}

func (l *Laurent) Mul(x, y *Laurent) *Laurent {
	prec := min(x.prec, y.prec)
	if x.degree != y.degree {
		panic(fmt.Sprint("Laurent.Mul: different degrees: ", x.degree, y.degree))
	}
	degree := x.degree
	coefs := karatsuba(x.coef[:prec], y.coef[:prec], degree, prec)
	coefs = coefs[:prec]
	l.val = x.val + y.val
	l.prec = prec
	l.coef = coefs
	l.degree = degree
	l.assertValid()
	return l
}

func (l *Laurent) MulScalar(x *Laurent, y *XInt) *Laurent {
	coef := make([]*XInt, x.prec)
	degree := x.degree
	for i := range x.prec {
		coef[i] = new(XInt).Mul(x.coef[i], y)
	}
	l.val = x.val
	l.prec = x.prec
	l.coef = coef
	l.degree = degree
	l.assertValid()
	return l
}

func (l *Laurent) Pow(x *Laurent, n int) *Laurent {
	degree := x.degree
	if n == 0 {
		return l.SetInt(new(XInt).SetInt(big.NewInt(1), degree), x.prec)
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
	l.assertValid()
	return l
}

func (l *Laurent) Inv(x *Laurent) *Laurent {
	degree := x.degree
	if !x.coef[0].Eq(new(XInt).SetInt(big.NewInt(1), degree)) {
		panic("constant term must be 1")
	}
	tmp := make([]*XInt, x.prec)
	for i := 1; i < x.prec; i++ {
		tmp[i] = new(XInt).Neg(x.coef[i])
	}
	coef := make([]*XInt, x.prec)
	coef[0] = new(XInt).SetInt(big.NewInt(1), degree)
	for i := 1; i < x.prec; i++ {
		coef[i] = new(XInt).Set(tmp[i])
		for j := i + 1; j < x.prec; j++ {
			tmp[j].Sub(tmp[j], new(XInt).Mul(coef[i], x.coef[j-i]))
		}
	}

	l.val = -x.val
	l.prec = x.prec
	l.coef = coef
	l.degree = degree

	l.assertValid()
	return l
}

func (l *Laurent) Set(x *Laurent) *Laurent {
	l.val = x.val
	l.prec = x.prec
	coef := make([]*XInt, x.prec)
	for i := range x.prec {
		coef[i] = new(XInt).Set(x.coef[i])
	}
	l.coef = coef
	l.degree = x.degree
	l.assertValid()
	return l
}

// Val returns the valuation of the Laurent series.
func (l *Laurent) Val() (val int, isZero bool) {
	allZero := true
	minNonZero := -1
	degree := l.degree
	zero := new(XInt).SetInt(big.NewInt(0), degree)
	for i, v := range l.coef {
		if !v.Eq(zero) {
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
	x.assertValid()
	y.assertValid()
	if x.degree != y.degree {
		panic(fmt.Sprint("Laurent.Add: different degrees", x.degree, y.degree))
	}
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
	coef := make([]*XInt, prec)
	for i := range prec {
		coef[i] = new(XInt).Add(x.Coef(i+xVal), y.Coef(i+xVal))
	}
	l.val = xVal
	l.prec = prec
	l.coef = coef
	l.degree = x.degree
	l.assertValid()
	return l
}

// Sub returns the difference of two Laurent series.
func (l *Laurent) Sub(x, y *Laurent) *Laurent {
	degree := y.degree
	tmp := new(Laurent).MulScalar(y, new(XInt).SetInt(big.NewInt(-1), degree))
	return l.Add(x, tmp)
}

// MulShift returns `x(yq)“.
func (l *Laurent) MulShift(x *Laurent, y, yInv *XInt) *Laurent {
	if y.Degree() != yInv.Degree() {
		panic("degree mismatch")
	}
	mul := new(XInt).Mul(y, yInv)
	if !mul.Eq(new(XInt).SetInt(big.NewInt(1), y.Degree())) {
		panic(fmt.Sprint("y * yInv != 1: ", y, yInv, mul))
	}
	val := x.val
	degree := x.degree
	prec := x.prec
	coef := make([]*XInt, prec)
	cur := new(XInt).SetInt(big.NewInt(1), degree)
	for i := -1; i >= val; i-- {
		cur.Mul(cur, yInv)
		if i-val < prec {
			coef[i-val] = new(XInt).Mul(x.coef[i-val], cur)
		}
	}
	cur.SetInt(big.NewInt(1), degree)
	for i := range val + prec {
		if i >= val {
			coef[i-val] = new(XInt).Mul(x.coef[i-val], cur)
		}
		cur.Mul(cur, y)
	}
	l.val = x.val
	l.prec = prec
	l.coef = coef
	l.degree = degree
	l.assertValid()
	return l
}

func (l *Laurent) Prec() int {
	allZero := true
	minNonZero := -1
	degree := l.degree
	zero := new(XInt).SetInt(big.NewInt(0), degree)
	for i, v := range l.coef {
		if !v.Eq(zero) {
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

func (l *Laurent) ShrinkToNormalForm(x *Laurent) *Laurent {
	val, isZero := x.Val()
	degree := x.degree
	if isZero {
		prec := x.prec + x.val
		l.val = 0
		l.prec = prec
		l.coef = make([]*XInt, prec)
		for i := range prec {
			l.coef[i] = new(XInt).SetInt(big.NewInt(0), degree)
		}
		return l
	}
	prec := x.prec - (val - x.val)
	coef := make([]*XInt, prec)
	for i := range prec {
		coef[i] = new(XInt).Set(x.coef[i+val-x.val])
	}
	l.val = val
	l.prec = prec
	l.coef = coef
	l.degree = degree
	l.assertValid()
	return l
}

func (l *Laurent) TruncatePrec(x *Laurent, newPrec int) int {
	if newPrec > x.prec {
		panic(fmt.Sprint("newPrec > x.prec: ", newPrec, x.prec))
	}
	degree := x.degree
	coef := make([]*XInt, newPrec)
	for i := range newPrec {
		coef[i] = new(XInt).Set(x.coef[i])
	}
	l.val = x.val
	l.prec = newPrec
	l.coef = coef
	l.degree = degree
	l.assertValid()
	return newPrec
}

func (l *Laurent) Coef(i int) *XInt {
	degree := l.degree
	if i < l.val || i >= l.prec+l.val {
		return new(XInt).SetInt(big.NewInt(0), degree)
	}
	return new(XInt).Set(l.coef[i-l.val])
}

func (l *Laurent) Degree() int {
	return l.degree
}

// V returns x(q^v).
func (l *Laurent) V(x *Laurent, v int) *Laurent {
	degree := x.degree
	newCoef := make([]*XInt, x.prec*v)
	newVal := x.val * v
	for i := range newCoef {
		newCoef[i] = new(XInt).SetInt(big.NewInt(0), degree)
	}
	for i := range x.prec {
		newCoef[i*v].Set(x.coef[i])
	}
	l.coef = newCoef
	l.val = newVal
	l.prec = len(newCoef)
	l.degree = degree
	l.assertValid()
	return l
}

func quo(a, b int) int {
	r := (a%b + b) % b
	return (a - r) / b
}

func (l *Laurent) InvV(x *Laurent, v int) *Laurent {
	degree := x.degree
	zero := new(XInt).SetInt(big.NewInt(0), degree)
	newVal := quo(x.val+v-1, v)
	newPrec := quo(x.val+x.prec+v-1, v) - newVal
	newCoef := make([]*XInt, newPrec)
	for i := range newCoef {
		newCoef[i] = new(XInt).SetInt(big.NewInt(0), degree)
	}
	for i := x.val; i < x.val+x.prec; i++ {
		if i%v != 0 {
			if !x.Coef(i).Eq(zero) {
				panic(fmt.Sprint("non-zero coefficient at ", i))
			}
		} else {
			newCoef[(i-x.val)/v] = new(XInt).Set(x.Coef(i))
		}
	}
	l.coef = newCoef
	l.val = newVal
	l.prec = len(newCoef)
	l.degree = degree
	l.assertValid()
	return l
}

func (l *Laurent) Resize(x *Laurent, newDegree int) *Laurent {
	coef := make([]*XInt, x.prec)
	for i := range x.prec {
		coef[i] = new(XInt).Resize(x.coef[i], newDegree)
	}
	l.val = x.val
	l.prec = x.prec
	l.coef = coef
	l.degree = newDegree
	l.assertValid()
	return l
}

func (l *Laurent) assertValid() {
	degree := l.degree
	for i := range l.prec {
		if l.coef[i].Degree() != degree {
			panic(fmt.Sprintf("invalid degree: %d", l.coef[i].Degree()))
		}
	}
}

func (l *Laurent) String() string {
	var s string
	for i := l.val; i < l.val+l.prec; i++ {
		s += fmt.Sprintf("%d => %v\n", i, l.Coef(i).String())
	}
	return s
}
