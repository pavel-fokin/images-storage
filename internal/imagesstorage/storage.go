package imagesstorage

import (
	"context"
	"io"
)

type Storage interface {
	List(ctx context.Context) ([]Image, error)
	Upload(
		ctx context.Context, uuid string, contenttype string, data io.Reader,
	) (Image, error)
	DoesExist(
		ctx context.Context, uuid string,
	) bool
	Download(
		ctx context.Context, uuid string,
	) (io.Reader, string, error)
	Metadata(
		ctx context.Context, uuid string,
	) (Image, error)
}
