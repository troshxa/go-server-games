# go-server-games

REST API сервер для управління колекцією ігор, написаний на Go.

## Технології

- **Go** — мова програмування
- **chi** — HTTP роутер
- **pgx** — драйвер PostgreSQL
- **golang-migrate** — міграції бази даних
- **viper** — конфігурація через змінні середовища
- **Docker** — запуск PostgreSQL

## Запуск

### 1. Клонуй репозиторій

```bash
git clone <repo-url>
cd go-server-games
```

### 2. Створи `.env` файл

```bash
cp .env.example .env
```

### 3. Запусти базу даних

```bash
docker compose up -d
```

### 4. Запусти сервер

```bash
go run .
```

При старті сервер автоматично застосує міграції та підніметься на `http://localhost:8080`.

## Ендпоінти

| Метод | URL | Опис |
|---|---|---|
| `GET` | `/games` | Отримати всі ігри |
| `GET` | `/games/{id}` | Отримати гру за ID |
| `POST` | `/games` | Створити гру |
| `PUT` | `/games/{id}` | Повністю оновити гру |
| `PATCH` | `/games/{id}` | Частково оновити гру |
| `DELETE` | `/games/{id}` | Видалити гру |

## Приклад запиту

```bash
curl -X POST http://localhost:8080/games \
  -H "Content-Type: application/json" \
  -d '{"title":"Hades","release_date":"2020-09-17T00:00:00Z","price":24.99,"rating":10}'
```

## Структура проекту

```
.
├── migrations/         # SQL міграції
├── bruno/              # Колекція запитів для Bruno
├── main.go             # Точка входу
├── migration.go        # Запуск міграцій
├── game.go             # Модель
├── game_repo.go        # Робота з базою даних
├── game_handler.go     # HTTP хендлери
├── docker-compose.yml  # PostgreSQL
└── .env.example        # Шаблон змінних середовища
```
