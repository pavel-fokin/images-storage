# Images Storage

`images-storage` is a service that lets you store, retrieve, and cutout images.

This service can be deployed as the `Google Cloud Run`
and uses `Google Cloud Storage` as the database to store images and related metadata.

## Requirements

- Go 1.19
- Google Cloud Storage

## Run

`Makefile` can be used to run a few simple commands.

```sh
$ make help
Makefile help:

run               Run service.
tests             Run unittests.
swagger           Generate swagger documentation.
                  Requires installed https://github.com/swaggo/swag
```

### Google Cloud Storage

To run locally it will be required to create a `Google Cloud Storage` and 
add related credentials that have permissions to `Read/Write` to this bucket.
To set up required permissions it will be required to add `GOOGLE_APPLICATION_CREDENTIALS` environment variable.

### Environment Variables

Service supports configuration with environment variables either directly, or with `*.env` file.

```
GOOGLE_APPLICATION_CREDENTIALS=<credentials-to-run-locally>

IMAGES_STORAGE_ENV_FILE=<path-to-the-env-file>
IMAGES_STORAGE_GOOGLE_BUCKET_NAME=<bucket-name>
```

## Deployment

Service uses `Google Cloud` products to manage deployments and infrastructure.

- `Cloud Build` - building containers and CI/CD 
- `Artifacts Registry` - repo with Docker images

- `Cloud Run` runs `images-storage` service
- `Secret Manager` stores environment variables

- `Cloud Storage` - database for images and its metadata
