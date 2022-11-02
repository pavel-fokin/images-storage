.PHONY: all
all: help

.PHONY: help
help:
	@echo 'Makefile help:                                                         '
	@echo '                                                                       '
	@echo 'run               Run service                                       '
	@echo 'tests             Run unittests                                     '
	@echo 'swagger           Generate swagger documentation.                   '
	@echo '                  Requires installed https://github.com/swaggo/swag '

.PHONY: tests
tests:
	@go test -coverprofile=coverage.out ./... -count=1

.PHONY: swagger
swagger:
	swag init -d cmd/images-storage,internal/server

.PHONY: run
run:
	@go run cmd/images-storage/main.go