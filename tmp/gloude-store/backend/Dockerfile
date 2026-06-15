# Этап сборки
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git ca-certificates

# Копируем файлы модуля и загружаем зависимости
COPY go.mod ./
RUN go mod download && go mod tidy

# Копируем исходный код
COPY . .
RUN go mod tidy

# Собираем бинарный файл (статическая линковка для alpine)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/server \
    ./cmd/server

# Финальный образ (минимальный)
FROM alpine:3.19

WORKDIR /app

# Устанавливаем ca-certificates для HTTPS запросов
RUN apk add --no-cache ca-certificates wget tzdata

# Копируем бинарник из этапа сборки
COPY --from=builder /app/server .

# Создаем директорию для загрузок
RUN mkdir -p /app/uploads

# Создаем непривилегированного пользователя
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Устанавливаем владельца директорий
RUN chown -R appuser:appgroup /app

USER appuser

# Порт сервера
EXPOSE 8080

# Запускаем сервер
CMD ["./server"]
