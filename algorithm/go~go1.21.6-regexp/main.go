package main

import (
	"fmt"
	"regexp/syntax"
)

func main() {
	result, err := syntax.Parse("a{2,}.*", 0)
	if err != nil {
		panic(err)
	}
	result = result.Simplify()
	fmt.Printf("%#v\n", result)
	prog, err := syntax.Compile(result)
	if err != nil {
		panic(err)
	}
	fmt.Println(prog.String())
}
