package stroage

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/minio/minio-go/v7"
)

// TestStorageImplementations verifies that all implementations satisfy the Storage interface
func TestStorageImplementations(t *testing.T) {
	// This is a compile-time check to ensure implementations satisfy the Storage interface
	var _ Storage = (*FileStorage)(nil)
	var _ Storage = (*MinioStorage)(nil)
}

// TestStorageCommon contains common tests that should work on any Storage implementation
func TestStorageCommon(t *testing.T) {
	t.Run("FileStorage", func(t *testing.T) {
		// Create a temporary directory for testing
		tempDir, err := os.MkdirTemp("", "file-storage-test")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tempDir)

		storage, err := NewFileStorage(tempDir)
		if err != nil {
			t.Fatalf("Failed to create FileStorage: %v", err)
		}

		runCommonStorageTests(t, storage)
	})

	// Note: MinioStorage test is commented out as it requires a real Minio server
	// The MinioStorage implementation is tested separately with mocks in minio_test.go

	//Uncomment to test with a real Minio server
	t.Run("MinioStorage", func(t *testing.T) {

		// Skip if credentials are not available
		endpoint := "localhost:9000"
		accessKey := "minioadmin"
		secretKey := "minioadmin"
		bucketName := "depgit-testing"

		if endpoint == "" || accessKey == "" || secretKey == "" || bucketName == "" {
			t.Skip("Minio credentials not available, skipping test")
		}

		storage, err := NewMinioStorage(endpoint, accessKey, secretKey, bucketName, false)
		if err != nil {
			t.Fatalf("Failed to create MinioStorage: %v", err)
		}

		// Delete all objects in bucket before running tests
		objectCh := storage.client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{Recursive: true})
		for obj := range objectCh {
			err := storage.client.RemoveObject(context.Background(), bucketName, obj.Key, minio.RemoveObjectOptions{})
			if err != nil {
				t.Fatalf("Failed to delete object %s: %v", obj.Key, err)
			}
		}

		runCommonStorageTests(t, storage)
	})
}

func runCommonStorageTests(t *testing.T, storage Storage) {
	// Test basic operations
	namespace := "test-namespace-common"
	objName := "common-test-file.txt"
	content := "common test content"
	reader := strings.NewReader(content)

	// Test put
	err := storage.Put(namespace, objName, reader)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	// Test get
	result, err := storage.Get(namespace, objName)
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

	// Test list
	files, err := storage.List(namespace)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	found := false
	for _, file := range files {
		if file == objName {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("File %s not found in listing", objName)
	}

	// Test putting an existing file
	reader = strings.NewReader("new content")
	err = storage.Put(namespace, objName, reader)
	if err != os.ErrExist {
		t.Errorf("Expected os.ErrExist when putting existing file, got %v", err)
	}
}
