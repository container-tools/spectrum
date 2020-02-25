
REGISTRY := quay.io
IMAGE_NAME := container-tools/spectrum

build:
	go build ./cmd/spectrum/

image:
	docker build -t $(REGISTRY)$(IMAGE_NAME) .

release: image
	docker push $(REGISTRY)$(IMAGE_NAME)

test-e2e:
	go test -timeout 30m -v ./e2e
