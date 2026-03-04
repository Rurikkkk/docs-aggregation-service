# Docs Aggregation Service

Сервис агрегации документов, написанный на Go, работающий с базой данных документов MongoDB компании "Астрал-Софт". Проект разработан с применением чистой архитектуры.

## Основные возможности и особенности

- Агрегация документов коллекции по диапазону дат и значению `fiscalDriveNumber`
- Асинхронный бэкенд сервиса на основе менеджера задач
- UI на основе REST API по HTTP
- Сбор метрик веб-сервера и менеджера задач через Prometheus

## HTTP API

1. `POST /aggregate` — создание новой задачи на агрегацию документов

   **Тело запроса (multipart/form-data):**
   - `startDate`: начальная дата диапазона выборки документов в формате `YYYY-MM-DD hh:mm:ss`
   - `endDate`: конечная дата диапазона выборки документов в формате `YYYY-MM-DD hh:mm:ss`
   - `filters`: xls-файл с требуемыми значениями `fiscalDriveNumber`

   **Ответ:**
   - `202 Accepted` (application/json) — задача успешно запущена 
      ```json
      {
        "id": <ID запущенной задачи>
      }
      ```
   - `400 Bad Request` (application/json) — некорректное тело запроса
      ```json
      {
        "error": <описание ошибки>
      }
      ```
   - `405 Method Not Allowed` (no body) — использован неверный HTTP-метод
   - `500 Internal Server Error` (application/json) — внутренняя ошибка сервера
      ```json
      {
        "error": <описание ошибки>
      }
      ```

2. `GET /aggregate/status` — получение статуса одной или всех задач

   **Параметры запроса:**
   - `id` (опционально): ID требуемой задачи

   **Ответ:**
   - `200 OK` (application/json) - возврат статуса
      ```json
      [
        {
            "id": <taskID>,
            "status": <статус задачи>,
            "error": <error>
        }
      ]
      ```
   - `405 Method Not Allowed` (no body) — использован неверный HTTP-метод
   - `500 Internal Server Error` (application/json) — внутренняя ошибка сервера
      ```json
      {
        "error": <описание ошибки>
      }
      ```

3. `GET /aggregate/result` — получение результата по завершенной задаче с указанным ID

   **Параметры запроса:**
   - `id`: ID требуемой задачи

   **Ответ:**
   - `200 OK` (text/csv) - возврат CSV с результатами агрегации
   - `400 Bad Request` (application/json) — некорректные параметры запроса или задача не завершена
      ```json
      {
        "error": <описание ошибки>
      }
      ```
   - `405 Method Not Allowed` (no body) — использован неверный HTTP-метод
   - `500 Internal Server Error` (application/json) — внутренняя ошибка сервера
      ```json
      {
        "error": <описание ошибки>
      }
      ```

4. `GET /metrics` — получение Prometheus метрик сервиса

   **Ответ:**
   - `200 OK` (text/plain) - возврат метрик

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
- `MONGO_URL` — URL сервера MongoDB (по умолчанию: `mongodb://localhost:27017`)
- `MONGO_DB_NAME` — имя базы данных MongoDB (по умолчанию: `docs-aggregation-service`)
- `DOCS_COLLECTION_NAME` — имя коллекции документов (по умолчанию: `docs`)
- `TASKS_COLLECTION_NAME` — имя коллекции задач (по умолчанию: `tasks`)
- `AGGREGATION_DIRPATH` — путь к директории результатов агрегации (по умолчанию: `aggregations`)
- `SERVER_ADDR` — адрес http-сервера (применительно к http.Server) (по умолчанию: `:8080`)
4. Соберите и запустите проект:
```sh
go build -o build/app cmd/main.go && ./build/app
```