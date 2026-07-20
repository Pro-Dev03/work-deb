FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOMEMLIMIT=512MiB go build -trimpath -ldflags="-s -w" -o worktrack-api ./cmd/api

FROM alpine:3.21
WORKDIR /app
RUN addgroup -S worktrack && adduser -S -G worktrack worktrack \
    && mkdir -p /app/uploads \
    && chown -R worktrack:worktrack /app
COPY --from=builder --chown=worktrack:worktrack /app/worktrack-api ./worktrack-api
USER worktrack
EXPOSE 8080
CMD ["./worktrack-api"]
