FROM golang:1.17.3-alpine AS builder

RUN apk add --update --no-cache ca-certificates bash make gcc musl-dev git openssh wget curl

WORKDIR /go/src/ratelimiter

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make build

################

FROM alpine:3.14.2

RUN apk --no-cache add ca-certificates

COPY --from=builder /go/src/ratelimiter/bin/linux-amd64/ratelimiter /usr/local/bin/

VOLUME /ssl

EXPOSE 8080 8443

ENTRYPOINT ["ratelimiter"]
