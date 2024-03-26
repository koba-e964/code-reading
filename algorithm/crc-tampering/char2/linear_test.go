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
