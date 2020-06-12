package main

import (
	"fmt"
	"red-envelope/infra/algo"
)

func main() {
	count, amount := int64(10), int64(100)
	for i := int64(0); i < count; i++ {
		x := algo.SimpleRand(count, amount*100)
		fmt.Print(float64(x)/float64(100), ",")
	}
	fmt.Println()
}
