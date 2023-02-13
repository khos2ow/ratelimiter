FROM alpine:3.17.2

RUN apk --no-cache add ca-certificates

COPY ratelimiter /usr/local/bin/ratelimiter

VOLUME /ssl

EXPOSE 8080 8443

ENTRYPOINT ["ratelimiter"]
