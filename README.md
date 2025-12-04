# Subscription Service

REST-сервис для агрегации данных об онлайн-подписках пользователей.

## 📋 Описание

Сервис предоставляет API для управления подписками пользователей с возможностью:
- CRUD операции над записями о подписках
- Подсчет суммарной стоимости подписок за период с фильтрацией
- Работа с данными в формате JSON
- Использование PostgreSQL в качестве СУБД

## 🚀 Быстрый старт

### Требования
- Go 1.20+
- PostgreSQL 12+
- Docker (опционально)

### 1. Установка зависимостей
```bash
go mod tidy
```

### 2. Настройка базы данных
Создайте базу данных PostgreSQL и настройте переменные окружения:

```bash
# Windows PowerShell
$env:DATABASE_URL = "postgres://user:password@localhost:5432/subscription_service?sslmode=disable"
$env:PORT = "8080"
$env:LOG_LEVEL = "info"
```

### 3. Запуск миграций
```bash
# Windows PowerShell
./scripts/migrate.ps1
```

### 4. Запуск сервиса
```bash
# Windows PowerShell
./scripts/run.ps1

# Или напрямую
go run cmd/server/main.go
```

Сервис будет доступен по адресу: `http://localhost:8080`

## 📚 API Документация

### Endpoints

#### 1. Создание подписки
```http
POST /api/v1/subscriptions
Content-Type: application/json

{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025",
    "end_date": "12-2025"
}
```

#### 2. Получение всех подписок
```http
GET /api/v1/subscriptions
```

#### 3. Получение подписки по ID
```http
GET /api/v1/subscriptions/{id}
```

#### 4. Обновление подписки
```http
PUT /api/v1/subscriptions/{id}
Content-Type: application/json

{
    "service_name": "Updated Service",
    "price": 500,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "08-2025"
}
```

#### 5. Удаление подписки
```http
DELETE /api/v1/subscriptions/{id}
```

#### 6. Подсчет суммарной стоимости
```http
GET /api/v1/subscriptions/cost?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&service_name=Yandex Plus&period=07-2025
```

### Формат данных

#### Модель подписки
- `service_name` (string) - название сервиса
- `price` (integer) - стоимость в рублях (целое число)
- `user_id` (UUID) - идентификатор пользователя в формате UUID
- `start_date` (string) - дата начала подписки в формате MM-YYYY
- `end_date` (string, optional) - дата окончания подписки в формате MM-YYYY

## 🐳 Docker

### Запуск с Docker Compose
```bash
docker-compose up --build
```

### Отдельный запуск PostgreSQL
```bash
docker run -d \
  --name postgres-subscription \
  -e POSTGRES_USER=user \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=subscription_service \
  -p 5432:5432 \
  postgres:15
```

## 🧪 Тестирование

### Автоматическое тестирование API
```bash
# Windows PowerShell
./scripts/test-api.ps1
```

### Unit тесты
```bash
go test ./tests/unit/...
```

### Ручное тестирование

1. **Health Check**
```bash
curl http://localhost:8080/health
```

2. **Создание подписки**
```bash
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```

3. **Получение всех подписок**
```bash
curl http://localhost:8080/api/v1/subscriptions
```

4. **Подсчет стоимости**
```bash
curl "http://localhost:8080/api/v1/subscriptions/cost?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba"
```

### 📁 **Структура проекта**

```
.
├── cmd/server/             # Точка входа приложения
├── internal/
│   ├── api/
│   │   ├── handlers/       # HTTP обработчики
│   │   └── middleware/     # Middleware для логирования и CORS
│   ├── config/            # Конфигурация приложения
│   ├── logger/            # Логирование
│   ├── model/             # Модели данных
│   ├── repository/        # Слой работы с БД
│   └── service/           # Бизнес-логика
├── db/migrations/         # SQL миграции
├── api/                   # OpenAPI спецификация
├── scripts/               # Скрипты для запуска и тестирования
├── tests/                 # Тесты
├── docker/                # Docker конфигурация
├── docker-compose.yml     # Docker Compose конфигурация
├── go.mod                 # Go модуль
└── README.md              # Документация
```

## 🔧 Конфигурация

Переменные окружения:
- `DATABASE_URL` - строка подключения к PostgreSQL
- `PORT` - порт для HTTP сервера (по умолчанию 8080)
- `LOG_LEVEL` - уровень логирования (info, debug, error)

## 📋 Особенности реализации

1. **Архитектура**: Clean Architecture с разделением на слои
2. **Фреймворк**: Gin для HTTP сервера
3. **База данных**: PostgreSQL с нативным драйвером lib/pq
4. **Валидация**: Валидация дат в формате MM-YYYY
5. **Логирование**: Структурированное логирование запросов
6. **API**: RESTful API с OpenAPI документацией
7. **Тестирование**: Unit тесты с mock репозиторием

## 🚨 Примечания

- Стоимость подписки указывается в рублях как целое число (копейки не учитываются)
- Даты используют формат MM-YYYY (например, "07-2025")
- User ID должен быть в формате UUID
- Проверка существования пользователей не требуется согласно заданию
```
bash scripts/migrate.sh
```

## Scripts

- `scripts/run.sh`: Starts the application.
- `scripts/migrate.sh`: Runs database migrations.

## Docker

To build the Docker image, use:
```
docker build -t subscription-service -f docker/Dockerfile .
```

To run the application using Docker Compose:
```
docker-compose up
```

## Testing

To run tests, use:
```
go test ./...
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the LICENSE file for details.