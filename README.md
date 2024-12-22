# CalcAPI - Go

### Navigation: English Version
- [Description](#description)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
  - [CURL Examples](#curl-examples)
    - [Success Example (200 STATUS)](#success-example-200-status)
    - [Invalid Expression (422 STATUS: Expression is not valid)](#invalid-expression-422-status-expression-is-not-valid)
    - [Division by Zero (422 STATUS: Division by zero)](#division-by-zero-422-status-division-by-zero)
    - [Server Error (500 STATUS)](#server-error-500-status)

### Навигация: Русская версия
- [Описание](#описание)
- [Начало работы](#начало-работы)
  - [Необходимые требования](#необходимые-требования)
  - [Установка](#установка)
- [Использование](#использование)
  - [Примеры CURL](#примеры-curl)
    - [Успешный пример (200 STATUS)](#успешный-пример-200-status)
    - [Некорректное выражение (422 STATUS: Expression is not valid)](#некорректное-выражение-422-status-expression-is-not-valid)
    - [Деление на ноль (422 STATUS: Division by zero)](#деление-на-ноль-422-status-division-by-zero)
    - [Ошибка сервера (500 STATUS)](#ошибка-сервера-500-status)

---

## Description

CalcAPI is a Go-based API calculator that processes mathematical expressions provided in JSON format. It evaluates the expressions and returns either the result or an error in JSON format. The server handles requests via the POST method and supports expressions with parentheses and basic arithmetic operations.

## Getting Started

### Prerequisites
- Installed and working Go environment.

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/dolezha/golang-web-calculator.git
   cd <repository_folder>
   ```
2. Run the application:
   ```bash
   go run main.go
   ```

## Usage

### CURL Examples

#### Success Example (200 STATUS)
```bash
curl -X POST -H "Content-Type: application/json" -d '{"expression": "2+2*2"}' http://localhost:8080/api/v1/calculate
```

#### Invalid Expression (422 STATUS: Expression is not valid)
```bash
curl -X POST -H "Content-Type: application/json" -d '{"expression": "2++2"}' http://localhost:8080/api/v1/calculate
```

#### Division by Zero (422 STATUS: Division by zero)
```bash
curl -X POST -H "Content-Type: application/json" -d '{"expression": "10/0"}' http://localhost:8080/api/v1/calculate
```

#### Server Error (500 STATUS)
Occurs for other errors in the expression or during calculations.

```bash
curl -X POST -H "Content-Type: application/json" -d '{"expression": "(some invalid expression)"}' http://localhost:8080/api/v1/calculate
```

---

# CalcAPI - Go (Русская версия)

## Описание

CalcAPI — это API калькулятор на Go, который принимает математическое выражение в формате JSON, вычисляет его и возвращает результат или ошибку в формате JSON. Сервер обрабатывает запросы по методу POST и поддерживает выражения со скобками и базовыми арифметическими операциями.

## Начало работы

### Необходимые требования
- Установленная и работающая среда Go.

### Установка
1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/dolezha/golang-web-calculator.git
   cd <repository_folder>
   ```
2. Запустите приложение:
   ```bash
   go run main.go
   ```

## Использование

### Примеры CURL

#### Успешный пример (200 STATUS)
```bash
curl -X POST -H "Content-Type: application/json" -d '{"expression": "2+2*2"}' http://localhost:8080/api/v1/calculate
```

#### Некорректное выражение (422 STATUS: Expression is not valid)
```bash
curl -X POST -H "Content-Type: application/json" -d '{"expression": "2++2"}' http://localhost:8080/api/v1/calculate
```

#### Деление на ноль (422 STATUS: Division by zero)
```bash
curl -X POST -H "Content-Type: application/json" -d '{"expression": "10/0"}' http://localhost:8080/api/v1/calculate
```

#### Ошибка сервера (500 STATUS)
Возникает при иных ошибках в выражении или вычислениях.

```bash
curl -X POST -H "Content-Type: application/json" -d '{"expression": "(неверное выражение)"}' http://localhost:8080/api/v1/calculate
```
