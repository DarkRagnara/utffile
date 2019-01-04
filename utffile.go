package utffile

import (
	"io"
	"os"
)

func Wrap(r io.Reader) io.Reader {
	switch r.(type) {
	case io.ReadCloser:
		return readcloser{reader{Reader: r}}
	default:
		return reader{Reader: r}
	}
}

func Open(name string) (io.Reader, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return Wrap(file), nil
}

var utf8bom = [3]byte{0xef, 0xbb, 0xbf}

type reader struct {
	io.Reader
	initialRead bool
}

func (r reader) Read(b []byte) (int, error) {
	if !r.initialRead {
		r.initialRead = true
		var buf [3]byte
		n, err := r.Reader.Read(buf[:])
		if err != nil || n < 3 {
			copy(b, buf[:n])
			return n, err
		}
		if buf != utf8bom {
			copy(b[:3], buf[:n])
			n, err := r.Reader.Read(b[3:])
			return n + 3, err
		}
	}
	return r.Reader.Read(b)
}

type readcloser struct {
	reader
}

var _ io.ReadCloser = readcloser{}

func (rc readcloser) Close() error {
	return rc.reader.Reader.(io.Closer).Close()
}
