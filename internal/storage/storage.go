package storage

import (
	"context"
	"io"
)

type StorageUploader interface {
	Upload(ctx context.Context, filename string, contenttype string, data io.Reader) error
}
