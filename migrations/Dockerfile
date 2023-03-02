FROM golang:1.20.1-alpine AS builder

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM scratch

WORKDIR /

COPY --from=builder /go/bin/goose /goose
COPY *.sql .

CMD ["./goose", "up"]
