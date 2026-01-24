# API Каталога Книг
> **Production-ready backend на Go для управления каталогом книг**  
> Построен по принципам чистой архитектуры, с PostgreSQL и Docker-first подходом.  
> Разработан как демонстрация практик middle-уровня в backend-разработке.

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-336791?logo=postgresql)](https://www.postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker)](https://www.docker.com)

## Ключевые достижения

- **Разработал** многослойную архитектуру (`api` → `service` → `repository`) по принципам Clean Architecture
- **Внедрил** production-ready миграции БД с поддержкой zero-downtime деплоя
- **Спроектировал** гибкую систему авторизации: публичные эндпоинты для авторов/жанров, защищённые — для книг
- **Обеспечил** полное покрытие тестами: unit-тесты с фейковым репозиторием + интеграционные тесты с тестовой БД + E2E-тесты разделенные по логике http запросов books/author/genre
- **Оптимизировал** работу с БД через connection pooling (`pgx`), таймауты контекста и индекс на частые SELECT запросы.
- **Устранил** риски при остановке сервиса: реализован graceful shutdown (обработка SIGTERM/SIGINT)
- **Гарантировал** целостность данных: внешние ключи, check-ограничения в PostgreSQL

## Используемые технологии

| Уровень        | Технология                     |
|----------------|--------------------------------|
| Язык           | Go 1.24+                       |
| Веб-фреймворк  | `net/http` + Gorilla Mux       |
| База данных    | PostgreSQL 17                  |
| Миграции       | `golang-migrate`               |
| Тестирование   | `testify`, интеграция с БД     |
| Инфраструктура | Docker, Docker Compose         |
| Логирование    | Структурированное (`slog`)     |


##  Структура проекта
``` bash
.
├── cmd/
│ ├── main/ # Основной HTTP-сервер
│ └── migrate/ # Утилита для миграций
├── migrations/ # SQL-миграции
├── pkg/
│ ├── api/ # HTTP-хендлеры, middleware, E2E-тесты
│ ├── service/ # Бизнес-логика, unit-тесты
│ ├── repository/ # Работа с данными, интеграционные тесты
│ └── models/ # DTO и доменные модели
├── .env # Переменные окружения (пример)
├── .env.example # Шаблон для новых разработчиков
├── docker-compose.yml
└── Dockerfile
```

## Быстрый запуск

### Требования
- Go 1.24+
- Docker + Docker Compose

### Запуск в режиме разработки
```bash
# 1. Запустить PostgreSQL в Docker
docker-compose up -d postgres

# 2. Применить миграции
go run ./cmd/migrate up

# 3. Запустить сервис
go run ./cmd/main

```
```bash
# Можно отменить миграции
go run ./cmd/migrate down
```
### Полная изоляция: миграции ≠ сервиc
Сервис никогда не трогает схему БД — он только читает/пишет данные.

### Запуск в Docker

```bash
# Сборка и запуск всего стека с автоматическим применением миграций
docker-compose up --build
```

### Публичные эндпоинты (не требуют токена):

| Метод | Путь                        | Описание                     |
|-------|-----------------------------|------------------------------|
| GET   | `/api/books`                | Список всех книг             |
| GET   | `/api/books/withauthors`    | Список книг с авторами       |
| GET   | `/api/authors`              | Список авторов               |
| POST  | `/api/authors`              | Добавление нового автора     |
| GET   | `/api/genres`               | Список жанров                |
| POST  | `/api/genres`               | Добавление нового жанра      |

###  Приватные эндпоинты (требуют токена)

| Метод | Путь                        | Описание                     |
|-------|-----------------------------|------------------------------|
| POST  | `/api/books`                | Создание новой книги         |
| PATCH | `/api/books?id={id}`        | Частичное обновление книги   |
| DELETE| `/api/books?id={id}`        | Удаление книги               |



>  Все приватные эндпоинты требуют заголовок:  
> `Authorization: <ваш_токен>` (см. раздел [Авторизация](#авторизация-authorization))


## Авторизация (Authorization)

Авторизация через заголовок Authorization: <токен> (настраивается через переменную окружения AUTH_TOKEN) *adminToken по умолчанию*.
### Примеры запросов с авторизацией
``` bash
curl -X PATCH "http://localhost:8080/api/books?id=5" \
  -H "Authorization: adminToken" \
  -H "Content-Type: application/json" \
  -d '{"price": 599}'
```
``` bash
curl -X DELETE "http://localhost:8080/api/books?id=5" \
  -H "Authorization: adminToken"
```


## Стратегия тестирования
* Unit-тесты: изолированная проверка бизнес-логики с использованием фейкового репозитория
* Integration-тесты 
  * Проверка работы репозитория с реальной PostgreSQL
  * Автоматическое применение и откат миграций перед каждым тестом
  * Валидация ограничений БД (уникальность, check-ограничения)
  * Полная изоляция состояния через TRUNCATE ... RESTART IDENTITY
* E2E-тесты: сквозная проверка HTTP-слоя с реальной БД и middleware
Покрытие: более 75% по операторам во всех пакетах

### Запуск всех тестов
``` bash 
# Только unit-тесты (без БД)
go test ./...

# Только интеграционные тесты
go test -v -tags=integration ./pkg/repository/postgres/...

# Только E2E-тесты
go test -v -tags=e2e ./pkg/api/...

# Все тесты
go test -v -tags="integration e2e" ./...
```
ИЛИ
``` bash 
# Все тесты
./goTests.sh
```

### Документация API
Доступна по адресу: http://localhost:8080/swagger/index.html

