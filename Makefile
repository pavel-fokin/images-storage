.PHONY: tests
tests:
	@go test -coverprofile=coverage.out ./... -count=1

.PHONY: run
run:
	@go run cmd/images-storage/main.go
