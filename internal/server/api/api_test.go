package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Images struct {
	mock.Mock
}

func (m *Images) Add(io.Reader, string) error {
	m.Called()
	return nil
}

func (m *Images) List() error {
	m.Called()
	return nil
}

func Test_ImagesGet(t *testing.T) {
	// setup
	req, _ := http.NewRequest("", "", nil)
	w := httptest.NewRecorder()

	images := &Images{}
	images.On("List").Return()

	// test
	ImagesGet(images)(w, req)

	// assert
	resp := w.Result()
	assert.Equal(t, 200, resp.StatusCode)
}

func Test_ImagesPostValidationError(t *testing.T) {
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

func Test_ImagesPostSuccess(t *testing.T) {
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
	assert.Equal(t, 200, resp.StatusCode)
}
