package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello from gx!")
	fmt.Println("ARgs:", os.Args[1:])
	os.Exit(0)
}
