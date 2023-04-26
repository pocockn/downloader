FROM golang:1.20.3-alpine3.17 AS builder

RUN apk update && apk add make git libc-dev

ADD . /src/

WORKDIR /src/

RUN go build -o downloader

FROM alpine:3.17

COPY --from=builder /src/ /src/

WORKDIR /src

ENTRYPOINT ["/src/downloader"]