FROM golang:1.19-alpine

RUN go env -w GO111MODULE="on"
RUN apk add git
RUN apk add build-base
RUN apk add --no-cache ca-certificates && update-ca-certificates
# RUN apk add --no-cache ffmpeg
COPY . /go/src
WORKDIR /go/src
CMD ["go","run","./cmd/server/main.go"]