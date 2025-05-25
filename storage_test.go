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

func TestMemoryCachePutGetDelete(t *testing.T) {
	mc := NewMemoryCache()
	ctx := context.Background()

	testKey := "testcert"
	testData := []byte("certificate data")

	// Test Put and Get
	if err := mc.Put(ctx, testKey, testData); err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	retrievedData, err := mc.Get(ctx, testKey)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !bytes.Equal(retrievedData, testData) {
		t.Errorf("Get returned wrong data: got %q, want %q", retrievedData, testData)
	}

	// Test Get for non-existent key
	_, err = mc.Get(ctx, "nonexistentkey")
	if err != os.ErrNotExist {
		t.Errorf("Get for non-existent key returned wrong error: got %v, want %v", err, os.ErrNotExist)
	}

	// Test Delete
	if err := mc.Delete(ctx, testKey); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = mc.Get(ctx, testKey)
	if err != os.ErrNotExist {
		t.Errorf("Get after Delete returned wrong error: got %v, want %v", err, os.ErrNotExist)
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
