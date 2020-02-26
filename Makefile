# Spectrum Makefile

build:
	go build ./cmd/spectrum/

test-e2e:
	go test -timeout 30m -v ./e2e
