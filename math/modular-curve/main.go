package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
)

func main() {
	delta := new(Laurent).JInv(20 - 1)
	val, _ := delta.Val()
	for i := val; i < val+delta.Prec(); i++ {
		fmt.Printf("%d => %v\n", i, delta.Coef(i))
	}
	degrees := []int{2, 3, 5, 7, 11, 13}
	for _, degree := range degrees {
		coefs := ModularBrute(degree)
		type Entry struct {
			XDeg int    `json:"x_deg"`
			YDeg int    `json:"y_deg"`
			Coef string `json:"coef"`
		}
		entries := make([]Entry, len(coefs))
		for i, coef := range coefs {
			entries[i] = Entry{
				XDeg: coef.XDeg,
				YDeg: coef.YDeg,
				Coef: coef.Coef.String(),
			}
		}
		slices.Reverse(entries)
		jsonString, err := json.MarshalIndent(entries, "", "  ")
		if err != nil {
			panic(err)
		}
		jsonString = append(jsonString, '\n')
		filename := fmt.Sprint(degree, ".json")
		os.WriteFile(filename, jsonString, 0o644)
		log.Println("wrote", filename)
	}

}
