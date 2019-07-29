FROM golang:1.12-alpine as builder

ENV CGO_ENABLED=1
ENV VERSION=2.1.0
ENV GO111MODULE=auto
ENV GOMOD=/root/go.mod

# build
WORKDIR /root/
COPY . .
RUN apk add --no-cache gcc musl-dev
RUN GOOS=linux GOARCH=amd64 go build -mod=vendor -a -o ks-pipeline-scheduler /root/cmd

# runtime image
FROM alpine:latest

ENV DATA_PATH=/data/scheduler.db

RUN mkdir /data

COPY --from=builder /root/ks-pipeline-scheduler .

CMD ["./ks-pipeline-scheduler", "--logtostderr=true", "--v=6"]