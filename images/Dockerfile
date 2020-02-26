FROM golang:1.14.0-alpine AS builder

RUN apk add --update --no-cache ca-certificates bash make gcc musl-dev git openssh wget curl bzr

RUN mkdir /build
WORKDIR /build

COPY . /build

RUN make build

################

FROM alpine:3.11.3

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app

COPY --from=builder /build/bin/linux-amd64/ratelimiter /app/

VOLUME /app/ssl

EXPOSE 8000 8443

ENV RATE_LIMIT $RATE_LIMIT
ENV RATE_INTERVAL $RATE_INTERVAL
ENV RATE_TIMEUNIT $RATE_TIMEUNIT

ENV USE_REDIS $USE_REDIS
ENV REDIS_URL $REDIS_URL
ENV REDIS_PORT $REDIS_PORT
ENV REDIS_PASSWORD $REDIS_PASSWORD
ENV BACKEND_SERVER $BACKEND_SERVER

CMD ["/app/ratelimiter"]
