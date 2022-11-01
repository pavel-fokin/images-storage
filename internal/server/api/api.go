package api

import (
	"context"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juju/errors"

	"github.com/pavel-fokin/images-storage/internal/imagesstorage"
	"github.com/pavel-fokin/images-storage/internal/log"
	"github.com/pavel-fokin/images-storage/internal/server/httputil"
)

var (
	ErrValidate = errors.New("'Content-Type' is required")
	ErrUpload   = errors.New("Coudn't upload an image")
	ErrNotFound = errors.New("Image not found")
	ErrUnknown  = errors.New(`Unkown error ¯\_(ツ)_/¯`)
)

type ImagesStorage interface {
	List(ctx context.Context) ([]imagesstorage.Image, error)
	Add(
		ctx context.Context, data io.Reader, contenttype string,
	) (imagesstorage.Image, error)
	Update(
		ctx context.Context, uuid string, data io.Reader, contenttype string,
	) (imagesstorage.Image, error)
	Metadata(
		ctx context.Context, uuid string,
	) (imagesstorage.Image, error)
	Data(
		ctx context.Context, uuid string, bbox imagesstorage.BBox,
	) (data io.Reader, contenttype string, err error)
}

func StatusOK(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func ImagesGet(images ImagesStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		images, err := images.List(r.Context())
		if err != nil {
			switch {
			case errors.Is(err, imagesstorage.ErrImageNotExist):
				httputil.AsErrorResponse(w, ErrNotFound, http.StatusNotFound)
			default:
				log.Error(r.Context(), err, "")
				httputil.AsErrorResponse(w, ErrUpload, http.StatusInternalServerError)
			}
			return
		}

		resp := asImagesGetResponse(images)

		httputil.AsSuccessResponse(w, resp, http.StatusOK)
	}
}

func ImagesPutByID(images ImagesStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		contenttype, ok := r.Header["Content-Type"]
		if !ok {
			httputil.AsErrorResponse(w, ErrValidate, http.StatusBadRequest)
			return
		}

		image, err := images.Update(r.Context(), id, r.Body, contenttype[0])
		if err != nil {
			switch {
			case errors.Is(err, imagesstorage.ErrImageNotExist):
				httputil.AsErrorResponse(w, ErrNotFound, http.StatusNotFound)
			default:
				log.Error(r.Context(), err, "")
				httputil.AsErrorResponse(w, ErrUpload, http.StatusInternalServerError)
			}
			return
		}

		resp := asImagesPostResponse(image)

		httputil.AsSuccessResponse(w, resp, http.StatusOK)
	}
}

func ImagesGetByID(images ImagesStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		image, err := images.Metadata(r.Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, imagesstorage.ErrImageNotExist):
				httputil.AsErrorResponse(w, ErrNotFound, http.StatusNotFound)
			default:
				log.Error(r.Context(), err, "")
				httputil.AsErrorResponse(w, ErrUpload, http.StatusInternalServerError)

			}
			return
		}

		resp := asImagesPostResponse(image)

		httputil.AsSuccessResponse(w, resp, http.StatusOK)
	}
}

func ImagesGetDataByID(images ImagesStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		// TODO Add validation
		bboxParam := r.URL.Query().Get("bbox")
		bbox := imagesstorage.BBox{}
		if bboxParam != "" {
			bbox.X, bbox.Y, bbox.W, bbox.H = ParseBBoxParam(bboxParam)
		}

		// TODO More errors handling
		data, contenttype, err := images.Data(r.Context(), id, bbox)
		if err != nil {
			switch {
			case errors.Is(err, imagesstorage.ErrImageNotExist):
				httputil.AsErrorResponse(w, ErrNotFound, http.StatusNotFound)
			default:
				log.Error(r.Context(), err, "")
				httputil.AsErrorResponse(w, ErrUpload, http.StatusInternalServerError)
			}
			return
		}

		httputil.AsDataRepsonse(w, contenttype, data)
	}
}

func ImagesPost(images ImagesStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// TODO More validation
		contenttype, ok := r.Header["Content-Type"]
		if !ok {
			httputil.AsErrorResponse(w, ErrValidate, http.StatusBadRequest)
			return
		}

		image, err := images.Add(r.Context(), r.Body, contenttype[0])
		if err != nil {
			log.Error(r.Context(), err, "")
			httputil.AsErrorResponse(w, ErrUpload, http.StatusInternalServerError)
			return
		}

		resp := asImagesPostResponse(image)

		httputil.AsSuccessResponse(w, resp, http.StatusCreated)
	}
}
