FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o spectrum ./cmd/spectrum/
FROM registry.access.redhat.com/ubi8/ubi:8.1
COPY --from=builder /build/spectrum /app/
WORKDIR /app
ENTRYPOINT ["./spectrum"]
