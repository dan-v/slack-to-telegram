FROM golang:1.14 AS build-env
ADD . /src
RUN cd /src && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o slack-to-telegram .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build-env /src/slack-to-telegram /
RUN addgroup -g 1000 -S slacktotelegram \
 && adduser -u 1000 -S slacktotelegram -G slacktotelegram \
 && apk add --no-cache ca-certificates \
 && rm -rf /var/cache/apk/*
USER slacktotelegram
VOLUME /config.toml
CMD ["/slack-to-telegram", "-config", "config.toml"]