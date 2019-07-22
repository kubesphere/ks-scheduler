FROM golang:1.10-alpine as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV VERSION=2.1.0

# build
WORKDIR /root/
COPY . .
RUN go build -mod=vendor -a -o ks-schduler ./cmd

# runtime image
FROM gcr.io/google_containers/ubuntu-slim:0.14
COPY --from=builder /root/ks-schduler .
ENTRYPOINT ["./ks-scheduler"]