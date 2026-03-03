# Docs Aggregation Service

Сервис агрегации документов, написанный на Go, работающий с базой данных документов MongoDB компании"Астрал-Софт". Проект разработан с применением чистой архитектуры.

## Основные возможности и особенности

- Агрегация документов коллекции по диапазону дат и полю `fiscalDriveNumber`
- Асинхронный бэкенд сервиса на основе менеджера задач
- UI на основе REST API по HTTP
- Сбор метрик веб-сервера и менеджера задач через Prometheus

## Структура репозитория

```
docs-aggregation-service/
├── cmd/
|   └── main.go     # package main и точка входа в сервис
├── internal/
|   ├── adapters/
|   |   └── ...     # Слой адаптеров
|   ├── app/
|   |   └── ...     # Слой приложения
|   ├── domains/
|   |   └── ...     # Домены с общими структурами и переменными
|   ├── metrics/
|   |   └── ...     # Сбор метрик через Prometheus
|   ├── ui/
|   |   └── ...     # Слой UI (http-сервер)
|   └── usecases/
|       └── ...     # Слой юзкейсов
├── .gitignore
├── go.mod
├── go.sum
├── LISENCE
└── README.md
```

## Быстрый старт

1. Клонируйте репозиторий:
```sh
git clone https://github.com/Rurikkkk/docs-aggregation-service
```
2. Перейдите в корневую папку проекта:
```sh
cd docs-aggregation-service
```
3. Сконфигурируйте переменные окружения в файле .env:
- `MONGO_URL` - URL сервера MongoDB (по умолчанию: `mongodb://localhost:27017`)
- `MONGO_DB_NAME` - имя базы данных MongoDB (по умолчанию: `docs-aggregation-service`)
- `DOCS_COLLECTION_NAME` - имя коллекции документов (по умолчанию: `docs`)
- `TASKS_COLLECTION_NAME` - имя коллекции задач (по умолчанию: `tasks`)
- `FILTERS_FILEPATH` - путь к файлу фильтра `fiscalDriveNumber` (по умолчанию: `filters.xls`)
- `AGGREGATION_DIRPATH` - путь к директории результатов агрегация (по умолчанию: `aggregations`)
- `SERVER_ADDR` - адрес http-сервера (применительно к http.Server) (по умолчанию: `:8080`)
4. Соберите и запустите проект:
```sh
go build -o build/app cmd/main.go && ./build/app
```