FROM golang:1.15-alpine3.12 AS builder
RUN apk add --update git curl libarchive-dev build-base linux-headers nodejs nodejs-npm

WORKDIR /app/server

# Copy go mod dependency information and download them so we can
# cache it for use when dependencies do not change.
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy rest of code for build.
COPY . .

RUN go get -u github.com/gobuffalo/packr/... && \
    (cd ui; npm install && npm run build) && \
    packr install

FROM alpine:3.12

RUN apk add --update libarchive

COPY --from=builder /go/bin/updatehub-ce /usr/bin/updatehub-ce

ENTRYPOINT ["/usr/bin/updatehub-ce"]
