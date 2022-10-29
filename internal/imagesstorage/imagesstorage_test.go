package imagesstorage

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (m *StorageMock) List(ctx context.Context) ([]Image, error) {
	m.Called()
	return []Image{}, nil
}

func (m *StorageMock) Upload(ctx context.Context, filename string, contenttype string, data io.Reader) error {
	m.Called()
	return nil
}

func Test_ImagesAdd(t *testing.T) {
	// setup
	storage := &StorageMock{}
	storage.On("Upload").Return()

	images := New(storage)

	// test
	err := images.Add(context.TODO(), bytes.NewReader([]byte("")), "image/png")

	// assert
	assert.NoError(t, err)
}

func Test_ImagesList(t *testing.T) {
	// setup
	storage := &StorageMock{}
	storage.On("List").Return()

	images := New(storage)

	// test
	_, err := images.List(context.TODO())

	// assert
	assert.NoError(t, err)
}
