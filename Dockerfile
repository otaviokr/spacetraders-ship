FROM golang:alpine as builder

WORKDIR $GOPATH/src/github.com/otaviokr/spacetraders-ship/
COPY component/ component/
COPY kafka/ kafka/
COPY web/ web/
COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go

RUN apk --no-cache add ca-certificates && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/spacetraders-ship .

FROM scratch

WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/spacetraders-ship /app/

ENTRYPOINT ["./spacetraders-ship"]
