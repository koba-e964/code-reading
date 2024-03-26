package char2

import "hash/crc32"

// GF(2^32), but in reversed representation.
//
//	0x8000_0000 = 1
//	0x4000_0000 = x
//	...
//	0x1 = x^31
type GF232 uint32

// Add sets g to a + b and returns g. Note that a - b = a + b.
func (g *GF232) Add(a, b GF232) *GF232 {
	*g = a ^ b
	return g
}

// MulX sets g to a * x and returns g.
func (g *GF232) MulX(a GF232) *GF232 {
	*g = a >> 1
	if a&1 != 0 {
		*g ^= crc32.IEEE
	}
	return g
}

// Mul sets g to a * b and returns g.
func (g *GF232) Mul(a, b GF232) *GF232 {
	*g = 0
	for i := 31; i >= 0; i-- {
		if b&(1<<i) != 0 {
			*g ^= a
		}
		a.MulX(a)
	}
	return g
}

// Pow sets g to a^n and returns g.
func (g *GF232) Pow(a GF232, n uint64) *GF232 {
	*g = 0x8000_0000
	for n > 0 {
		if n&1 != 0 {
			g.Mul(*g, a)
		}
		a.Mul(a, a)
		n >>= 1
	}
	return g
}

// Weight32 sets g to x^{8 * pos} and returns g.
func (g *GF232) Weight32(pos int) *GF232 {
	g.Pow(0x80_0000, uint64(pos))
	return g
}
