package stroage

import (
	"context"
	"io"
)

type Storage interface {
	Put(ctx context.Context, namespace string, objname string, obj io.Reader) error
	Get(ctx context.Context, namespace string, objname string) (io.Reader, error)
	List(ctx context.Context, namespace string) ([]string, error)
}
