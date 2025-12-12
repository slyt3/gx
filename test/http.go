package main

import (
	"fmt"
	"net/http"
)

func main() {
	resp, err := http.Get("https://example.com")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Status:", resp.Status)
}
