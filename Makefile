.PHONY: tests
tests:
	@go test -coverprofile=coverage.out ./... -count=1

.PHONY: swagger
swagger:
	swag init -d cmd/images-storage,internal/server


.PHONY: run
run:
	@go run cmd/images-storage/main.go