FROM golang:1.20.1-alpine AS build

WORKDIR /go/src/build

COPY go.mod go.sum ./
RUN go mod download

COPY src ./src

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o app_binary src/main.go

# hadolint ignore=DL3006
FROM gcr.io/distroless/base AS app

WORKDIR /app

COPY --from=build /go/src/build/app_binary .

CMD ["./app_binary"]
