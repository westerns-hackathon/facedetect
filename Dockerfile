# Используем образ golang с базой Alpine
FROM golang:1.23-alpine

# Устанавливаем необходимые пакеты, включая make, gcc, и другие зависимости для компиляции
RUN apk update && apk add --no-cache \
    make \
    gcc \
    g++ \
    libc-dev \
    bash \
    git \
    sudo \
    && rm -rf /var/cache/apk/*

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем все файлы проекта в контейнер
COPY . .

# Проверяем, что Makefile скопирован в контейнер
RUN ls -l /app

# Загружаем зависимости Go
RUN go mod tidy

# Устанавливаем зависимости с помощью make install
# Используем команду yes для автоматического подтверждения всех запросов
RUN yes | make install

# Стартуем приложение с помощью make run
CMD ["make", "run"]
