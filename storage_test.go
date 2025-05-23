package easyhttps

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFileCachePutCreatesDir(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "cachetest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	// Define a cache directory path that is a subdirectory of tmpDir but doesn't exist yet
	cacheDir := filepath.Join(tmpDir, "testcache")

	// Instantiate FileCache
	fc := FileCache{Dir: cacheDir}

	// Define test data
	testKey := "testfile"
	testData := []byte("hello world")

	// Call fc.Put. Check for errors.
	if err := fc.Put(context.Background(), testKey, testData); err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	// Verify the directory was created
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		t.Errorf("cache directory %s was not created", cacheDir)
	} else if err != nil {
		t.Errorf("failed to stat cache directory %s: %v", cacheDir, err)
	}

	// Verify the file was created
	filePath := filepath.Join(cacheDir, testKey)
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", filePath, err)
	}

	// Verify file content
	if !bytes.Equal(fileData, testData) {
		t.Errorf("file content mismatch: got %q, want %q", fileData, testData)
	}
}

func TestFileCachePutExistingDir(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "cachetest_existing")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	// Define a cache directory path and explicitly create it
	cacheDir := filepath.Join(tmpDir, "testcache_exists")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		t.Fatalf("Failed to pre-create cache dir %s: %v", cacheDir, err)
	}

	// Instantiate FileCache
	fc := FileCache{Dir: cacheDir}

	// Define test data
	testKey := "testfile_existing"
	testData := []byte("hello again")

	// Call fc.Put. Check for errors.
	if err := fc.Put(context.Background(), testKey, testData); err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	// Verify the file was created
	filePath := filepath.Join(cacheDir, testKey)
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", filePath, err)
	}

	// Verify file content
	if !bytes.Equal(fileData, testData) {
		t.Errorf("file content mismatch: got %q, want %q", fileData, testData)
	}
}
