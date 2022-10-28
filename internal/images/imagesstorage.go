package images

import (
	"context"
	"io"
	"log"
	// "image"
	// _ "image/jpeg"

	"github.com/google/uuid"
)

type ImageRaw io.Reader

type ImagesStorage struct {
	storage ListUploader
}

func New(storage ListUploader) *ImagesStorage {
	return &ImagesStorage{
		storage: storage,
	}
}

func (i *ImagesStorage) Add(data io.Reader, contenttype string) error {
	// m, _, err := image.Decode(data)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(m.Bounds())
	filename := uuid.New().String()

	err := i.storage.Upload(context.Background(), filename, contenttype, data)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (i *ImagesStorage) List(ctx context.Context) error {
	i.storage.List(ctx)
	return nil
}
