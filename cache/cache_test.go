package cache

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGetCacheDir(t *testing.T) {
	dir, err := GetCacheDir()
	if err != nil {
		t.Fatalf("GetCacheDir failed: %v", err)
	}

	t.Logf("cache directory: %s", dir)

	if !strings.Contains(dir, "gx") {
		t.Errorf("Expected cache directory to contain 'gx', got: %s", dir)
	}
}

func TestGenerateCachkey(t *testing.T) {
	key1 := GenerateCacheKey("/home/user/script.go")
	key2 := GenerateCacheKey("/home/user/other/script.go")

	t.Logf("Key 1: %s", key1)
	t.Logf("Key 2: %s", key2)

	if key1 == key2 {
		t.Error("Expected different keys for different paths")
	}
}

func TestStore(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-binary-*.bin")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString("fake binary content")
	tmpFile.Close()

	scriptPath := "/home/user/test.go"
	modTime := time.Now()
	fileSize := int64(len("fake binary content"))

	err = Store(scriptPath, tmpFile.Name(), modTime, fileSize)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	cacheDir, _ := GetCacheDir()
	cacheKey := GenerateCacheKey(scriptPath)

	binPath := filepath.Join(cacheDir, cacheKey+".bin")
	metaPath := filepath.Join(cacheDir, cacheKey+".meta")

	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		t.Error("Binary file was not created in cache")
	}

	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		t.Error("Meta file was not created in cache")
	}

	t.Logf("Cache files created successfully!")
	t.Logf("Binary: %s", binPath)
	t.Logf("Meta: %s", metaPath)
}

func TestCheck(t *testing.T) {
	// fake folder creation
	tmpFile, err := os.CreateTemp("", "test-*.bin")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("test content")
	tmpFile.Close()

	//  test data
	scriptPath := "/home/user/checktest.go"
	fileSize := int64(12)
	modTime := time.Now()

	// Store binary in cache
	err = Store(scriptPath, tmpFile.Name(), modTime, fileSize)
	if err != nil {
		t.Fatal(err)
	}

	// Should find cached binary
	binPath, valid, err := Check(scriptPath, fileSize, modTime)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}
	if !valid {
		t.Error("Expected cache to be valid")
	}
	if binPath == "" {
		t.Error("Expected binary path")
	}
	t.Logf("✓ Test 1 passed: Found cached binary at %s", binPath)

	// Should NOT find cached binary
	_, valid, err = Check(scriptPath, 999, modTime)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}
	if valid {
		t.Error("Expected cache to be invalid (wrong hash)")
	}
	t.Logf("✓ Test 2 passed: Correctly rejected wrong hash")

	// Should NOT find cached binary
	wrongTime := time.Now().Add(time.Hour)
	_, valid, err = Check(scriptPath, 999, wrongTime)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}
	if valid {
		t.Error("Expected cache to be invalid (wrong time)")
	}
	t.Logf("✓ Test 3 passed: Correctly rejected wrong time")
}

func TestClean(t *testing.T) {
	_, err := GetCacheDir()
	if err != nil {
		t.Fatal(err)
	}

	err = Clean()
	if err != nil {
		t.Fatalf("Clean failed: %v", err)
	}
	t.Log("Cache cleaned succesfully")
}
