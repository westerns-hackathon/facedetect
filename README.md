<a name="top-of-readme"></a>
<div align="center">

<h1 align="center">Face Detect</h1>

[![PostgreSQL][PostgreSQL]][PostgreSQL-url]
[![Go][Go]][Go-url]
[![Ubuntu][Ubuntu]][Ubuntu-url]
[![Docker][Docker]][Docker-url]
[![Git][Git]][Git-url]
[![TensorFlow][TensorFlow]][TensorFlow-url]
[![Makefile][Makefile]][Makefile-url]
[![YAML][YAML]][YAML-url]

<p align="center">
   Face Detect — проект, созданный в рамках хакатона.
Задача проекта — разработать систему распознавания лиц, которая находит в базе данных картину, схожую с изображением,
загруженным пользователем. Это приложение использует несколько библиотек для обработки изображений и работы с базой данных,
а также интеграцию с контейнерами через Docker.
</p>

</div>

Метод | Endpoint            | Смысл                                                                               | 
--- |---------------------|-------------------------------------------------------------------------------------| 
POST | /v1/app/photo       | Клиент делает запрос с фото, сервер сохраняет фото и выдает распознанные лица на нем | 
POST | /v1/app/photo/match | Клиент делает запрос с двумя фото, сервер выдает коефицент сходства (от 0 или > 0)  |
POST | /v1/app/photo/find | Клиент делает запрос с фото, сервер выдает фотографии, схожие с лицами на фото      


## Установка и запуск

Для запуска проекта выполните следующие шаги:

> [!IMPORTANT]
>
> Необходимо установить все зависимости (goCV, dlib).

1. Клонируйте репозиторий:
   ```bash
   git clone git@github.com:westerns-hackathon/facedetect.git
   cd face-detect
   make install 
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

[PostgreSQL]: https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white
[PostgreSQL-url]: https://www.postgresql.org/
[Ubuntu]: https://img.shields.io/badge/Ubuntu-E95420?style=for-the-badge&logo=ubuntu&logoColor=white
[Ubuntu-url]: https://ubuntu.com/
[Go]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[Go-url]: https://go.dev/
[TensorFlow]: https://img.shields.io/badge/TensorFlow-FF6F00?style=for-the-badge&logo=tensorflow&logoColor=white
[TensorFlow-url]: https://www.tensorflow.org/
[Git]: https://img.shields.io/badge/GIT-E44C30?style=for-the-badge&logo=git&logoColor=white
[Git-url]: https://git-scm.com/
[Docker]: https://img.shields.io/badge/Docker-0c49c2?style=for-the-badge&logo=docker&logoColor=white
[Docker-url]: https://www.docker.com/
[Makefile]: https://img.shields.io/badge/Makefile-eaeaea?style=for-the-badge&logo=makefile&logoColor=white
[Makefile-url]: https://www.gnu.org/software/make/manual/make.html
[YAML]: https://img.shields.io/badge/yaml-red?style=for-the-badge&logo=yaml&logoColor=white
[YAML-url]: https://yaml.org/
[stats]: https://starchart.cc/{AnwarSaginbai}/{https://github.com/westerns-hackathon/facedetect}.svg