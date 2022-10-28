package storage

import (
	"context"
	"io"
	"io/ioutil"
	"log"

	"cloud.google.com/go/storage"
)

type Config struct {
	BucketName string `env:"IMAGES_STORAGE_GOOGLE_BUCKET_NAME" envDefault:""`
}

type Storage struct {
	client *storage.Client
	bucket *storage.BucketHandle
}

func New(ctx context.Context, config Config) *Storage {
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	bucket := client.Bucket(config.BucketName)

	return &Storage{
		client: client,
		bucket: bucket,
	}
}

func (s *Storage) Upload(
	ctx context.Context, filename string, contenttype string, data io.Reader,
) error {
	wc := s.bucket.Object(filename).NewWriter(ctx)
	wc.ContentType = contenttype
	defer wc.Close()

	buf, err := ioutil.ReadAll(data)
	if err != nil {
  	log.Println(err)
		return err
	}

	wc.Write(buf)

	return nil
}
