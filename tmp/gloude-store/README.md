# 🌩️ Gloude Store — Cloud Storage

Production-ready облачное хранилище файлов на Go Fiber + React.

## Стек технологий

**Backend:**
- Go 1.22 + Fiber v2
- PostgreSQL + GORM
- JWT аутентификация
- Локальное файловое хранилище

**Frontend:**
- React 18 + TypeScript + Vite
- Tailwind CSS
- Zustand (state management)
- Axios + react-hot-toast

**Инфраструктура:**
- Docker + Docker Compose
- Nginx (reverse proxy + SPA)

## Быстрый старт

```bash
# 1. Клонируйте репозиторий
git clone <repo>
cd gloude-store

# 2. Создайте .env файл
cp .env.example .env

# 3. Отредактируйте секреты в .env
nano .env

# 4. Запустите все сервисы
docker-compose up -d --build

# 5. Откройте в браузере
open http://localhost
```

## Структура проекта

```
gloude-store/
├── backend/
│   ├── cmd/server/        # Точка входа
│   ├── internal/
│   │   ├── config/        # Конфигурация
│   │   ├── handler/       # HTTP обработчики
│   │   ├── middleware/    # Auth middleware
│   │   ├── models/        # GORM модели
│   │   ├── repository/    # Слой БД
│   │   └── service/       # Бизнес-логика
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── api/           # Axios + API функции
│   │   ├── components/    # UI компоненты
│   │   ├── hooks/         # Кастомные хуки
│   │   ├── pages/         # Страницы
│   │   ├── types/         # TypeScript типы
│   │   └── utils/         # Утилиты
│   └── Dockerfile
├── nginx/
│   ├── nginx.conf
│   └── default.conf
├── docker-compose.yml
└── .env.example
```

## API Endpoints

| Метод | URL | Описание |
|-------|-----|----------|
| POST | /api/v1/auth/register | Регистрация |
| POST | /api/v1/auth/login | Вход |
| POST | /api/v1/auth/logout | Выход |
| GET | /api/v1/account/me | Текущий пользователь |
| POST | /api/v1/storage/upload | Загрузка файла |
| GET | /api/v1/storage/files | Список файлов |
| GET | /api/v1/storage/download/:id | Скачать файл |
| DELETE | /api/v1/storage/:id | Удалить файл |
| GET | /api/v1/storage/quota | Квота |
| GET | /api/v1/storage/activity | Активность |
| PATCH | /api/v1/storage/:id/favorite | Избранное |

## Функциональность

- ✅ Аутентификация (JWT + HttpOnly Cookie)
- ✅ Загрузка файлов (Drag & Drop с прогресс-баром)
- ✅ Список файлов (Grid/List режимы)
- ✅ Фильтрация (по расширению, размеру, поиск)
- ✅ Скачивание файлов (streaming)
- ✅ Удаление с подтверждением
- ✅ Квоты (1 ГБ на пользователя)
- ✅ Избранное
- ✅ GitHub-style heatmap активности
- ✅ Защита от Directory Traversal
- ✅ UUID имена файлов
- ✅ Skeleton-загрузчики
- ✅ Toast уведомления

## Управление

```bash
# Запуск
docker-compose up -d

# Логи
docker-compose logs -f backend
docker-compose logs -f frontend

# Остановка
docker-compose down

# Полная очистка (включая данные!)
docker-compose down -v
```

## Переменные окружения

| Переменная | По умолчанию | Описание |
|-----------|-------------|----------|
| POSTGRES_USER | gloude | Пользователь БД |
| POSTGRES_PASSWORD | gloude_secret | Пароль БД |
| POSTGRES_DB | gloude_store | Имя БД |
| JWT_SECRET | change_me | JWT секрет (ОБЯЗАТЕЛЬНО СМЕНИТЬ!) |
| MAX_QUOTA_BYTES | 1073741824 | Квота (1 ГБ) |
| ENV | production | Окружение |
