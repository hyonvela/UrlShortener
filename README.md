# URL Shortener Service

Сервис для сокращения URL-адресов с поддержкой HTTP и gRPC API, использующий PostgreSQL для постоянного хранения или Redis в качестве in-memory решения.

## Особенности
- Поддержка HTTP и gRPC API
- Хранение данных в PostgreSQL или Redis на выбор
- Метрики Prometheus
- Логирование logrus
- Конфигурация через YAML файл

## API Documentation
### HTTP API

- **POST /v1/url_shortener**
  - **Request Body**: 
    ```json
    {
      "long_url": "https://example.com"
    }
    ```
  - **Response**:
    ```json
    {
      "short_url": "abc123"
    }
    ```

- **GET /v1/url_shortener**
  - **Request Body**: 
    ```json
    {
      "short_url": "abc123"
    }
    ```
  - **Response**:
    ```json
    {
      "long_url": "https://example.com"
    }
    ```

### gRPC API

- **ShortenUrl**
  - **Request**: 
    ```protobuf
    message UrlRequest {
      string url = 1;
    }
    ```
  - **Response**: 
    ```protobuf
    message UrlResponse {
      string url = 1;
    }
    ```

- **GetLongUrl**
  - **Request**: 
    ```protobuf
    message UrlRequest {
      string url = 1;
    }
    ```
  - **Response**: 
    ```protobuf
    message UrlResponse {
      string url = 1;
    }
    ```

**Пример proto-файла:** `proto/url_shortener.proto`

## Запуск проекта

### Запустить сервисы:
```bash
make up
```

### Остановить сервисы:
```bash
make down
```

### Тестирование:
```bash
make test
```

### Миграции баз данных:
```bash
make migrate
```

## Конфигурация

Настройки сервиса можно изменить в файле `config/config.yaml`:

```yaml
logs_format: text
logs_lvl: debug
storage: postgres
listen:
  bind_ip: 0.0.0.0
  http_port: 8080
  grpc_port: 9090
  write_timeout: 15
  read_timeout: 15
database:
  db_host: db
  db_port: 5432
  username: postgres
  password: password
  db_name: testdb
  ssl_mode: disable
cache:
  redis_host: redis
  redis_port: 6379
  redis_db: 0
  life_time: 10
```

## Технологии
- Go 1.23
- PostgreSQL
- Redis
- gRPC
- Gin Framework
- Prometheus
- Docker