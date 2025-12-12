package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//  TODO: add imports as we will need them later

// Holds information about a cached script
type CacheEntry struct {

	// Original script location
	ScriptPath string `json:"script_path"`

	// Hash of script content
	ScriptHash string `json:"script_hash"`

	// Unix timestamp of last modification
	ModTime time.Time `json:"mod_time"`

	// Its where the compiled binary will be stored
	BinaryPath string `json:"binary_path"`
}

// GetCacheDIr return cache directory path (~/.cache/gx/)
func GetCacheDir() (string, error) {

	// get home directory
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)

	}

	// combines path
	cachePath := filepath.Join(dir, ".cache", "gx")

	// Creates are directory
	err = os.MkdirAll(cachePath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to make directory: %w", err)
	}

	return cachePath, nil

}

// GenerateCacheKey creates uniqeu cache key from script path
func GenerateCacheKey(scriptPath string) string {

	filename := filepath.Base(scriptPath)
	name := strings.TrimSuffix(filename, ".go")

	// SHA256 hash of the script path
	hasher := sha256.New()
	hasher.Write([]byte(scriptPath))
	hashBytes := hasher.Sum(nil)

	// hex to string and first 8 char
	hashStr := hex.EncodeToString(hashBytes)[:8]

	return name + "-" + hashStr
}

// Store saves a compiled binary and its metadata to cache
func Store(scriptPath string, compiledBinaryPath string, modTime time.Time, scriptHash string) error {

	// cache directory
	cacheDir, err := GetCacheDir()
	if err != nil {
		return fmt.Errorf("failed to get cache directory: %w", err)
	}

	cacheKey := GenerateCacheKey(scriptPath)

	binaryPath := filepath.Join(cacheDir, cacheKey+".bin")
	metaPath := filepath.Join(cacheDir, cacheKey+".meta")

	// compile binay for reading
	srcFile, err := os.Open(compiledBinaryPath)
	if err != nil {
		return fmt.Errorf("failed to open compiled binary: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(binaryPath)
	if err != nil {
		return fmt.Errorf("failed to create the file: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy binary to cache: %w", err)
	}

	entry := CacheEntry{
		ScriptPath: scriptPath,
		ScriptHash: scriptHash,
		ModTime:    modTime,
		BinaryPath: binaryPath,
	}

	jsonBytes, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("Failed to convert to JSON: %w", err)
	}

	err = os.WriteFile(metaPath, jsonBytes, 0644)
	if err != nil {
		return fmt.Errorf("Failed to write JSON to .meta file: %w", err)
	}

	return nil
}

// Check looks up a cached binary for the given script
func Check(scriptPath string, scriptHash string, modTime time.Time) (string, bool, error) {

	cacheDir, err := GetCacheDir()
	if err != nil {
		return "", false, fmt.Errorf("failed to get cache directory: %w", err)

	}

	cacheKey := GenerateCacheKey(scriptPath)
	metaPath := filepath.Join(cacheDir, cacheKey+".meta")

	_, err = os.Stat(metaPath)
	if os.IsNotExist(err) {
		return "", false, nil
	}

	data, err := os.ReadFile(metaPath)
	if err != nil {
		return "", false, fmt.Errorf("failed to read meta file: %w", err)
	}

	var entry CacheEntry
	err = json.Unmarshal(data, &entry)
	if err != nil {
		return "", false, fmt.Errorf("failed to parse meta file: %w", err)
	}

	if entry.ScriptHash != scriptHash {
		return "", false, nil
	}

	if !entry.ModTime.Equal(modTime) {
		return "", false, nil
	}

	return entry.BinaryPath, true, nil
}

func Clean() error {

	cacheDir, err := GetCacheDir()
	if err != nil {
		return fmt.Errorf("Failed to get cache directory: %w", err)
	}
	err = os.RemoveAll(cacheDir)
	if err != nil {
		return fmt.Errorf("Failed to remove whole directory: %w", err)
	}

	return nil
}
