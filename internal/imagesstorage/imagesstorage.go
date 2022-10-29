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

func (i *ImagesStorage) Add(
	ctx context.Context, data io.Reader, contenttype string,
) (Image, error) {
	// m, _, err := image.Decode(data)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(m.Bounds())
	uuid := uuid.New().String()

	image, err := i.storage.Upload(ctx, uuid, contenttype, data)
	if err != nil {
		log.Fatal(err)
	}

	return image, nil
}

func (i *ImagesStorage) Metadata(
	ctx context.Context, uuid string,
) (Image, error) {
	image, err := i.storage.Metadata(ctx, uuid)
	if err != nil {
		log.Fatal(err)
	}

	return image, nil
}

func (i *ImagesStorage) List(ctx context.Context) ([]Image, error) {
	images, err := i.storage.List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return images, nil
}
