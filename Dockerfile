FROM registry.access.redhat.com/ubi8/ubi:8.1
ADD ./spectrum /app/
WORKDIR /app
ENTRYPOINT ["./spectrum"]
