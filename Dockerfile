FROM golang:1.12

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

# Cleanup
WORKDIR /
RUN rm -Rf /build

ENTRYPOINT ["/bin/s3bucket_exporter"]
