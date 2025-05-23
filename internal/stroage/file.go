package stroage

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// FileStorage implements the Storage interface using the local filesystem
type FileStorage struct {
	basePath string
}

// NewFileStorage creates a new file storage client
func NewFileStorage(basePath string) (*FileStorage, error) {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0o750); err != nil {
		return nil, err
	}

	return &FileStorage{
		basePath: basePath,
	}, nil
}

// Put stores an object in the filesystem
func (s *FileStorage) Put(_ context.Context, namespace, objname string, obj io.Reader) error {
	// Create directory for namespace if it doesn't exist
	dir := filepath.Join(s.basePath, namespace)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return err
	}

	// Check if file already exists
	filePath := filepath.Join(dir, objname)
	if _, err := os.Stat(filePath); err == nil {
		return os.ErrExist
	}

	// Create parent directories for the file if needed
	fileDir := filepath.Dir(filePath)
	if err := os.MkdirAll(fileDir, 0o750); err != nil {
		return err
	}

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Printf("Error closing file: %v", closeErr)
			if err == nil {
				err = closeErr
			}
		}
	}()

	// Copy data to file
	_, err = io.Copy(file, obj)
	return err
}

// Get retrieves an object from the filesystem
func (s *FileStorage) Get(_ context.Context, namespace, objname string) (io.Reader, error) {
	filePath := filepath.Join(s.basePath, namespace, objname)
	return os.Open(filePath)
}

// List returns all objects in a namespace from the filesystem
func (s *FileStorage) List(_ context.Context, namespace string) ([]string, error) {
	var objects []string
	dir := filepath.Join(s.basePath, namespace)

	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return objects, nil
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Remove the base path and namespace prefix
		relativePath := strings.TrimPrefix(path, filepath.Join(s.basePath, namespace)+string(filepath.Separator))
		objects = append(objects, relativePath)
		return nil
	})

	return objects, err
}
