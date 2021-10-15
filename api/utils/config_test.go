package utils

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
)

type badReader struct{}

func (badReader) Read(p []byte) (int, error) {
	return 0, errors.New("test error")
}

type badWriter struct{}

func (badWriter) Write(p []byte) (int, error) {
	return 0, errors.New("test error")
}

type mockReaderWriter struct {
	nextReader io.Reader
	nextWriter io.Writer
	lastWrite  []byte
}

func (m *mockReaderWriter) setNextWriter(w io.Writer) {
	m.nextWriter = w
}

func (m *mockReaderWriter) setNextReader(r io.Reader) {
	m.nextReader = r
}

func (m *mockReaderWriter) setNextString(s string) {
	m.nextReader = strings.NewReader(s)
}

func (m *mockReaderWriter) Read(p []byte) (int, error) {
	if m.nextReader != nil {
		return m.nextReader.Read(p)
	}

	return 0, io.EOF
}

func (m *mockReaderWriter) Write(p []byte) (int, error) {
	if m.nextWriter != nil {
		return m.nextWriter.Write(p)
	}

	m.lastWrite = p
	return len(p), nil
}

func TestCreateDefaultConfig(t *testing.T) {
	rw := mockReaderWriter{}

	config, err := LoadOrCreateConfig(&rw)

	assert.Equal(t, err, nil)
	assert.Equal(t, config.Database.Path, "crateapi.db3")
	assert.Equal(t, config.Listener.Address, "0.0.0.0")
	assert.Equal(t, config.Listener.Port, uint(9080))

	assert.Equal(t, string(rw.lastWrite), `
[Database]
  Path = "crateapi.db3"

[Listener]
  Address = "0.0.0.0"
  Port = 9080
`)
}

func TestReadFullConfig(t *testing.T) {
	rw := mockReaderWriter{}

	rw.setNextString(`
[Database]
  Path = "other.db3"

[Listener]
  Address = "10.128.1.1"
  Port = 4343
`)

	config, err := LoadOrCreateConfig(&rw)

	assert.Equal(t, err, nil)
	assert.Equal(t, config.Database.Path, "other.db3")
	assert.Equal(t, config.Listener.Address, "10.128.1.1")
	assert.Equal(t, config.Listener.Port, uint(4343))
	assert.Equal(t, rw.lastWrite, nil)
}

func TestReadPartialConfig(t *testing.T) {
	rw := mockReaderWriter{}

	rw.setNextString(`
[Listener]
  Address = "10.128.1.1"
  Port = 4343
`)

	config, err := LoadOrCreateConfig(&rw)

	assert.Equal(t, err, nil)
	assert.Equal(t, config.Database.Path, "crateapi.db3")
	assert.Equal(t, config.Listener.Address, "10.128.1.1")
	assert.Equal(t, config.Listener.Port, uint(4343))
	assert.Equal(t, string(rw.lastWrite), `
[Database]
  Path = "crateapi.db3"

[Listener]
  Address = "10.128.1.1"
  Port = 4343
`)
}

func TestReadBadConfig(t *testing.T) {
	rw := mockReaderWriter{}

	rw.setNextString(`
[Database
  Path = other.db3"

[Listener]
  Address "10.128.1.1"
  Port = 4343
`)

	config, err := LoadOrCreateConfig(&rw)

	assert.NotEqual(t, err, nil)
	assert.Equal(t, config, nil)
	assert.Equal(t, rw.lastWrite, nil)
}

func TestReadErrorConfig(t *testing.T) {
	rw := mockReaderWriter{}
	rw.setNextReader(badReader{})

	config, err := LoadOrCreateConfig(&rw)

	assert.NotEqual(t, err, nil)
	assert.Equal(t, config, nil)
	assert.Equal(t, rw.lastWrite, nil)
}

func TestWriteErrorConfig(t *testing.T) {
	rw := mockReaderWriter{}
	rw.setNextWriter(badWriter{})

	config, err := LoadOrCreateConfig(&rw)

	assert.NotEqual(t, err, nil)
	assert.NotEqual(t, config, nil)
}
