package images

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Storage struct {
	mock.Mock
}

func (m *Storage) List(ctx context.Context) error {
	m.Called()
	return nil
}

func (m *Storage) Upload(ctx context.Context, filename string, contenttype string, data io.Reader) error {
	m.Called()
	return nil
}

func Test_ImagesAdd(t *testing.T) {
	// setup
	storage := &Storage{}
	storage.On("Upload").Return()

	images := New(storage)

	// test
	err := images.Add(bytes.NewReader([]byte("")), "image/png")

	// assert
	assert.NoError(t, err)
}

func Test_ImagesList(t *testing.T) {
	// setup
	storage := &Storage{}
	storage.On("List").Return()

	images := New(storage)

	// test
	err := images.List(context.TODO())

	// assert
	assert.NoError(t, err)
}
