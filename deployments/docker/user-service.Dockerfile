FROM golang:1.25.6-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-w -s" -trimpath -o main ./cmd/user-api


FROM alpine:3.23.2

COPY --from=builder /app/main /bin/main

ENTRYPOINT [ "/bin/main" ]
