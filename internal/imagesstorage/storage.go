package imagesstorage

import (
	"context"
	"io"
)

type Storage interface {
	List(ctx context.Context) ([]Image, error)
	Upload(
		ctx context.Context, filename string, contenttype string, data io.Reader,
	) (Image, error)
	Metadata(
		ctx context.Context, uuid string,
	) (Image, error)
}
