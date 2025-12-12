package main

import (
	"fmt"
	"math"
)

func main() {
	sum := 0.0
	for i := 0; i < 10000000; i++ {
		sum += math.Sqrt(float64(i))
	}
	fmt.Printf("Result: %.2f\n", sum)
}
