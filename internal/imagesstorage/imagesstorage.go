package imagesstorage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"
)

var (
	ErrImageNotExist = errors.New("Image doesn't exist")
)

// ImageStorage keeps the logic for images-storage functionality.
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
	image, err := NewImage(uuid.New(), contenttype, data)
	if err != nil {
		return Image{}, fmt.Errorf("Add(): %w", err)
	}

	image, err = i.storage.Upload(ctx, image)
	if err != nil {
		return Image{}, fmt.Errorf("Add(): %w", err)
	}

	return image, nil
}

func (i *ImagesStorage) Update(
	ctx context.Context, id string, data io.Reader, contenttype string,
) (Image, error) {
	if doesExist := i.storage.DoesExist(ctx, id); !doesExist {
		return Image{}, ErrImageNotExist
	}

	image, err := NewImage(uuid.MustParse(id), contenttype, data)
	if err != nil {
		return Image{}, fmt.Errorf("Update(): %w", err)
	}

	image, err = i.storage.Upload(ctx, image)
	if err != nil {
		return Image{}, fmt.Errorf("Update(): %w", err)
	}

	return image, nil
}

func (i *ImagesStorage) Metadata(
	ctx context.Context, uuid string,
) (Image, error) {
	image, err := i.storage.Metadata(ctx, uuid)
	if err != nil {
		return Image{}, fmt.Errorf("Metadata(): %w", err)
	}

	return image, nil
}

func (i *ImagesStorage) Data(
	ctx context.Context, uuid string, bbox BBox,
) (io.Reader, string, error) {
	data, contenttype, err := i.storage.Download(ctx, uuid)
	if err != nil {
		return bytes.NewReader([]byte{}), "", fmt.Errorf("Data(): %w", err)
	}

	data, err = CutOut(data, bbox)
	if err != nil {
		return bytes.NewReader([]byte{}), "", fmt.Errorf("Data(): %w", err)
	}

	return data, contenttype, nil
}

func (i *ImagesStorage) List(ctx context.Context) ([]Image, error) {
	images, err := i.storage.List(ctx)
	if err != nil {
		return []Image{}, fmt.Errorf("List(): %w", err)
	}
	return images, nil
}
