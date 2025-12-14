package cmd

import (
	"fmt"

	"github.com/slyt3/gx/cache"
)

// Clean removes all cached compiled binaries
func Clean() error {

	err := cache.Clean()
	if err != nil {
		return fmt.Errorf("cant clean the cache: %v", err)
	}

	fmt.Println("Cache cleaned successfully")

	return nil
}
