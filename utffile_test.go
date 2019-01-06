package utffile

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func TestEmptyFile(t *testing.T) {
	testFile(t, "samples/empty.txt", "")
}

func TestShortFile1(t *testing.T) {
	testFile(t, "samples/1byte.txt", "a")
}

func TestShortFile2(t *testing.T) {
	testFile(t, "samples/2byte.txt", "ab")
}

func TestShortFile3(t *testing.T) {
	testFile(t, "samples/3byte.txt", "abc")
}

func TestShortFile4(t *testing.T) {
	testFile(t, "samples/4byte.txt", "abcd")
}

func TestUTF8(t *testing.T) {
	testFile(t, "samples/utf8.txt", expectedString)
}

func TestUTF8BOM(t *testing.T) {
	testFile(t, "samples/utf8-bom.txt", expectedString)
}

func TestUTF16LE(t *testing.T) {
	testFile(t, "samples/utf16-le.txt", expectedString)
}

func TestUTF16BE(t *testing.T) {
	testFile(t, "samples/utf16-be.txt", expectedString)
}

func TestEmptyUTF8BOM(t *testing.T) {
	testFile(t, "samples/empty-utf8-bom.txt", "")
}

func TestEmptyUTF16LE(t *testing.T) {
	testFile(t, "samples/empty-utf16-le.txt", "")
}

func TestEmptyUTF16BE(t *testing.T) {
	testFile(t, "samples/empty-utf16-be.txt", "")
}

func TestEmptyFileSlowRead(t *testing.T) {
	testFileSlowRead(t, "samples/empty.txt", "")
}

func TestShortFile1SlowRead(t *testing.T) {
	testFileSlowRead(t, "samples/1byte.txt", "a")
}

func TestShortFile2SlowRead(t *testing.T) {
	testFileSlowRead(t, "samples/2byte.txt", "ab")
}

func TestShortFile3SlowRead(t *testing.T) {
	testFileSlowRead(t, "samples/3byte.txt", "abc")
}

func TestShortFile4SlowRead(t *testing.T) {
	testFileSlowRead(t, "samples/4byte.txt", "abcd")
}

func TestUTF8SlowRead(t *testing.T) {
	testFileSlowRead(t, "samples/utf8.txt", expectedString)
}

func TestUTF8BOMSlowRead(t *testing.T) {
	testFileSlowRead(t, "samples/utf8-bom.txt", expectedString)
}

func TestUTF16LESlowRead(t *testing.T) {
	testFileSlowRead(t, "samples/utf16-le.txt", expectedString)
}

func TestUTF16BESlowRead(t *testing.T) {
	testFileSlowRead(t, "samples/utf16-be.txt", expectedString)
}

func TestWrapCloser(t *testing.T) {
	sr := strings.NewReader("abc")
	noCloser := Wrap(sr)

	if _, ok := noCloser.(io.ReadCloser); ok {
		t.Fatal("io.Reader should not be wrapped into io.ReadCloser")
	}

	src := ioutil.NopCloser(sr)
	closer := Wrap(src)

	if _, ok := closer.(io.ReadCloser); !ok {
		t.Fatal("io.ReadCloser should be wrapped as io.ReadCloser")
	}
}

func testFileSlowRead(t *testing.T, name string, expected string) {
	file, err := Open(name)
	assertNoError(t, err)
	defer file.(io.ReadCloser).Close()

	var contents strings.Builder
	var buf [1]byte
	for {
		n, err := file.Read(buf[:])
		if n == 1 {
			contents.WriteByte(buf[0])
		}
		if err == io.EOF {
			break
		}
		assertNoError(t, err)
	}

	assertNoError(t, err)

	assertStringEqual(t, expected, contents.String())
}

func testFile(t *testing.T, name string, expected string) {
	file, err := Open(name)
	assertNoError(t, err)
	defer file.(io.ReadCloser).Close()

	contents, err := ioutil.ReadAll(file)
	assertNoError(t, err)

	assertStringEqual(t, expected, string(contents))
}

func assertStringEqual(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Fatalf("Expected '%v' %v, but got '%v' %v", expected, []byte(expected), actual, []byte(actual))
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

var expectedString = "Hello, viele Grüße & さよなら!"
