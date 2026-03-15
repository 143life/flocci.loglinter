# flocciloglint

[![Go Reference](https://pkg.go.dev/badge/github.com/143life/flocci.loglinter.svg)](https://pkg.go.dev/github.com/143life/flocci.loglinter)

Линтер для проверки лог-сообщений на соответствие корпоративным стандартам:
- сообщения начинаются со строчной буквы;
- сообщения только на английском языке (ASCII);
- отсутствие эмодзи и специальных символов;
- отсутствие чувствительных данных (пароли, токены и т.п.).

## Возможности

- Гибкая настройка через YAML/JSON.
- Работа как самостоятельный анализатор (`go vet`).
- Интеграция с `golangci-lint` через систему модульных плагинов (начиная с v2.0).

## Установка

### 1. Как отдельный анализатор (для `go vet`)

```bash
go install github.com/143life/flocci.loglinter/cmd/flocciloglinter@latest
```

Запуск:
```bash
go vet -vettool=$(which flocciloglinter) ./...
```

### 2. Как плагин для golangci-lint

Для использования с `golangci-lint` требуется собрать кастомную версию, включающую этот плагин.

#### 2.1. Подготовьте файл `.custom-gcl.yml` в корне вашего проекта:

```yaml
version: v2.0.0
name: custom-gcl
destination: ./bin

plugins:
  - module: 'github.com/143life/flocci.loglinter'
    import: 'github.com/143life/flocci.loglinter/pkg/golangci'
    path: .   # путь до локальной копии, если вы разрабатываете плагин
```

#### 2.2. Соберите кастомный бинарник:

```bash
golangci-lint custom
```

Будет создан исполняемый файл `./bin/custom-gcl` (или указанное вами имя).

#### 2.3. Настройте `.golangci.yml`:

```yaml
version: "2"
linters:
  default: none
  enable:
    - flocciloglint
  settings:
    custom:
      flocciloglint:
        type: "module"
        description: "Log message linter"
        settings:
          sensitive_patterns:
            - "password"
            - "token"
            - "secret"
            - "api.?key"
          check_first_lowercase: true
          forbid_emoji: true
          allow_only_ascii: true
          use_regex: false   # если хотите использовать регулярные выражения в sensitive_patterns
```

#### 2.4. Запустите проверку:

```bash
./bin/custom-gcl run ./...
```

## Конфигурация

Все настройки задаются в секции `settings` внутри `.golangci.yml` или в отдельном YAML/JSON файле, путь к которому можно передать через флаг `-config`.

### Доступные параметры

| Поле | Тип | По умолчанию | Описание |
|------|-----|--------------|----------|
| `sensitive_patterns` | `[]string` | `["password","token","secret","key"]` | Список слов/паттернов, считающихся чувствительными |
| `use_regex` | `bool` | `false` | Если `true`, элементы `sensitive_patterns` интерпретируются как регулярные выражения |
| `check_first_lowercase` | `bool` | `true` | Проверять, что сообщение начинается со строчной буквы |
| `forbid_emoji` | `bool` | `true` | Запрещать эмодзи |
| `allow_only_ascii` | `bool` | `false` | Разрешать только ASCII символы (полезно для английского языка) |

### Пример конфигурационного файла (`config.yaml`)

```yaml
sensitive_patterns:
  - "password"
  - "token"
  - "secret.*key"
  - "credit.?card"
use_regex: true
check_first_lowercase: true
forbid_emoji: true
allow_only_ascii: true
```

## Примеры

### Плохо

```go
log.Printf("user password: %s", pwd)                 // sensitive: password
log.Println("Ошибка подключения")                     // non-ASCII
log.Info("server started 😊")                          // emoji
log.Error("Failed to connect!!!")                     // special chars
log.Warn("warning: something went wrong...")          // special chars
slog.Error("Starting server on port 8080")            // uppercase start
```

### Хорошо

```go
log.Printf("user authenticated successfully")
log.Println("connection timeout")
log.Info("server started")
slog.Error("failed to connect to database")
```

## Лицензия

MIT