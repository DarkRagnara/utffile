package utffile

import (
	"bytes"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"math"
	"os"
)

func Wrap(r io.Reader) io.Reader {
	switch r.(type) {
	case io.ReadCloser:
		return readcloser{newReader(r)}
	default:
		return newReader(r)
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
var utf16lebom = [2]byte{0xff, 0xfe}
var utf16bebom = [2]byte{0xfe, 0xff}

type reader struct {
	io.Reader
	initialRead *bool
	buf         *[]byte
	utf16reader **transform.Reader
}

func newReader(r io.Reader) reader {
	return reader{Reader: r, buf: new([]byte), initialRead: new(bool), utf16reader: new(*transform.Reader)}
}

func (r reader) Read(b []byte) (int, error) {
	if !*r.initialRead {
		*r.buf = make([]byte, 3)
		*r.initialRead = true
		n, err := r.Reader.Read((*r.buf)[:])
		*r.buf = (*r.buf)[:n]

		wanted := int(math.Min(float64(n), float64(len(b))))

		if err != nil || n < 2 {
			copy(b, (*r.buf)[:wanted])
			(*r.buf) = (*r.buf)[wanted:]
			return wanted, err
		}
		if equalSlice((*r.buf)[:], utf8bom[:]) {
			*r.buf = nil
		} else if equalSlice((*r.buf)[:2], utf16lebom[:]) {
			r.setUTF16Decoder(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM))
		} else if equalSlice((*r.buf)[:2], utf16bebom[:]) {
			r.setUTF16Decoder(unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM))
		}
	}
	if len(*r.buf) != 0 {
		return r.readFromBufFirst(b)
	} else if *r.utf16reader != nil {
		return (*r.utf16reader).Read(b)
	}
	return r.Reader.Read(b)
}

func (r reader) setUTF16Decoder(utf16 encoding.Encoding) {
	reader := io.MultiReader(bytes.NewReader((*r.buf)[2:]), r.Reader)
	*r.utf16reader = transform.NewReader(reader, utf16.NewDecoder())
	*r.buf = (*r.buf)[0:0]
}

func (r reader) readFromBufFirst(b []byte) (int, error) {
	wanted := int(math.Min(float64(len(*r.buf)), float64(len(b))))
	copy(b[:wanted], (*r.buf)[:wanted])
	*r.buf = (*r.buf)[wanted:]
	n, err := r.Reader.Read(b[wanted:])
	return n + wanted, err
}

func equalSlice(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type readcloser struct {
	reader
}

var _ io.ReadCloser = readcloser{}

func (rc readcloser) Close() error {
	return rc.reader.Reader.(io.Closer).Close()
}
