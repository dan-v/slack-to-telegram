FROM golang:1.10

WORKDIR /go/src/github.com/dan-v/slack-to-telegram/

COPY main.go .
COPY vendor vendor

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates  
WORKDIR /root/
COPY --from=0 /go/src/github.com/dan-v/slack-to-telegram/app .
CMD ["./app", "-config", "/config.toml"]