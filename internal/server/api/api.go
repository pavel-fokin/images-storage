package api

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"

	"pavel-fokin/images-storage/internal/imagesstorage"
	"pavel-fokin/images-storage/internal/server/httputil"
)

type ImagesStorage interface {
	List(ctx context.Context) ([]imagesstorage.Image, error)
	Add(data io.Reader, contenttype string) error
}

var (
	ErrValidate = errors.New("'Content-Type' is required")
)

func StatusOK(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func ImagesGet(images ImagesStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		images, err := images.List(r.Context())
		if err != nil {
			log.Fatal(err)
		}

		resp := asImagesGetResponse(images)

		httputil.AsSuccessResponse(w, resp, http.StatusOK)
	}
}

func ImagesPost(images ImagesStorage) http.HandlerFunc {
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
