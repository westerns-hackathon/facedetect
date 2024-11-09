<div align="center">
  <h1>Face Detect</h1>
  <img src="https://upload.wikimedia.org/wikipedia/commons/a/a5/Go_Logo.svg" alt="Go" width="50" height="50"/>
  <img src="https://upload.wikimedia.org/wikipedia/commons/e/e3/Docker_logo.svg" alt="Docker" width="50" height="50"/>
  <img src="https://upload.wikimedia.org/wikipedia/commons/e/e5/YAML_logo.svg" alt="YAML" width="50" height="50"/>
  <img src="https://raw.githubusercontent.com/jmoiron/sqlx/master/docs/images/sqlx-logo.png" alt="sqlx" width="50" height="50"/>
  <img src="https://upload.wikimedia.org/wikipedia/commons/2/2a/Pqsql-logo.svg" alt="pq" width="50" height="50"/>
  <img src="https://raw.githubusercontent.com/jeffallen/cleanenv/master/docs/assets/cleanenv-logo.png" alt="cleanenv" width="50" height="50"/>
  <img src="https://gocv.io/img/gocv.svg" alt="gocv" width="50" height="50"/>
  <img src="https://raw.githubusercontent.com/caspervonb/face-recognition-go/master/docs/images/logo.png" alt="go-face" width="50" height="50"/>
  <img src="https://go-recognizer.com/assets/img/recognizer-logo.png" alt="go-recognizer" width="50" height="50"/>
</div>

## Описание

**Face Detect** — проект, созданный в рамках хакатона. Задача проекта — разработать систему распознавания лиц, которая находит в базе данных картину, схожую с изображением, загруженным пользователем. Это приложение использует несколько библиотек для обработки изображений и работы с базой данных, а также интеграцию с контейнерами через Docker.

## Стек технологий

- **Go** — основной язык разработки.
- **Docker** — для контейнеризации приложения.
- **Makefile** — для автоматизации команд и настройки.
- **YAML** — для конфигурации приложения.
- **SQLx** — для работы с базой данных.
- **pq** — драйвер для работы с PostgreSQL.
- **Cleanenv** — для управления конфигурациями в проекте.
- **GoCV** — для работы с OpenCV и обработки изображений.
- **Go-Face** — для распознавания лиц.
- **Go-Recognizer** — для дополнительной обработки лиц и их сравнения.

## Установка и запуск

Для запуска проекта выполните следующие шаги:

1. Клонируйте репозиторий:
   ```bash
   git clone git@github.com:westerns-hackathon/facedetect.git
   cd face-detect
   ```
2. Запустите контейнеры с помощью Docker Compose:
    ```bash
    docker compose up -d
    ```
3. Запустите сервер приложения:
    ```bash
    go run ./cmd/
    ```
---

**© 2024 Face Detect Project**