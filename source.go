package ini

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sync/atomic"
)

var errSourceLocked = errors.New("ini: the data source was locked")

type DataSource interface {
	Open() (io.ReadCloser, error)
}

type dataSource struct {
	lock       int32
	readCloser io.ReadCloser
	reader     io.Reader
	bytes      []byte
	path       string
	source     DataSource
	factory    func() (io.ReadCloser, error)
}

func (s *dataSource) Lock() {
	atomic.StoreInt32(&s.lock, 1)
}

func (s *dataSource) Unlock() {
	atomic.StoreInt32(&s.lock, 0)
}

func (s *dataSource) Open() (io.ReadCloser, error) {
	if atomic.LoadInt32(&s.lock) == 1 {
		return nil, errSourceLocked
	}
	if s.readCloser != nil {
		return s.readCloser, nil
	}
	if s.reader != nil {
		return io.NopCloser(s.reader), nil
	}
	if s.bytes != nil {
		return io.NopCloser(bytes.NewReader(s.bytes)), nil
	}
	if s.path != "" {
		return os.Open(s.path)
	}
	if s.factory != nil {
		return s.factory()
	}
	return s.source.Open()
}

func (s *dataSource) reload(m *Manager) error {
	rc, err := s.Open()
	if err != nil {
		// In loose mode, we create an empty default section for nonexistent files.
		if os.IsNotExist(err) && m.options.Loose {
			return nil
		}
		if errors.Is(err, errSourceLocked) {
			return nil
		}
		return err
	}
	defer rc.Close()
	return m.parse(rc)
}

func parseDataSource(source any) (*dataSource, error) {
	switch s := source.(type) {
	case string:
		return &dataSource{path: s}, nil
	case []byte:
		return &dataSource{bytes: s}, nil
	case io.Reader:
		return &dataSource{reader: s}, nil
	case io.ReadCloser:
		return &dataSource{readCloser: s}, nil
	case DataSource:
		return &dataSource{source: s}, nil
	case func() (io.ReadCloser, error):
		return &dataSource{factory: s}, nil
	default:
		return nil, fmt.Errorf("error parsing data source: unknown type %q", s)
	}
}
