# Kefir SSO — Сервис аутентификации и авторизации

Высокопроизводительный SSO-сервис на Go с gRPC API, PostgreSQL и Redis. Обеспечивает безопасную аутентификацию и авторизацию для микросервисной архитектуры.

---

## Ключевые возможности

| Функция | Описание |
|---------|----------|
| **Регистрация** | Создание новых пользователей с валидацией email и пароля |
| **Вход** | Выдача JWT-токенов после проверки учетных данных |
| **Админ-панель** | Проверка прав администратора |
| **Валидация токенов** | Проверка подлинности JWT |
| **Выход** | Добавление токенов в blacklist (Redis) |
| **Кэширование** | Redis для быстрого доступа к пользователям |

---

## Быстрый старт

### Предварительные требования
- Docker & Docker Compose
- Make (опционально)

### Запуск одной командой

```bash
git clone https://github.com/Kefir4c/sso-service
cd sso-service
docker compose up -d
```
Проверка работы
```bash
# Проверить логи
docker compose logs -f app
```
### Запустить тесты
```bash
docker compose run --rm test go test ./tests/... -v
```
# Структура проекта
```bash
sso-service/
│
├── cmd/                                    # Точки входа
│   ├── sso/                                # Основное приложение
│   │   └── main.go                         # Запуск gRPC сервера
│   └── migrator/                           # Миграции БД
│       └── main.go                         # Запуск миграций
│
├── config/                                 # Конфигурация
│   ├── local.yaml                          # Для локальной разработки
│   ├── remote.yaml                         # Для Docker/prod
│   └── example.yaml                        # Пример конфига
│
├── internal/                               # Внутренний код (private)
│   │
│   ├── app/                                # Сборка приложения
│   │   ├── app.go                          # Инициализация всех компонентов
│   │   └── grpc/
│   │       └── app.go                      # Настройка gRPC сервера
│   │
│   ├── config/                             # Загрузка конфигов
│   │   └── config.go                       # Парсинг YAML
│   │
│   ├── domain/                             # Модели данных
│   │   └── models/
│   │       ├── user.go                     # Модель пользователя
│   │       └── app.go                      # Модель приложения
│   │
│   ├── grpc/                               # gRPC слой
│   │   └── auth/
│   │       └── server.go                   # Хендлеры Register, Login, IsAdmin...
│   │
│   ├── lib/                                # Вспомогательные пакеты
│   │   ├── jwt/
│   │   │   └── jwt.go                      # Генерация/валидация JWT
│   │   ├── logger/                         # Логирование
│   │   │   ├── handlers/
│   │   │   │   └── slogpretty/
│   │   │   │       └── slogpretty.go       # Красивое форматирование логов
│   │   │   └── sl/
│   │   │       └── error.go                # Хелперы для ошибок
│   │   └── validation/                     # Валидация
│   │       ├── email.go                    # Проверка email
│   │       └── password.go                 # Проверка пароля
│   │
│   ├── logger/                             # Настройка логгера
│   │   └── logger.go                       # Инициализация
│   │
│   ├── services/                           # Бизнес-логика
│   │   └── auth/
│   │       └── auth.go                     # Регистрация, логин, IsAdmin...
│   │
│   └── storage/                            # Хранилища данных
│       ├── cache/                          # Кэш (Redis)
│       │   ├── cache.go                    # Интерфейс кэша
│       │   └── redis/
│       │       └── redis.go                # Реализация Redis
│       └── postgres/                       # PostgreSQL
│           ├── conn.go                     # Подключение к БД
│           ├── user_repo.go                # Работа с пользователями
│           └── app_repo.go                 # Работа с приложениями
│
├── migrations/                             # SQL миграции
│   ├── 1_init.up.sql                       # Создание таблиц
│   ├── 1_init.down.sql                     # Откат создания таблиц
│   ├── 2_test_app.up.sql                   # Добавление тестовых данных
│   └── 2_test_app.down.sql                 # Откат тестовых данных
│
├── tests/                                  # Интеграционные тесты
│   ├── suite/                              # Обвязка для тестов
│   │   └── suite.go                        # Настройка тестового окружения
│   ├── testdata/                           # Тестовые данные
│   │   └── testdata.go                     # Константы для тестов
│   ├── auth_test.go                        # Основные тесты
│   └── register_login_test.go              # Тесты регистрации/логина
│
├── docker-compose.yaml                     # Docker Compose конфиг
├── Dockerfile.app                          # Dockerfile для приложения
├── Dockerfile.migrator                     # Dockerfile для мигратора
├── Dockerfile.test                         # Dockerfile для тестов
├── go.mod                                  # Зависимости
├── go.sum                                  # Контрольные суммы
├── Taskfile.yaml                           # Taskfile команды
├── .gitignore                              # Что не коммитить
└── README.md                               # Документация
```
## Конфигурация
Создайте config/local.yaml на основе примера:
```bash
yaml
grpc:
  port: 44044
  timeout: 5s

storage:
  host: localhost
  port: 5432
  dbname: sso
  username: postgres

redis:
  addr: localhost:6379
  password: ""
  db: 0
  ```
## Тестирование
```bash
# Запустить все тесты
docker compose run --rm test go test ./tests/... -v

# Запустить конкретный тест
docker compose run --rm test go test ./tests/... -v -run TestHappyPath_RegisterLogin
```
```bash
Текущий статус тестов
TestHappyPath_RegisterLogin	        ✅ Успешно
TestRegister_InvalidPasswordLength	✅ Успешно
TestRegister_PasswordComplexity	    ✅ Успешно
TestRegister_Duplicate	            ✅ Успешно
TestLogin_FailCases(частично)	    ⚠️ В работе
```
# Технологический стек
* Язык: Go 1.23+

* API: gRPC + Protocol Buffers

* База данных: PostgreSQL 18

* Кэш: Redis 8

* Аутентификация: JWT + bcrypt

* Контейнеризация: Docker + Docker Compose

* Миграции: golang-migrate

## Документация API
proto-файлы доступны в отдельном репозитории: github.com/Kefir4c/proto_sso

## Основные методы gRPC:
* Register — регистрация пользователя

* Login — вход и получение токена

* IsAdmin — проверка прав администратора

* ValidateToken — валидация JWT

* Logout — выход (добавление в blacklist)

## Docker-инфраструктура
```bash
# Сборка образов
docker compose build

# Запуск всех сервисов
docker compose up -d

# Просмотр логов
docker compose logs -f app

# Остановка
docker compose down
```
# Известные ограничения и планы
## Текущие ограничения
* Тесты IsAdmin и ValidateToken требуют доработки

* Отсутствует Swagger-документация

## Планы по развитию
* Добавить поддержку OAuth2

* Реализовать rate limiting

* Добавить метрики (Prometheus)

* Написать больше тестов

# Автор
* Kefir4c
* GitHub: @Kefir4c