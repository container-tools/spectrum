
build:
	go build ./cmd/spectrum/

image:
	docker build -t nferraro/spectrum .

release: image
	docker push nferraro/spectrum
