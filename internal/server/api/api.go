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
	ErrValidate     = errors.New("'Content-Type' is required")
	ErrUpload       = errors.New("Coudn't upload an image")
	ErrNotFound     = errors.New("Image not found")
	ErrUnknown      = errors.New(`Unknown error ¯\_(ツ)_/¯`)
	ErrBBoxValidate = errors.New("'bbox' parse error")
)

// ImageStorage is an interface to images-storage functionality.
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

// @Summary   List metadata for stored images.
// @Tags      images-storage
// @Produce   json
// @Success   200        {object}  ImagesGetResp
// @Failure   400        {object}  httputil.ErrorResponse
// @Failure   500        {object}  httputil.ErrorResponse
// @Router    /v1/images [get]
func ImagesGet(images ImagesStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		images, err := images.List(r.Context())
		if err != nil {
			log.Error(r.Context(), err, "")
			httputil.AsErrorResponse(w, ErrUnknown, http.StatusInternalServerError)
			return
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
		bbox, err := ParseBBoxParam(r)
		if err != nil {
			httputil.AsErrorResponse(w, ErrBBoxValidate, http.StatusBadRequest)
			return
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
