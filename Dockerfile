# refer https://github.com/nats-io/nats-replicator/blob/master/Dockerfile
ARG GO_VERSION=1.12.9
ARG UBUNTU_VERSION=19.04
# compile stage
FROM golang:${GO_VERSION} as immediate

WORKDIR /src/nats-bench
COPY . .

# binary build
RUN go mod download && \
    CGO_ENABLED=0 go build -v -a -tags netgo -installsuffix netgo -o /nats-bench

# final docker image building stage
FROM ubuntu:{UBUNTU_VERSION} as builder

RUN mkdir -p /nats/bin && mkdir /nats/conf
COPY --from=immediate /nats-bench /nats/bin/nats-bench
RUN ln -ns /nats/bin/nats-bench /bin/nats-bench

ENTRYPOINT ["/bin/nats-bench"]
CMD ["--help"]
