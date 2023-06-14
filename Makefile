# Spectrum Makefile

build:
	@echo "####### Running unit test..."
	go test ./pkg/...
	@echo "####### Building CLI..."
	go build ./cmd/spectrum/

test-e2e:
	go test -timeout 30m -v ./e2e
