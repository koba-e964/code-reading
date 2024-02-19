package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
)

// https://www.jstage.jst.go.jp/article/sicejl1962/28/9/28_9_803/_pdf

type GF byte

const Poly GF = 0x1d

func (a GF) mul(b GF) GF {
	var res GF
	for b != 0 {
		if (b & 1) != 0 {
			res ^= a
		}
		if (a & 0x80) != 0 {
			a = a<<1 ^ Poly
		} else {
			a <<= 1
		}
		b >>= 1
	}
	return res
}

func (a GF) pow(x int) GF {
	res := GF(1)
	cur := a
	for x > 0 {
		if x&1 != 0 {
			res = res.mul(cur)
		}
		cur = cur.mul(cur)
		x >>= 1
	}
	return res
}

func (a GF) inv() GF {
	if a == 0 {
		panic("division by zero (GF(2^8))")
	}
	return a.pow(254)
}

func polyXor(a, b []GF) []GF {
	if len(a) < len(b) {
		a, b = b, a
	}
	res := make([]GF, len(a))
	copy(res, a)
	for i, v := range b {
		res[i] ^= v
	}
	return res
}

func polyMul(a, b []GF) []GF {
	if len(a) == 0 || len(b) == 0 {
		return nil
	}
	res := make([]GF, len(a)+len(b)-1)
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			res[i+j] ^= a[i].mul(b[j])
		}
	}
	return res
}

func polyDiv(a, b []GF) ([]GF, []GF) {
	if len(b) == 0 {
		panic("division by zero (poly)")
	}
	tmp := make([]GF, len(a))
	copy(tmp, a)
	q := make([]GF, len(a)-len(b)+1)
	r := make([]GF, len(b)-1)
	invl := b[len(b)-1].inv()
	for i := len(a) - len(b); i >= 0; i-- {
		q[i] = tmp[i+len(b)-1].mul(invl)
		for j := 0; j < len(b); j++ {
			tmp[i+j] ^= q[i].mul(b[j])
		}
	}
	copy(r, tmp[:len(b)-1])
	return q, r
}

func polyOf(a []GF, x GF) GF {
	res := GF(0)
	cur := GF(1)
	for _, v := range a {
		res ^= cur.mul(v)
		cur = cur.mul(x)
	}
	return res
}

func polyReduce(a []GF) []GF {
	for i := len(a) - 1; i >= 0; i-- {
		if a[i] != 0 {
			return a[:i+1]
		}
	}
	return nil
}

// g = a * s + b * t
func polyGCD(a, b []GF) (g []GF, s []GF, t []GF) {
	a = polyReduce(a)
	b = polyReduce(b)
	if len(b) > 0 && len(a) > 2 {
		q, r := polyDiv(a, b)
		r = polyReduce(r)
		g, s, t := polyGCD(b, r)
		// g = s * b + t * r = t * a + (s + q * t) * b
		return g, t, polyReduce(polyXor(s, polyMul(q, t)))
	}
	return a, []GF{1}, []GF{}
}

// Only can find exactly 2 errors
// Ref: https://www.jstage.jst.go.jp/article/sicejl1962/28/9/28_9_803/_pdf
func rs2824Direct2(I, syndrome, table []GF) bool {
	logger := log.New(os.Stdout, "rs2824Direct2: ", 0)
	quad := make([]GF, 3)
	quad[0] = syndrome[1].mul(syndrome[3]) ^ syndrome[2].mul(syndrome[2])
	quad[1] = syndrome[0].mul(syndrome[3]) ^ syndrome[1].mul(syndrome[2])
	quad[2] = syndrome[0].mul(syndrome[2]) ^ syndrome[1].mul(syndrome[1])
	logger.Println("quad =", quad)
	if quad[2] != 0 {
		quad0inv := quad[2].inv()
		tmp := make([]GF, 3)
		for i := 0; i < 3; i++ {
			tmp[i] = quad[i].mul(quad0inv)
		}
		logger.Println("normalize(quad) =", tmp)
	} else {
		return false
	}
	loc := []int{}
	for i := 0; i < 28; i++ {
		if polyOf(quad, table[i]) == 0 {
			loc = append(loc, i)
		}
	}
	if len(loc) == 2 {
		i, j := loc[0], loc[1]
		ej := (table[i].mul(syndrome[0]) ^ syndrome[1]).mul((table[i] ^ table[j]).inv())
		ei := syndrome[0] ^ ej
		logger.Printf("error values = 0x%02x (@%d), 0x%02x (@%d)\n", ei, i, ej, j)
	} else {
		logger.Println("Invalid number of error locations:", loc)
	}
	return len(loc) == 2
}

// Only can find exactly 1 error
func rs2824Direct1(I, syndrome, table []GF) bool {
	if syndrome[0] == 0x00 {
		return false
	}
	logger := log.New(os.Stdout, "rs2824Direct1: ", 0)
	alphaPow := syndrome[1].mul(syndrome[0].inv())
	for i := 0; i < 28; i++ {
		if alphaPow == table[i] {
			logger.Printf("error value = 0x%02x (@%d)\n", syndrome[0], i)
			return true
		}
	}
	return false
}

