FROM golang:1.12 AS builder

# Dependencies
RUN apt-get update \
    && apt-get install -y bzr \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir /build
WORKDIR /build

RUN go get github.com/aws/aws-sdk-go \
    && go get github.com/prometheus/client_golang/prometheus \
    && go get github.com/Sirupsen/logrus

# Build
COPY ./*.go /build/
COPY controllers /build/controllers

RUN go build ./main.go
RUN cp ./main /bin/s3bucket_exporter

FROM debian:buster-slim
COPY --from=builder /bin/s3bucket_exporter /bin/s3bucket_exporter
WORKDIR /tmp
ENTRYPOINT ["/bin/s3bucket_exporter"]
