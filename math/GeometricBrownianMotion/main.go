package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
)

func geometricBrownianMotion(epoch int64, mu float64, sigma float64, s0 float64) float64 {
	// Generate a random number
	cur := 1.0
	for i := int64(0); i < epoch; i++ {
		cur = cur * (1 + mu + sigma*rand.NormFloat64())
	}
	return cur
}

func main() {
	trial := 1_000_000
	mu := 0.07 / 4
	sigma := 0.2 / 2
	ratios := []float64{0.05, 0.1, 0.5, 0.8, 1.0}
	div := int64(40)
	for _, ratio := range ratios {
		sum := 0.0
		sqsum := 0.0
		values := make([]float64, 0, trial)
		for i := 0; i < trial; i++ {
			val := geometricBrownianMotion(div, mu*ratio, sigma*ratio, 1)
			sum += val
			sqsum += val * val
			values = append(values, val)
		}
		slices.Sort(values)
		mean := sum / float64(trial)
		variance := sqsum/float64(trial) - mean*mean
		fmt.Println("ratio: ", ratio)
		fmt.Println("Mean: ", mean)
		fmt.Println("Variance: ", math.Sqrt(variance))
		fmt.Println("median: ", values[trial/2])
		fmt.Println("10%ile: ", values[trial/10])
	}
}
