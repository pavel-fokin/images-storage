package images

import (
	"context"
	"io"
	"log"
	// "image"
	// _ "image/jpeg"

	"github.com/google/uuid"

	"pavel-fokin/images-storage/internal/storage"
)

type ImageRaw io.Reader

type Images struct {
	storage storage.StorageUploader
}

func New(storage storage.StorageUploader) *Images {
	return &Images{
		storage: storage,
	}
}

func (i *Images) Add(data io.Reader, contenttype string) error {
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

func (i *Images) List() error {
	return nil
}
