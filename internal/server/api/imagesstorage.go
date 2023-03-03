package api

import (
	"context"
	"io"

	"github.com/pavel-fokin/images-storage/internal/imagesstorage"
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
