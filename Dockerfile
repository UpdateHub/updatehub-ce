FROM golang:1.11-alpine3.8 as builder

RUN apk add --update git curl libarchive-dev build-base linux-headers nodejs nodejs-npm

RUN mkdir -p $$GOPATH/bin && \
    curl https://glide.sh/get | sh

ADD . /go/src/github.com/UpdateHub/updatehub-ce
WORKDIR /go/src/github.com/UpdateHub/updatehub-ce

RUN glide i && \
    go get -u github.com/gobuffalo/packr/... && \
    (cd ui; npm install && npm run build) && \
    packr install

FROM alpine:3.8

RUN apk add --update libarchive

COPY --from=builder /go/bin/updatehub-ce /usr/bin/updatehub-ce

ENTRYPOINT ["/usr/bin/updatehub-ce"]
