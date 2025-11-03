# ---------- Builder ----------
FROM golang:1.25 AS builder
WORKDIR /src

# Go-Module zuerst für Caching
COPY go.mod go.sum ./
RUN go mod download

# Rest des Repos (wegen imports aus /pkg etc.)
COPY . .

# Nur den user-service bauen
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
    -o /out/app ./services/event-service

# ---------- Runtime ----------
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /out/app /app/app

ENV GIN_MODE=release
EXPOSE 8080
CMD ["/app/app"]
