FROM golang:1.26-alpine AS build

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/imghat ./cmd/server

FROM alpine:latest

RUN apk add --no-cache curl ca-certificates tzdata
COPY --from=build /app/imghat /usr/local/bin/imghat
COPY --from=build /build/.env.example /app/.env

WORKDIR /app
EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:3000/healthz || exit 1

CMD ["imghat"]
