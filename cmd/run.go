package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/slyt3/gx/cache"
)

// Run execute a Go script with caching
// Return the exit code from the script
func Run(scriptPath string, args []string) int {
	path, err := filepath.Abs(scriptPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	// check if path actualy exists
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	modTime := fileInfo.ModTime()

	// reading file content
	fileContent, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		return 1
	}

	hasher := sha256.New()
	hasher.Write(fileContent)
	hashBytes := hasher.Sum(nil)
	scriptHash := hex.EncodeToString(hashBytes)[:8]

	hashedPath, found, err := cache.Check(scriptPath, scriptHash, modTime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error cache error: %v\n", err)
		return 1
	}

	if found {

		fmt.Println("Cache hit! Using cached binary")

		// runing the cache binary, command returnsthe cmd struct
  // to execute the named program with the given arguments
		cmd := exec.Command(hashedPath, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()

		// exit code
		if err != nil {
			return 1
		} else {
			return 0
		}
	}

	fmt.Println("Cache hit! Using cached binary")

	tmpFile, err := os.CreateTemp("", "gx-*.bin")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating temp file: %v\n", err)
		return 1
	}

	tmpPath := tmpFile.Name()
	tmpFile.Close()

	// Compilie the script
	buildCmd := exec.Command("go", "build", "-o", tmpPath, path)
	buildCmd.Stderr = os.Stderr
	err = buildCmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Compilation failed: %v\n", err)
		return 1
	}

	// Store compiled binary in cache
	// help me god
	err = cache.Store(path, tmpPath, modTime, scriptHash)
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
