# Распределенный калькулятор

Сервис для асинхронного вычисления арифметических выражений с поддержкой параллельных вычислений.

## О проекте

Система состоит из двух компонентов:
1. **Оркестратор** - сервер, который принимает выражения и разбивает их на операции
2. **Агент** - вычислитель, который выполняет отдельные операции

Особенности:
- Поддержка базовых арифметических операций (+, -, *, /)
- Параллельное выполнение независимых операций
- Масштабирование через добавление агентов
- REST API для взаимодействия

## Быстрый старт

1. Клонируйте репозиторий:
```
git clone https://github.com/dolezha/golang-web-calculator.git
cd calculator
```
2. Запустите оркестратор:
```
go run cmd/calc_service/main.go
```
3. Запустите агент:
```
go run cmd/agent/main.go
```

## Примеры использования

### 1. Отправка выражения на вычисление
```
curl -X POST "http://localhost:8080/api/v1/calculate" \
-H "Content-Type: application/json" \
-d '{"expression": "2+22"}'
```

Ответ:
```
{"id": "1234567890"}
```

### 2. Проверка статуса вычисления
```
curl -X GET "http://localhost:8080/api/v1/expressions/1234567890"
```

Ответ:
```
{
    "expression": {
        "id": "1234567890",
        "status": "done",
        "result": 6
    }
}
```

### 3. Получение списка всех выражений
```
curl -X GET "http://localhost:8080/api/v1/expressions"
```
Ответ:
```
{
    "expressions": [
        {
        "id": "1234567890",
        "status": "done",
        "result": 6
        }
    ]
}
```

### Обработка ошибок

Невалидное выражение:
```
curl -X POST "http://localhost:8080/api/v1/calculate" \
-H "Content-Type: application/json" \
-d '{"expression": "2++2"}'
```
Ответ (422):
```
{"error": "невалидное выражение"}
```
Несуществующее выражение:
```
curl -X GET "http://localhost:8080/api/v1/expressions/nonexistent"
```

Ответ (404):
```
{"error": "нет такого выражения"}
```


## Конфигурация

Настройка через переменные окружения:

### Оркестратор
- `TIME_ADDITION_MS` - время выполнения сложения (по умолчанию 1000мс)
- `TIME_SUBTRACTION_MS` - время выполнения вычитания (по умолчанию 1000мс)
- `TIME_MULTIPLICATION_MS` - время выполнения умножения (по умолчанию 2000мс)
- `TIME_DIVISION_MS` - время выполнения деления (по умолчанию 2000мс)

### Агент
- `COMPUTING_POWER` - количество параллельных вычислителей (по умолчанию 4)

## Архитектура
```
┌─────────────┐ HTTP ┌──────────────┐
│   Клиент    │──⟷──│  Оркестратор │
└─────────────┘      └──────────────┘
                        ▲
                        │
                        │ HTTP 
                        │
                        ┌──────────────┐
                        │    Агент     │
                        └──────────────┘

1. Клиент отправляет выражение оркестратору
2. Оркестратор разбивает его на операции
3. Агенты забирают и выполняют операции
4. Клиент периодически проверяет статус вычисления
```

## Тесты

Запуск тестов:
```
go test ./...
```
