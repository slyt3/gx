package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/slyt3/gx/cache"
)

// Watch runs a Go script and automatically reruns it when the file changes.
func Watch(scriptPath string, args []string) error {
	// Get absolute path
	path, err := filepath.Abs(scriptPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// Create file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	// Add file to watch
	err = watcher.Add(path)
	if err != nil {
		return fmt.Errorf("failed to watch file: %w", err)
	}

	fmt.Printf("Watching %s for changes... (Ctrl+C to stop)\n", path)

	// Channel to signal script restart
	restart := make(chan bool)

	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Watch for file changes in background (ONLY ONCE)
	go func() {
		var lastRestart time.Time
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					// Ignore if we just restarted (within 500ms)
					if time.Since(lastRestart) < 500*time.Millisecond {
						continue
					}
					fmt.Println("\nFile changed, restarting...")
					lastRestart = time.Now()
					restart <- true
				}
			case err := <-watcher.Errors:
				if err != nil {
					fmt.Fprintf(os.Stderr, "Watcher error: %v\n", err)
				}
			}
		}
	}()

	// Track the currently running process
	var currentCmd *exec.Cmd
	var cmdMutex sync.Mutex

	// Function to run the script
	runScript := func() {
		cmdMutex.Lock()
		defer cmdMutex.Unlock()

		// Kill existing process if running
		if currentCmd != nil && currentCmd.Process != nil {
			currentCmd.Process.Kill()
			currentCmd.Wait()
		}

		// Get file info for caching
		fileInfo, err := os.Stat(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}

		modTime := fileInfo.ModTime()
		fileSize := fileInfo.Size()

		// Check cache
		binaryPath, found, err := cache.Check(path, fileSize, modTime)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cache error: %v\n", err)
			return
		}

		// If not cached, compile
		if !found {
			tmpFile, err := os.CreateTemp("", "gx-watch-*.bin")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return
			}
			tmpPath := tmpFile.Name()
			tmpFile.Close()

			buildCmd := exec.Command("go", "build", "-o", tmpPath, path)
			buildCmd.Stderr = os.Stderr
			if err := buildCmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Compilation failed: %v\n", err)
				return
			}

			// Store in cache
			cache.Store(path, tmpPath, modTime, fileSize)
			binaryPath = tmpPath
		}

		// Run the binary
		currentCmd = exec.Command(binaryPath, args...)
		currentCmd.Stdout = os.Stdout
		currentCmd.Stderr = os.Stderr

		go func() {
			currentCmd.Run()
		}()
	}

	// Run initially
	runScript()

	// Main loop - wait for restart signals or Ctrl+C
	for {
		select {
		case <-restart:
			runScript()
		case <-sigChan:
			fmt.Println("\nStopping watch...")
			cmdMutex.Lock()
			if currentCmd != nil && currentCmd.Process != nil {
				currentCmd.Process.Kill()
			}
			cmdMutex.Unlock()
			return nil
		}
	}
}
