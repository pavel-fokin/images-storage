package imagesstorage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"

	"github.com/google/uuid"
)

var (
	ErrImageNotExist = errors.New("Image doesn't exist")
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

	uuid := uuid.New().String()

	image, err := i.storage.Upload(ctx, uuid, contenttype, data)
	if err != nil {
		log.Fatal(err)
	}

	return image, nil
}

func (i *ImagesStorage) Update(
	ctx context.Context, uuid string, data io.Reader, contenttype string,
) (Image, error) {
	if doesExist := i.storage.DoesExist(ctx, uuid); !doesExist {
		return Image{}, ErrImageNotExist
	}

	image, err := i.storage.Upload(ctx, uuid, contenttype, data)
	if err != nil {
		return Image{}, err
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

func (i *ImagesStorage) Data(
	ctx context.Context, uuid string,
) (io.Reader, string, error) {
	data, contenttype, err := i.storage.Download(ctx, uuid)
	if err != nil {
		log.Println(err)
		return bytes.NewReader([]byte{}), "", err
	}

	return data, contenttype, nil
}

func (i *ImagesStorage) List(ctx context.Context) ([]Image, error) {
	images, err := i.storage.List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return images, nil
}
