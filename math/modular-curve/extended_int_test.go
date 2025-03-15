package main

import (
	"math/big"
	"testing"
)

func TestXIntAdd(t *testing.T) {
	a := new(XInt).SetInt(big.NewInt(1), 1)
	b := new(XInt).SetInt(big.NewInt(2), 1)
	c := new(XInt).Add(a, b)
	if c.Coef(0).Cmp(big.NewInt(3)) != 0 {
		t.Errorf("expected 3, got %v", c.Coef(0))
	}
}

func TestXIntSub(t *testing.T) {
	a := new(XInt).SetInt(big.NewInt(1), 1)
	b := new(XInt).SetInt(big.NewInt(2), 1)
	c := new(XInt).Sub(a, b)
	if c.Coef(0).Cmp(big.NewInt(-1)) != 0 {
		t.Errorf("expected -1, got %v", c.Coef(0))
	}
}

func TestXIntEq(t *testing.T) {
	a := new(XInt).SetInt(big.NewInt(1), 1)
	b := new(XInt).SetInt(big.NewInt(1), 1)
	if !a.Eq(b) {
		t.Errorf("expected true, got false")
	}
	a = new(XInt).SetInt(big.NewInt(1), 1)
	b = new(XInt).SetInt(big.NewInt(0), 1)
	if a.Eq(b) {
		t.Errorf("expected false, got true")
	}
}
