# Распределенный Калькулятор

Сервис распределенного калькулятора с поддержкой аутентификации пользователей, хранением данных и асинхронным вычислением математических выражений.

## Возможности

- Регистрация и аутентификация пользователей с использованием JWT токенов
- Хранение выражений и результатов вычислений в базе данных SQLite
- Асинхронное вычисление сложных математических выражений
- Многопользовательский режим с изоляцией данных
- REST API для всех операций

## Быстрый старт

### Требования

- Go 1.19 или новее
- SQLite3

### Установка

1. Клонируйте репозиторий
2. Перейдите в директорию проекта
3. Запустите сервис:

```bash
go run ./cmd/calc_service/...
```

По умолчанию сервис запустится на порту 8080.

## Документация API

### Аутентификация

#### Регистрация нового пользователя

```bash
curl --location 'http://localhost:8080/api/v1/register' \
--header 'Content-Type: application/json' \
--data '{
    "login": "testuser",
    "password": "testpass123"
}'
```

Успешный ответ (200 OK):
```json
{
    "message": "User registered successfully"
}
```

Ответ с ошибкой (400 Bad Request):
```json
{
    "error": "user with this login already exists"
}
```

#### Вход в систему

```bash
curl --location 'http://localhost:8080/api/v1/login' \
--header 'Content-Type: application/json' \
--data '{
    "login": "testuser",
    "password": "testpass123"
}'
```

Успешный ответ (200 OK):
```json
{
    "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

Ответ с ошибкой (401 Unauthorized):
```json
{
    "error": "invalid login or password"
}
```

### Операции калькулятора

#### Отправка выражения на вычисление

```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Authorization: Bearer YOUR_JWT_TOKEN' \
--header 'Content-Type: application/json' \
--data '{
    "expression": "2+2*2"
}'
```

Успешный ответ (202 Accepted):
```json
{
    "id": "expr_123",
    "status": "pending"
}
```

#### Получение результата вычисления

```bash
curl --location 'http://localhost:8080/api/v1/expressions/expr_123' \
--header 'Authorization: Bearer YOUR_JWT_TOKEN'
```

Успешный ответ (200 OK):
```json
{
    "id": "expr_123",
    "expression": "2+2*2",
    "status": "completed",
    "result": 6
}
```

#### Получение списка выражений пользователя

```bash
curl --location 'http://localhost:8080/api/v1/expressions' \
--header 'Authorization: Bearer YOUR_JWT_TOKEN'
```

Успешный ответ (200 OK):
```json
{
    "expressions": [
        {
            "id": "expr_123",
            "expression": "2+2*2",
            "status": "completed",
            "result": 6
        }
    ]
}
```

## Обработка ошибок

API использует стандартные HTTP коды состояния:

- 200: Успешное выполнение
- 201: Создано
- 202: Принято к обработке
- 400: Неверный запрос
- 401: Не авторизован
- 404: Не найдено
- 500: Внутренняя ошибка сервера

## Разработка

### Структура проекта

```
.
├── agent/          # Реализация агента-вычислителя
├── cmd/            # Точки входа приложения
├── handlers/       # Обработчики HTTP запросов
├── middleware/     # Промежуточное ПО
├── models/         # Модели данных
├── services/       # Бизнес-логика
└── utils/          # Вспомогательные функции
```

### Запуск тестов

```bash
go test ./...
```

## Безопасность

- Все пароли хешируются с использованием bcrypt перед сохранением
- Для аутентификации используются JWT токены
- Данные хранятся в базе данных SQLite
- Данные каждого пользователя изолированы от других