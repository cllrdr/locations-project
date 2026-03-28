# REST API Documentation

## Base URL
```
http://localhost:8080/api
```

## Домен "Локации" (аналог Services)

| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/api/locations?location={name}` | Список локаций с фильтрацией |
| GET | `/api/locations/:id` | Одна локация |
| POST | `/api/locations` | Создание локации |
| PUT | `/api/locations/:id` | Обновление локации |
| DELETE | `/api/locations/:id` | Удаление локации (soft delete) |
| POST | `/api/locations/:id/image` | Загрузка изображения |
| POST | `/api/locations/:id/addLocation` | Добавить локацию в заявку-черновик |

### Примеры:

**GET список:**
```bash
GET /api/locations?location=пляж
```

**POST создание:**
```json
POST /api/locations
{
  "name": "Пляж",
  "description": "Красивый пляж",
  "players": "2-4"
}
```

**POST загрузка изображения:**
```bash
POST /api/locations/1/image
Content-Type: multipart/form-data
image: <file>
```

---

## Домен "Заявки" (Requests)

| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/api/requests/cart` | Иконка корзины (id черновика + счётчик) |
| GET | `/api/requests?status=&start_date=&end_date=` | Список с фильтрацией |
| GET | `/api/requests/:id` | Одна заявка + локации |
| PUT | `/api/requests/:id` | Обновление полей заявки |
| DELETE | `/api/requests/:id` | Удаление заявки (статус → "удалён") |
| PUT | `/api/requests/:id/form` | Сформировать черновик |
| PUT | `/api/requests/:id/complete` | Завершить/отклонить заявку |

### Примеры:

**GET список с фильтрацией:**
```bash
GET /api/requests?status=сформирован&start_date=2025-01-01&end_date=2025-12-31
```

**PUT завершение:**
```json
PUT /api/requests/1/complete
{
  "approve": true
}
```

---

## Домен "М-М связь" (RequestLocations)

| Метод | URL | Описание |
|-------|-----|----------|
| PUT | `/api/requestlocations/:id/location/:locationId` | Обновление приоритета |
| DELETE | `/api/requestlocations/:id/location/:locationId` | Удалить локацию из заявки |

### Примеры:

**PUT приоритет:**
```json
PUT /api/requestlocations/1/location/5
{
  "priority": 10
}
```

---

## Домен "Пользователь" (Profile)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/api/profile/register` | Регистрация |
| POST | `/api/profile/login` | Аутентификация |
| POST | `/api/profile/logout` | Деавторизация |
| GET | `/api/profile/me` | Данные текущего пользователя |
| PUT | `/api/profile/me` | Обновление профиля |

### Примеры:

**POST регистрация:**
```json
POST /api/profile/register
{
  "email": "user@example.com",
  "name": "Иван",
  "password": "12345"
}
```

**POST логин:**
```json
POST /api/profile/login
{
  "email": "user@example.com",
  "password": "12345"
}
```

---

## Статусы заявок

- `черновик` → можно удалять, формировать
- `сформирован` → можно завершать, отклонять
- `завершён` → конечный статус
- `отклонён` → конечный статус
- `удалён` → конечный статус (soft delete)

---

## Запуск

1. Поднять Docker:
```bash
docker-compose up -d
```

2. Запустить приложение:
```bash
./lab
# или
go run ./cmd/lab/main.go
```

3. Проверить API:
```bash
curl http://localhost:8080/api/locations
```
