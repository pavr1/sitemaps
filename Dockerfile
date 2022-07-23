# Start from golang base image
FROM golang:1.18-alpine as builder

# We use this to point to a specific application for building and testing
ARG USER_DATA

# install git
RUN apk update \
  && apk add --no-cache curl git \
  && apk add --no-cache ca-certificates \
  && apk add --no-cache gcc musl-dev \
  && update-ca-certificates

# Set the current working directory inside the container
WORKDIR /build

# install the golangci-lint binary
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.44.2

# Copy go.mod, go.sum files and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy pkg code, lint it, and test it
# doing this in a standalone step _should_ speed up builds
# COPY pkg/ pkg/
# RUN golangci-lint run --timeout 10m0s ./pkg/...
# RUN DATABASE_TYPE=sqlite go test -v -cover ./pkg/...

# Copy sources to the working directory
COPY . .

# lint everything in the application directory
#RUN golangci-lint run --timeout 10m0s ./$USER_DATA/...

# Test everything in the app directory
# RUN DATABASE_TYPE=sqlite go test -v -cover $USER_DATA/...

# Build the Go app
RUN echo "building app in ${USER_DATA}"
RUN go build -a -v -o server ./src/cmd/

# Start a new stage from alpine
FROM alpine:latest

RUN apk update \
  && apk upgrade \
  && apk add ca-certificates tzdata \
  && update-ca-certificates 2>/dev/null || true

WORKDIR /dist

# Copy the build artifacts from the previous stage
COPY --from=builder /build/server .

# Run the executable
CMD ["./server"]
