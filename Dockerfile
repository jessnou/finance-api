# Используем официальный образ Go для сборки
FROM golang:1.23.4 AS builder

# Устанавливаем рабочую директорию для сборки
WORKDIR /app

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Устанавливаем Goose в /usr/local/bin (это стандартный путь для утилит)
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Копируем исходный код в контейнер
COPY . .

# Собираем бинарник
RUN go build -o finance-api ./cmd/main.go

# Используем более легкий образ для запуска
FROM ubuntu:22.04

WORKDIR /app

# Устанавливаем необходимые пакеты, включая libc6 для glibc
RUN apt-get update && apt-get install -y \
    libc6-dev \
    libstdc++6 \
    && apt-get clean

# Копируем скомпилированный бинарник из билдера
COPY --from=builder /app/finance-api .

# Копируем директорию миграций
COPY --from=builder /app/migrations ./migrations

# Копируем бинарник `goose` из билдера
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Устанавливаем права на выполнение
RUN chmod +x finance-api

# Добавляем стандартные пути для бинарных файлов Go в PATH
ENV PATH="/usr/local/go/bin:/usr/local/bin:${PATH}"

# Проверка установки goose
RUN echo $PATH && ls /usr/local/bin/

# Запускаем сервер
CMD ["./finance-api"]
