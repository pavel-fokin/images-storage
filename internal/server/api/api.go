package api

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"pavel-fokin/images-storage/internal/imagesstorage"
	"pavel-fokin/images-storage/internal/server/httputil"
)

var (
	ErrValidate = errors.New("'Content-Type' is required")
	ErrUpload   = errors.New("Coudn't upload an image")
)

type ImagesStorage interface {
	List(ctx context.Context) ([]imagesstorage.Image, error)
	Add(
		ctx context.Context, data io.Reader, contenttype string,
	) (imagesstorage.Image, error)
	Metadata(
		ctx context.Context, uuid string,
	) (imagesstorage.Image, error)
}

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

func ImagesGetByID(images ImagesStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		image, err := images.Metadata(r.Context(), id)
		if err != nil {
			log.Fatal(err)
		}

		resp := asImagesPostResponse(image)

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

		image, err := images.Add(r.Context(), r.Body, contenttype[0])
		if err != nil {
			httputil.AsErrorResponse(w, ErrUpload, http.StatusInternalServerError)
			return
		}

		resp := asImagesPostResponse(image)

		httputil.AsSuccessResponse(w, resp, http.StatusCreated)
	}
}
