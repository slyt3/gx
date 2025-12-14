package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/slyt3/gx/cache"
)

// Run executes a Go script with caching.
// Returns the exit code from the script.
func Run(scriptPath string, args []string) int {
	path, err := filepath.Abs(scriptPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	// Check if path actually exists
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	modTime := fileInfo.ModTime()
	fileSize := fileInfo.Size()

	hashedPath, found, err := cache.Check(path, fileSize, modTime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking cache: %v\n", err)
		return 1
	}

	if found {
		// Cache hit - run the cached binary
		cmd := exec.Command(hashedPath, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return 1
		}
		return 0
	}

	// Cache miss - compile the script
	tmpFile, err := os.CreateTemp("", "gx-*.bin")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating temp file: %v\n", err)
		return 1
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()

	buildCmd := exec.Command("go", "build", "-o", tmpPath, path)
	buildCmd.Stderr = os.Stderr
	err = buildCmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Compilation failed: %v\n", err)
		return 1
	}

	// Store compiled binary in cache
	err = cache.Store(path, tmpPath, modTime, fileSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to store in cache: %v\n", err)
		return 1
	}

	// Run the compiled binary
	cmd := exec.Command(tmpPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return 1
	}
	return 0
}
