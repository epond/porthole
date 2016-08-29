package test

import (
	"testing"
	"strconv"
	"os"
	"time"
	"io"
)

func Expect(t *testing.T, valueName string, expected string, actual string) {
	if actual != expected {
		t.Errorf("Expected %v to be %v but was %v", valueName, expected, actual)
	}
}

func ExpectInt(t *testing.T, valueName string, expected int, actual int) {
	Expect(t, valueName, strconv.Itoa(expected), strconv.Itoa(actual))
}

type DummyFileInfo struct {
	CannedName  string
	CannedIsDir bool
}
func (d DummyFileInfo) Name() string {
	return d.CannedName
}
func (d DummyFileInfo) Size() int64 {
	return 1
}
func (d DummyFileInfo) Mode() os.FileMode {
	if d.IsDir() {
		return os.ModeDir
	}
	return os.ModePerm
}
func (d DummyFileInfo) ModTime() time.Time {
	return time.Now()
}
func (d DummyFileInfo) IsDir() bool {
	return d.CannedIsDir
}
func (d DummyFileInfo) Sys() interface{} {
	return nil
}

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}