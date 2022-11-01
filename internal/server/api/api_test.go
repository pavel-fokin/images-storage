package api

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/pavel-fokin/images-storage/internal/imagesstorage"
)

type Images struct {
	mock.Mock
}

func (m *Images) Add(context.Context, io.Reader, string) (imagesstorage.Image, error) {
	m.Called()
	return imagesstorage.Image{}, nil
}

func (m *Images) Update(context.Context, string, io.Reader, string) (imagesstorage.Image, error) {
	m.Called()
	return imagesstorage.Image{}, nil
}

func (m *Images) Data(context.Context, string, imagesstorage.BBox) (io.Reader, string, error) {
	args := m.Called()
	return args.Get(0).(io.Reader), args.Get(1).(string), args.Error(2)
}

func (m *Images) Metadata(context.Context, string) (imagesstorage.Image, error) {
	args := m.Called()
	return args.Get(0).(imagesstorage.Image), args.Error(1)
}

func (m *Images) List(context.Context) ([]imagesstorage.Image, error) {
	args := m.Called()
	return args.Get(0).([]imagesstorage.Image), args.Error(1)
}

func Test_ImagesGet_Success(t *testing.T) {
	// setup
	req, _ := http.NewRequest("", "", nil)
	w := httptest.NewRecorder()

	images := &Images{}
	images.On("List").Return([]imagesstorage.Image{}, nil)

	// test
	ImagesGet(images)(w, req)

	// assert
	resp := w.Result()
	assert.Equal(t, 200, resp.StatusCode)
}

func Test_ImagesGet_Failure(t *testing.T) {
	// setup
	req, _ := http.NewRequest("", "", nil)
	w := httptest.NewRecorder()

	images := &Images{}
	images.On("List").Return([]imagesstorage.Image{}, errors.New("error"))

	// test
	ImagesGet(images)(w, req)

	// assert
	resp := w.Result()
	assert.Equal(t, 500, resp.StatusCode)
}

func Test_ImagesGetByID_Success(t *testing.T) {
	// setup
	req, _ := http.NewRequest("", "", nil)
	w := httptest.NewRecorder()

	images := &Images{}
	images.On("Metadata").Return(imagesstorage.Image{}, nil)

	// test
	ImagesGetByID(images)(w, req)

	// assert
	resp := w.Result()
	assert.Equal(t, 200, resp.StatusCode)
}

func Test_ImagesGetByID_NotFound(t *testing.T) {
	// setup
	req, _ := http.NewRequest("", "", nil)
	w := httptest.NewRecorder()

	images := &Images{}
	images.On("Metadata").Return(imagesstorage.Image{}, imagesstorage.ErrImageNotExist)

	// test
	ImagesGetByID(images)(w, req)

	// assert
	resp := w.Result()
	assert.Equal(t, 404, resp.StatusCode)
}

func Test_ImagesGetByID_UnkownError(t *testing.T) {
	// setup
	req, _ := http.NewRequest("", "", nil)
	w := httptest.NewRecorder()

	images := &Images{}
	images.On("Metadata").Return(imagesstorage.Image{}, errors.New("error"))

	// test
	ImagesGetByID(images)(w, req)

	// assert
	resp := w.Result()
	assert.Equal(t, 500, resp.StatusCode)
}

func Test_ImagesGetDataByID_Success(t *testing.T) {
	// setup
	req, _ := http.NewRequest("", "", nil)
	w := httptest.NewRecorder()

	images := &Images{}
	images.On("Data").Return(bytes.NewBuffer([]byte{}), "image/png", nil)

	// test
	ImagesGetDataByID(images)(w, req)

	// assert
	resp := w.Result()
	assert.Equal(t, 200, resp.StatusCode)
}

func Test_ImagesGetDataByID_WithBBox_Success(t *testing.T) {
	// setup
	req, _ := http.NewRequest("", "v1/images?bbox=0,0,10,10", nil)
	w := httptest.NewRecorder()

	images := &Images{}
	images.On("Data").Return(bytes.NewBuffer([]byte{}), "image/png", nil)

	// test
	ImagesGetDataByID(images)(w, req)

	// assert
	resp := w.Result()
	assert.Equal(t, 200, resp.StatusCode)
}

func Test_ImagesGetDataByID_WithBBox_Validation(t *testing.T) {
	// setup
	tests := []struct {
		url            string
		expectedStatus int
	}{
		{"v1/images?bbox=0,0,0,0", 400},
		{"v1/images?bbox=0,0,a,b", 400},
		{"v1/images?bbox=0,0", 400},
		{"v1/images?bbox=0,0,10,10", 200},
	}

	images := &Images{}
	images.On("Data").Return(bytes.NewBuffer([]byte{}), "image/png", nil)

	for _, test := range tests {
		req, _ := http.NewRequest("", test.url, nil)
		w := httptest.NewRecorder()

		// test
		ImagesGetDataByID(images)(w, req)

		// assert
		resp := w.Result()
		assert.Equal(t, test.expectedStatus, resp.StatusCode)
	}
}

func Test_ImagesPost_ValidationError(t *testing.T) {
	// setup
	req, _ := http.NewRequest("", "", nil)
	w := httptest.NewRecorder()

	images := &Images{}
	images.On("Add").Return()

	// test
	ImagesPost(images)(w, req)

	// assert
	resp := w.Result()
	assert.Equal(t, 400, resp.StatusCode)
}

func Test_ImagesPost_Success(t *testing.T) {
	// setup
	req, _ := http.NewRequest("", "", nil)
	req.Header["Content-Type"] = []string{"image/png"}

	w := httptest.NewRecorder()

	images := &Images{}
	images.On("Add").Return()

	// test
	ImagesPost(images)(w, req)

	// assert
	resp := w.Result()
	assert.Equal(t, 201, resp.StatusCode)
}
