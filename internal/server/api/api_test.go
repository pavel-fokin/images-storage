package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/mock"
)

func TestImagesGet(t *testing.T) {
	// setup
	req, _ := http.NewRequest("", "", nil)
	w := httptest.NewRecorder()

	// test
	ImagesGet()(w, req)

	// assert
	resp := w.Result()
	assert.Equal(t, 200, resp.StatusCode)
}

func TestImagesPost(t *testing.T) {
	// setup
	req, _ := http.NewRequest("", "", nil)
	w := httptest.NewRecorder()

	// test
	ImagesPost()(w, req)

	// assert
	resp := w.Result()
	assert.Equal(t, 200, resp.StatusCode)
}
