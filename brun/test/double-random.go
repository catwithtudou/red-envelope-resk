package main

import (
	"fmt"
	"red-envelope/infra/algo"
)

func main() {
	count, amount := int64(10), int64(100)
	for i := int64(0); i < count; i++ {
		x := algo.DoubleRandom(count, amount*100)
		fmt.Print(x, ",")

	}
	fmt.Println()
}
