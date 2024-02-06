FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install deps
RUN apk add --update gcc musl-dev
COPY go.mod go.sum ./
RUN go mod download

# Build app
RUN --mount=target=. \
  --mount=type=cache,target=/root/.cache/go-build \
  go build -v -tags "linux" -o /server

# Runner
FROM alpine:latest

COPY --from=builder /server /app/server
EXPOSE 8080

CMD [ "/app/server" ]
