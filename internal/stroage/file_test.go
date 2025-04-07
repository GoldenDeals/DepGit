package stroage

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileStorage(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "file-storage-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create FileStorage instance
	storage, err := NewFileStorage(tempDir)
	if err != nil {
		t.Fatalf("Failed to create FileStorage: %v", err)
	}

	t.Run("Put and Get", func(t *testing.T) {
		// Test data
		namespace := "test-namespace"
		objName := "test-file.txt"
		content := "hello world"
		reader := strings.NewReader(content)
		ctx := context.Background()

		// Put the object
		err := storage.Put(ctx, namespace, objName, reader)
		if err != nil {
			t.Fatalf("Put failed: %v", err)
		}

		// Verify file exists
		filePath := filepath.Join(tempDir, namespace, objName)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Fatalf("File wasn't created at expected path: %s", filePath)
		}

		// Get the object
		result, err := storage.Get(ctx, namespace, objName)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		// Read and verify content
		data, err := io.ReadAll(result)
		if err != nil {
			t.Fatalf("Failed to read result: %v", err)
		}

		if string(data) != content {
			t.Errorf("Expected content %q, got %q", content, string(data))
		}
	})

	t.Run("Put existing file", func(t *testing.T) {
		// Create test data and file
		namespace := "test-namespace"
		objName := "existing-file.txt"
		content := "existing content"
		reader := strings.NewReader(content)
		ctx := context.Background()

		// First Put should succeed
		err := storage.Put(ctx, namespace, objName, reader)
		if err != nil {
			t.Fatalf("First Put failed: %v", err)
		}

		// Second Put should fail with os.ErrExist
		reader = strings.NewReader("new content")
		err = storage.Put(ctx, namespace, objName, reader)
		if err != os.ErrExist {
			t.Errorf("Expected os.ErrExist, got %v", err)
		}
	})

	t.Run("Get non-existent file", func(t *testing.T) {
		ctx := context.Background()
		// Try to get a file that doesn't exist
		_, err := storage.Get(ctx, "non-existent", "file.txt")
		if err == nil {
			t.Errorf("Expected error when getting non-existent file, got nil")
		}
	})

	t.Run("List", func(t *testing.T) {
		// Create test namespace with multiple files
		namespace := "list-test"
		files := []string{"file1.txt", "file2.txt", "nested/file3.txt"}
		content := "test content"
		ctx := context.Background()

		// Put all test files
		for _, file := range files {
			err := storage.Put(ctx, namespace, file, strings.NewReader(content))
			if err != nil {
				t.Fatalf("Failed to put test file %s: %v", file, err)
			}
		}

		// List files
		listed, err := storage.List(ctx, namespace)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		// Verify file count
		if len(listed) != len(files) {
			t.Errorf("Expected %d files, got %d", len(files), len(listed))
		}

		// Verify all files are listed
		for _, file := range files {
			found := false
			for _, listedFile := range listed {
				if file == listedFile {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("File %s not found in listing", file)
			}
		}
	})

	t.Run("List empty namespace", func(t *testing.T) {
		ctx := context.Background()
		// List files in a namespace with no files
		listed, err := storage.List(ctx, "empty-namespace")
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(listed) != 0 {
			t.Errorf("Expected empty list, got %d items", len(listed))
		}
	})
}
