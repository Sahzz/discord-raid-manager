FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/app/
COPY . .

RUN go get -d -v
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/app

FROM alpine
COPY --from=builder /go/bin/app /go/bin/app
RUN apk add ca-certificates
ENTRYPOINT ["/go/bin/app"]
