package char2

import (
	"reflect"
	"testing"
)

func TestLinearSolve0(t *testing.T) {
	tests := []struct {
		A  [][]byte
		b  []byte
		ok bool
	}{
		{
			A: [][]byte{
				{0x01},
			},
			b:  []byte{0x01},
			ok: true,
		},
		{
			A: [][]byte{
				{0x01},
			},
			b:  []byte{0x02},
			ok: false,
		},
		{
			A: [][]byte{
				{0x01},
				{0x02},
				{0x03},
			},
			b:  []byte{0x03},
			ok: true,
		},
		{
			A: [][]byte{
				{0x01},
				{0x02},
				{0x03},
			},
			b:  []byte{0x04},
			ok: false,
		},
		{
			A: [][]byte{
				{0xa0, 0xcd},
			},
			b:  []byte{0x5a, 0x5a},
			ok: false,
		},
		{
			A: [][]byte{
				{0xa0, 0xcd},
				{0xa0, 0xcb},
			},
			b:  []byte{0x00, 0x06},
			ok: true,
		},
		{
			// Example from a real use case.
			A: [][]byte{
				{0x6f, 0x4c, 0xa5, 0x9b},
				{0x9f, 0x9e, 0x3b, 0xec},
				{0x7f, 0x3b, 0x6, 0x3},
				{0xfe, 0x76, 0xc, 0x6},
				{0xfc, 0xed, 0x18, 0xc},
				{0xf8, 0xdb, 0x31, 0x18},
				{0xf0, 0xb7, 0x63, 0x30},
				{0xe0, 0x6f, 0xc7, 0x60},
			},
			b:  []byte{0x0, 0x0, 0x0, 0x0},
			ok: true,
		},
		{
			// Example from a real use case.
			A: [][]byte{
				{0x6f, 0x4c, 0xa5, 0x9b},
				{0x9f, 0x9e, 0x3b, 0xec},
				{0x7f, 0x3b, 0x6, 0x3},
				{0xfe, 0x76, 0xc, 0x6},
				{0xfc, 0xed, 0x18, 0xc},
				{0xf8, 0xdb, 0x31, 0x18},
				{0xf0, 0xb7, 0x63, 0x30},
				{0xe0, 0x6f, 0xc7, 0x60},
				{0xc0, 0xdf, 0x8e, 0xc1},
				{0xc1, 0xb9, 0x6c, 0x58},
				{0x82, 0x73, 0xd9, 0xb0},
				{0x45, 0xe1, 0xc3, 0xba},
				{0xcb, 0xc4, 0xf6, 0xae},
				{0xd7, 0x8f, 0x9c, 0x86},
				{0xef, 0x19, 0x48, 0xd6},
				{0x9f, 0x35, 0xe1, 0x77},
				{0x3e, 0x6b, 0xc2, 0xef},
				{0x3d, 0xd0, 0xf5, 0x4},
				{0x7a, 0xa0, 0xeb, 0x9},
				{0xf4, 0x40, 0xd7, 0x13},
				{0xe8, 0x81, 0xae, 0x27},
				{0xd0, 0x3, 0x5d, 0x4f},
				{0xa0, 0x7, 0xba, 0x9e},
				{0x1, 0x9, 0x5, 0xe6},
				{0x43, 0x14, 0x7b, 0x17},
				{0x86, 0x28, 0xf6, 0x2e},
				{0xc, 0x51, 0xec, 0x5d},
				{0x18, 0xa2, 0xd8, 0xbb},
				{0x71, 0x42, 0xc0, 0xac},
				{0xa3, 0x82, 0xf1, 0x82},
				{0x7, 0x3, 0x92, 0xde},
				{0x4f, 0x0, 0x55, 0x66},
			},
			b:  []byte{0x2, 0x9c, 0x89, 0x49},
			ok: true,
		},
	}
	for _, tt := range tests {
		got, ok := SolveLinear(tt.A, tt.b)
		if ok != tt.ok {
			t.Errorf("SolveLinear(%v, %v) = %v, %v; want %v", tt.A, tt.b, got, ok, tt.ok)
		}
		if ok {
			prod := make([]byte, len(tt.b))
			for i, row := range tt.A {
				if got[i] {
					for j, v := range row {
						prod[j] ^= v
					}
				}
			}
			if !reflect.DeepEqual(prod, tt.b) {
				t.Errorf("SolveLinear(%v, %v) = %v, %v; but %v * %v = %v, want %v", tt.A, tt.b, got, ok, got, tt.A, prod, tt.b)
			}
		}
	}
}
