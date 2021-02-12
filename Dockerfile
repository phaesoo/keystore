
# build: docker build --force-rm -t phaesoo/shield .

### 1.Dependencies
FROM golang:1.15-alpine AS dependencies
LABEL maintainer "phaesoo <phaesoo@gmail.com>"
WORKDIR /app
# Set build flags
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
# Create appuser
RUN adduser -D -g '' appuser
# Install required binaries
RUN apk add --update --no-cache zip git make ca-certificates tzdata && update-ca-certificates
# Prepare timezone data for scratch image
# (required if using any other timezone than the system default)
RUN cd /usr/share/zoneinfo && zip -r -0 /zoneinfo.zip .
# Copy app sources
COPY go.mod go.mod
COPY go.sum go.sum
COPY Makefile Makefile
# Download all golang package dependencies
RUN make deps


### 2.Build
FROM dependencies AS build
# Copy source files
COPY . .
# Build executable
RUN make build-docker

### 3. Release
FROM alpine:latest AS release
WORKDIR /app
# Expose application port
# Import the user and group files to run the app as an unpriviledged user
COPY --from=build /etc/passwd /etc/passwd
# CA certificates are required to call HTTPS endpoints
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Timezone info
# Golang knows how to load this https://golang.org/pkg/time/#LoadLocation from the ZONEINFO var
ENV ZONEINFO /zoneinfo.zip
COPY --from=build /zoneinfo.zip /
# Config File
COPY --from=build /app/config.yaml /app/config.yaml
# Use an unprivileged user
# USER appuser
# Grab compiled binary from build
COPY --from=build /app/main /app/main

# Set entry point
CMD [ "./main" ]