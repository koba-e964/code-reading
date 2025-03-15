package main

import (
	"math/big"
	"testing"
)

func TestMulShift(t *testing.T) {
	minusOne := new(XInt).SetInt(big.NewInt(-1), 1)
	a := new(Laurent).SetInt(new(XInt).SetInt(big.NewInt(1), 1), 3)
	a.SetCoef(1, new(XInt).SetInt(big.NewInt(2), 1))
	a.SetCoef(2, new(XInt).SetInt(big.NewInt(3), 1))
	c := new(Laurent).MulShift(a, minusOne, minusOne)
	if !c.Coef(0).Eq(new(XInt).SetInt(big.NewInt(1), 1)) {
		t.Errorf("expected 1, got %v", c.Coef(0))
	}
	if !c.Coef(1).Eq(new(XInt).SetInt(big.NewInt(-2), 1)) {
		t.Errorf("expected -2, got %v", c.Coef(1))
	}
	if !c.Coef(2).Eq(new(XInt).SetInt(big.NewInt(3), 1)) {
		t.Errorf("expected 3, got %v", c.Coef(2))
	}
}
