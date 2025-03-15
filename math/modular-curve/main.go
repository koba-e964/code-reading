package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	delta := new(Laurent).JInv(20 - 1)
	val, _ := delta.Val()
	for i := val; i < val+delta.Prec(); i++ {
		fmt.Printf("%d => %v\n", i, delta.Coef(i).Coef(0))
	}
	degrees := []int{2, 3, 5, 7, 11, 13, 17, 19}
	for _, degree := range degrees {
		start := time.Now()
		coefs := ModularBrute(degree)
		elapsed := time.Since(start)
		type Entry struct {
			XDeg int    `json:"x_deg"`
			YDeg int    `json:"y_deg"`
			Coef string `json:"coef"`
		}
		type Output struct {
			Entries []Entry `json:"entries"`
			TimeMs  int64   `json:"time_ms"`
		}
		entries := make([]Entry, len(coefs))
		for i, coef := range coefs {
			entries[i] = Entry{
				XDeg: coef.XDeg,
				YDeg: coef.YDeg,
				Coef: coef.Coef.String(),
			}
		}
		output := Output{
			Entries: entries,
			TimeMs:  elapsed.Milliseconds(),
		}
		jsonString, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			panic(err)
		}
		jsonString = append(jsonString, '\n')
		filename := fmt.Sprint(degree, ".json")
		os.WriteFile(filename, jsonString, 0o644)
		log.Println("wrote", filename)
	}
}
