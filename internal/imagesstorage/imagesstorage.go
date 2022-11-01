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

	buf, err := io.ReadAll(data)
	if err != nil {
		return Image{}, fmt.Errorf("Add(): %w", err)
	}

	width, height, err := GetWidthHeight(bytes.NewReader(buf))
	if err != nil {
		return Image{}, fmt.Errorf("Add(): %w", err)
	}
	metadata := map[string]string{
		"ImageWidth":  fmt.Sprint(width),
		"ImageHeight": fmt.Sprint(height),
	}

	image, err := i.storage.Upload(ctx, uuid, contenttype, bytes.NewReader(buf), metadata)
	if err != nil {
		return Image{}, fmt.Errorf("Add(): %w", err)
	}

	return image, nil
}

func (i *ImagesStorage) Update(
	ctx context.Context, uuid string, data io.Reader, contenttype string,
) (Image, error) {
	if doesExist := i.storage.DoesExist(ctx, uuid); !doesExist {
		return Image{}, ErrImageNotExist
	}

	buf, err := io.ReadAll(data)
	if err != nil {
		return Image{}, fmt.Errorf("Updates(): %w", err)
	}

	width, height, err := GetWidthHeight(bytes.NewReader(buf))
	if err != nil {
		return Image{}, fmt.Errorf("Update(): %w", err)
	}
	metadata := map[string]string{
		"ImageWidth":  fmt.Sprint(width),
		"ImageHeight": fmt.Sprint(height),
	}

	image, err := i.storage.Upload(ctx, uuid, contenttype, bytes.NewReader(buf), metadata)
	if err != nil {
		return Image{}, fmt.Errorf("Updates(): %w", err)
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

	if bbox.Valid() {
		data, err = CutOut(data, bbox)
		if err != nil {
			return bytes.NewReader([]byte{}), "", fmt.Errorf("Data(): %w", err)
		}
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
