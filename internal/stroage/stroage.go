package stroage

import "io"

type Storage interface {
	Put(namespace string, objname string, obj io.Reader) error
	Get(namespace string, objname string) (io.Reader, error)
	List(namespace string) ([]string, error)
}
