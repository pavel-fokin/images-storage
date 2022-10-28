package imagesstorage

import (
	"context"
	"io"
)

type Storage interface {
	Lister
	Uploader
}

type Lister interface {
	List(ctx context.Context) ([]Image, error)
}

type Uploader interface {
	Upload(ctx context.Context, filename string, contenttype string, data io.Reader) error
}
