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

func (m *StorageMock) Upload(
	ctx context.Context,
	uuid string,
	contenttype string,
	data io.Reader,
	metadata map[string]string,
) (Image, error) {
	m.Called()
	return Image{}, nil
}

func (m *StorageMock) DoesExist(
	ctx context.Context, uuid string,
) bool {
	m.Called()
	return true
}

func (m *StorageMock) Download(
	ctx context.Context, uuid string,
) (io.Reader, string, error) {
	m.Called()
	return bytes.NewReader([]byte{}), "", nil
}

func (m *StorageMock) Metadata(
	ctx context.Context, uuid string,
) (Image, error) {
	m.Called()
	return Image{}, nil
}

func Test_ImagesAdd(t *testing.T) {
	// setup
	storage := &StorageMock{}
	storage.On("Upload").Return()

	images := New(storage)

	// test
	_, err := images.Add(context.TODO(), bytes.NewReader([]byte{}), "image/png")

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
