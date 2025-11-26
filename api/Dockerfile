FROM golang:1.25.4 AS builder

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o build/api

FROM debian:stable-slim

COPY --from=builder /app/build/api /api

CMD [ "/api" ]
