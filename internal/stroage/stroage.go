// Package stroage provides interfaces and implementations for object storage.
// It supports different storage backends like local filesystem and MinIO.
package stroage

import (
	"context"
	"io"
)

// Storage defines the interface for object storage operations.
// Implementations must support basic operations like Put, Get, and List
// while handling namespacing for object organization.
type Storage interface {
	// Put stores an object in the specified namespace with the given name.
	Put(ctx context.Context, namespace string, objname string, obj io.Reader) error

	// Get retrieves an object from the specified namespace by its name.
	Get(ctx context.Context, namespace string, objname string) (io.Reader, error)

	// List returns all objects in the specified namespace.
	List(ctx context.Context, namespace string) ([]string, error)
}
