package api

import (
	"errors"
	"io"
	"net/http"

	"pavel-fokin/images-storage/internal/server/httputil"
)

type ImagesAdder interface {
	Add(io.Reader, string) error
}

type ImagesLister interface {
	List() error
}

var (
	ErrValidate = errors.New("'Content-Type' is required")
)

func StatusOK(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func ImagesGet(images ImagesLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func ImagesPost(images ImagesAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		contenttype, ok := r.Header["Content-Type"]
		if !ok {
			httputil.AsErrorResponse(w, ErrValidate, http.StatusBadRequest)
			return
		}

		images.Add(r.Body, contenttype[0])

		w.WriteHeader(http.StatusOK)
	}
}
