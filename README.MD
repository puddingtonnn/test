# Chat API Service

Тестовое задание: API чатов и сообщений.
Сервис реализует REST API для создания чатов, отправки сообщений и получения истории.

## Технологии

* **Language:** Go 1.22 (net/http, slog)
* **Database:** PostgreSQL 16
* **ORM:** GORM
* **Migrations:** Goose (embedded)
* **Infrastructure:** Docker & Docker Compose
* **Testing:** Testify + Mockery

## Быстрый запуск

Для запуска требуется установленный Docker и Docker Compose.

1. **Запуск сервиса:**
   ```bash
   make run
   # или
   docker-compose up --build