func rs2824Euclidean(I, syndrome, table []GF) {
	logger := log.New(os.Stdout, "rs2824Euclidean: ", 0)

	if reflect.DeepEqual(syndrome, make([]GF, 4)) {
		logger.Println("no error")
		return
	}
	x4 := make([]GF, 5)
	x4[4] = 1
	gcd, _, s := polyGCD(x4, syndrome)
	invs0 := s[0].inv()
	for i := range s {
		s[i] = s[i].mul(invs0)
	}
	for i := range gcd {
		gcd[i] = gcd[i].mul(invs0)
	}
	logger.Println("Omega =", gcd)
	logger.Println("s (error locator) =", s)

	loc := []int{}

	for i := 0; i < 28; i++ {
		if polyOf(s, table[i].inv()) == 0 {
			loc = append(loc, i)
		}
	}
	logger.Println("error locations =", loc)

	// error values
	s_der := make([]GF, len(s)-1)
	for i := 0; i < len(s)-1; i++ {
		if (i+1)%2 != 0 {
			s_der[i] = s[i+1]
		}
	}
	logger.Println("s' =", s_der)
	for _, l := range loc {
		x := table[l].inv()
		// In Forney algorithm, multiplication by table[l] is necessary
		// because we start from c = 0
		// What we want as the denominator: \prod_j (1-a_i^{-1}a_j)
		// s'(a_i^{-1}): (-1)^{n-1}a_0...a_{n-1} \prod_j (a_i^{-1} - a_j^{-1}) = a_i\prod_j (1-a_i^{-1}a_j)
		// https://en.wikipedia.org/wiki/Forney_algorithm
		e := polyOf(gcd, x).mul(polyOf(s_der, x).inv()).mul(table[l])
		logger.Printf("error value = 0x%02x (@%d)\n", e, l)
	}
}

func main() {
	table := make([]GF, 255)
	table[0] = 1
	for i := 1; i < 255; i++ {
		table[i] = table[i-1].mul(2)
	}
	for i, v := range table {
		fmt.Printf("%d: 0x%02x\n", i, v)
	}
	// CD generating polynomial
	g := []GF{0x01}
	for i := 0; i < 4; i++ {
		g = polyMul(g, []GF{table[i], 0x01})
	}
	// a^6 + a^78x + a^249x^2 + a^75x^3 + x^4 = 0x40 + 0x78x + 0x36x^2 + 0x0fx^3 + x^4
	fmt.Println("g(x) =", g)
	if !reflect.DeepEqual(g, []GF{0x40, 0x78, 0x36, 0x0f, 0x01}) {
		panic("g(x) is wrong")
	}

	// encode: RS(28, 24) with parity in the middle
	// x^{-12} R(x) mod g(x) where R(x) = upper * x^16 + lower mod g(x)
	msg := []byte("Hello, World!")
	I := make([]GF, 28)
	for i := 0; i < 12; i++ {
		if i < len(msg) {
			I[i] = GF(msg[i])
		}
	}
	for i := 12; i < 24; i++ {
		if i < len(msg) {
			I[i+4] = GF(msg[i])
		}
	}
	_, R := polyDiv(I, g)
	xinv := make([]GF, 4)
	{
		g0inv := g[0].inv()
		for i := 0; i < 4; i++ {
			xinv[i] = g[i+1].mul(g0inv)
		}
	}
	for i := 0; i < 12; i++ {
		_, R = polyDiv(polyMul(R, xinv), g)
	}
	fmt.Println("R(x) =", R)
	copy(I[12:16], R)
	fmt.Println("I(x) =", I)
	if _, rem := polyDiv(I, g); !reflect.DeepEqual(rem, make([]GF, 4)) {
		panic("I(x) is not a multiple of g(x)")
	}

	// perturb
	perturbations := [][]struct {
		pos int
		val GF
	}{
		{},
		{{pos: 4, val: 0x10}},
		{{pos: 4, val: 0x10}, {pos: 8, val: 0x20}},
	}
	for _, perturbation := range perturbations {
		copyI := make([]GF, 28)
		copy(copyI, I)
		for _, p := range perturbation {
			copyI[p.pos] ^= p.val
		}
		syndrome := make([]GF, 4)
		for i := 0; i < 4; i++ {
			syndrome[i] = polyOf(copyI, table[i])
		}
		fmt.Println("syndrome =", syndrome)

		// directly find the error locator polynomial
		// only fine for 2 errors
		if !rs2824Direct2(copyI, syndrome, table) {
			fmt.Println("Direct2 failed")
			rs2824Direct1(copyI, syndrome, table)
		}

		// Euclidean algorithm
		rs2824Euclidean(copyI, syndrome, table)
	}
}
