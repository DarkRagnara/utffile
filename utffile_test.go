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
