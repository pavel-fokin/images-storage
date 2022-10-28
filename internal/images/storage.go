package images

import (
	"context"
	"io"
)

type ListUploader interface {
	Lister
	Uploader
}

type Lister interface {
	List(ctx context.Context) error
}

type Uploader interface {
	Upload(ctx context.Context, filename string, contenttype string, data io.Reader) error
}
