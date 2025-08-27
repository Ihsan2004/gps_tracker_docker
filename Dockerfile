# Adım 1: Temel Dockerfile - Sadece Go uygulaması
FROM golang:1.24-alpine AS builder

# Gerekli araçları yükle
RUN apk add --no-cache git gcc musl-dev

# Çalışma dizini
WORKDIR /app

# Go mod dosyalarını kopyala (caching için)
COPY go.mod go.sum ./

# Bağımlılıkları indir
RUN go mod download

# Tüm kaynak kodunu kopyala
COPY . .
# Uygulamayı build et
RUN CGO_ENABLED=1 GOOS=linux go build -o gps_tracker_docker ./main.go

# Production stage
FROM alpine:latest

# Gerekli paketleri yükle
RUN apk --no-cache add ca-certificates tzdata curl

# Çalışma dizini
WORKDIR /app

# Binary'yi kopyala
COPY --from=builder /app/gps_tracker_docker .
# Docs klasörünü kopyala (Swagger için)
COPY --from=builder /app/docs ./docs

# Port'u expose et
EXPOSE 3000

# Basit health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:3000/ || exit 1

# Uygulamayı çalıştır
CMD ["./gps_tracker_docker"]