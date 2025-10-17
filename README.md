# Warehouse Control System

Система управления складом с аудитом изменений и ролевой моделью доступа.

## 🚀 Быстрый старт

### 1. Настройка окружения

Создайте файл `.env` в корне проекта:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=warehouse-control
JWT_SECRET=your_jwt_secret_key_here
GOOSE_DRIVER=postgres
GOOSE_MIGRATION_DIR=migrations
```

### 2. Запуск через Docker Compose

```bash
# Запуск всех сервисов
docker-compose up --build

# Или в фоновом режиме
docker-compose up -d --build
```

### 3. Доступ к приложению

- **Веб-интерфейс**: http://localhost:8080
- **API**: http://localhost:8080/api
- **База данных**: localhost:5432

## 📋 Функциональность

### Аутентификация
- Регистрация пользователей с ролями
- Вход в систему через JWT токены
- Автоматическое сохранение сессии

### Управление товарами
- **Создание** товаров (admin, manager)
- **Просмотр** списка товаров (все роли)
- **Редактирование** товаров (admin, manager)
- **Удаление** товаров (только admin)

### Аудит изменений
- Автоматическое логирование всех изменений через триггеры PostgreSQL
- Просмотр истории изменений с фильтрацией
- Экспорт истории в CSV формат
- Отображение кто, что и когда изменил

### Роли пользователей
- **Admin** - полный доступ ко всем функциям
- **Manager** - создание/редактирование товаров, просмотр истории
- **Viewer** - только просмотр товаров и истории

## 🛠 API Endpoints

### Аутентификация
- `POST /api/auth/login` - Вход в систему
- `POST /api/auth/register` - Регистрация пользователя

### Товары
- `GET /api/items` - Получить все товары
- `GET /api/items/:id` - Получить товар по ID
- `POST /api/items` - Создать товар (admin, manager)
- `PUT /api/items/:id` - Обновить товар (admin, manager)
- `DELETE /api/items/:id` - Удалить товар (admin)

### Аудит
- `GET /api/audit` - Получить историю изменений
- `GET /api/audit?item_id=123` - История по конкретному товару
- `GET /api/audit/export` - Экспорт истории в CSV (admin, manager)

## 🏗 Архитектура

### Backend (Go)
- **Handlers** - HTTP обработчики
- **Services** - Бизнес-логика
- **Repository** - Работа с базой данных
- **Middleware** - JWT аутентификация и авторизация
- **Models** - Структуры данных

### Frontend (Vanilla JS)
- **HTML** - Структура страниц
- **CSS** - Современные стили с градиентами и анимациями
- **JavaScript** - Логика приложения и API взаимодействие

### База данных (PostgreSQL)
- **Триггеры** - Автоматическое логирование изменений
- **Миграции** - Управление схемой БД
- **Индексы** - Оптимизация запросов

## 🔧 Разработка

### Локальная разработка

1. **Установите зависимости:**
```bash
go mod download
```

2. **Запустите PostgreSQL:**
```bash
docker run --name postgres-warehouse -e POSTGRES_PASSWORD=password -e POSTGRES_DB=warehouse-control -p 5432:5432 -d postgres
```

3. **Примените миграции:**
```bash
goose postgres "host=localhost port=5432 user=postgres password=password dbname=warehouse-control sslmode=disable" up
```

4. **Запустите приложение:**
```bash
go run cmd/warehouse-control/main.go
```

### Структура проекта

```
warehouse-control/
├── cmd/warehouse-control/     # Точка входа приложения
├── internal/
│   ├── api/handlers/         # HTTP обработчики
│   ├── api/router/           # Маршрутизация
│   ├── api/server/           # HTTP сервер
│   ├── config/               # Конфигурация
│   ├── dto/                  # Data Transfer Objects
│   ├── middleware/           # Middleware (JWT, CORS)
│   ├── model/                # Модели данных
│   ├── repository/           # Слой данных
│   └── service/              # Бизнес-логика
├── migrations/               # SQL миграции
├── pkg/jwt/                  # JWT утилиты
├── web/                      # Фронтенд
│   ├── index.html
│   ├── styles.css
│   └── app.js
├── docker-compose.yml
├── Dockerfile
└── README.md
```

## 🎯 Особенности реализации

### Антипаттерн с триггерами
Как требовалось в задании, для аудита используются триггеры PostgreSQL - это антипаттерн, который демонстрирует проблемы такого подхода:

- **Проблемы**: Сложность отладки, привязка к конкретной СУБД, сложность тестирования
- **Правильный подход**: Логирование на уровне приложения через события/интерцепторы

### Ролевая модель
- JWT токены содержат информацию о роли пользователя
- Middleware проверяет права доступа для каждого запроса
- UI динамически показывает/скрывает функции в зависимости от роли

### Современный фронтенд
- Адаптивный дизайн для всех устройств
- Плавные анимации и переходы
- Уведомления о действиях пользователя
- Модальные окна для форм

## 🐛 Устранение неполадок

### Проблемы с базой данных
```bash
# Проверьте статус PostgreSQL
docker ps | grep postgres

# Посмотрите логи
docker logs postgres-warehouse
```

### Проблемы с миграциями
```bash
# Примените миграции вручную
goose postgres "host=localhost port=5432 user=postgres password=password dbname=warehouse-control sslmode=disable" up
```

### Проблемы с JWT
- Убедитесь, что JWT_SECRET установлен в .env файле
- Проверьте, что токен не истек
