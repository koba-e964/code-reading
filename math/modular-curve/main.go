package main

import "fmt"

func main() {
	delta := new(Laurent).JInv(20 - 1)
	for i := delta.Val(); i < delta.Val()+delta.Prec(); i++ {
		fmt.Printf("%d => %v\n", i, delta.Coef(i))
	}
}
