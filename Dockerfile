FROM golang:1.14.2-alpine AS builder

RUN apk add --update --no-cache ca-certificates bash make gcc musl-dev git openssh wget curl bzr

RUN mkdir /build
WORKDIR /build

COPY . /build

RUN make build

################

FROM alpine:3.11.6

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app

COPY --from=builder /build/bin/linux-amd64/ratelimiter /app/

VOLUME /app/ssl

EXPOSE 8080 8443

CMD ["/app/ratelimiter"]
