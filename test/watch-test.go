package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Version 18 - Updated!")
	for i := 0; i < 10; i++ {
		fmt.Printf("Count: %d\n", i)
		time.Sleep(1 * time.Second)
	}
}
