# ---- Build stage ----
FROM golang:1.24-alpine AS builder

# Bağımlılıklar
RUN apk add --no-cache git

# RAM azaldığında sorun yaşamamak için derlemeyi kısıtla
ENV GOFLAGS="-p=1" GOMAXPROCS=1
ENV GOPROXY=https://proxy.golang.org,direct

WORKDIR /src

# Modülleri önceden indir (cache için)
COPY go.mod go.sum ./
RUN go mod download

# Kaynakları kopyala
COPY . .

# main.go kökte
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o app .

# ---- Runtime stage ----
FROM alpine:latest
WORKDIR /app
COPY --from=builder /src/app /app/app

ENV APP_PORT=3000
ENV GIN_MODE=release

EXPOSE 3000
CMD ["/app/app"]