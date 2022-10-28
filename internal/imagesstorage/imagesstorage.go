package imagesstorage

import (
	"context"
	"io"
	"log"

	// "image"
	// _ "image/jpeg"

	"github.com/google/uuid"
)

type Image struct {
	Name        string
	ContentType string
	Width       int
	Height      int
	Size        int
	UploadedAt  string
}

type ImagesStorage struct {
	storage Storage
}

func New(storage Storage) *ImagesStorage {
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

func (i *ImagesStorage) List(ctx context.Context) ([]Image, error) {
	images, err := i.storage.List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return images, nil
}
