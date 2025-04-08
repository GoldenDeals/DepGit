package git

import (
	"io"
	"sync"

	"github.com/GoldenDeals/DepGit/internal/share/errors"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/packfile"
	"github.com/go-git/go-git/v5/storage/memory"
)

type memObjectMetadata struct {
	pos int64
	plumbing.MemoryObject
}

func (o *memObjectMetadata) SetPos(pos int64) {
	o.pos = pos
}

type objectObserver struct {
	count   uint32
	objects []plumbing.EncodedObject

	objHeader struct {
		t    plumbing.ObjectType
		size int64
		pos  int64
	}
	// IDK can git parcer parce object header and content in parallel
	// so we need to lock it
	objHeaderMu sync.Mutex
}

func (o *objectObserver) OnHeader(count uint32) error {
	o.objects = make([]plumbing.EncodedObject, 0, count)
	o.count = count

	return nil
}

func (o *objectObserver) OnInflatedObjectHeader(t plumbing.ObjectType, objSize int64, pos int64) error {
	o.objHeaderMu.Lock()
	defer o.objHeaderMu.Unlock()

	o.objHeader.t = t
	o.objHeader.size = objSize
	o.objHeader.pos = pos

	return nil
}

func (o *objectObserver) OnInflatedObjectContent(h plumbing.Hash, pos int64, crc uint32, content []byte) error {
	o.objHeaderMu.Lock()
	defer o.objHeaderMu.Unlock()

	obj := &memObjectMetadata{}
	obj.SetType(o.objHeader.t)
	obj.SetSize(o.objHeader.size)
	obj.SetPos(pos)

	n, err := obj.Write(content)
	if err != nil {
		return err
	}
	if n != len(content) {
		// TODO: better error
		return errors.ErrBadData
	}

	o.objects = append(o.objects, obj)

	return nil
}

func (o *objectObserver) OnFooter(h plumbing.Hash) error {
	return nil
}

// We are storing all the packfile's objects in-memory. Maybe it will be a problem
// for large repositories.
func ParsePackfile(re io.Reader) ([]plumbing.EncodedObject, error) {
	observer := &objectObserver{}
	parser, err := packfile.NewParserWithStorage(
		packfile.NewScanner(re),
		memory.NewStorage(),
		observer,
	)
	if err != nil {
		return nil, err
	}

	if _, err := parser.Parse(); err != nil {
		return nil, err
	}

	log.
		WithField("count", len(observer.objects)).
		Trace("Packfile parsed")

	return observer.objects, nil
}
