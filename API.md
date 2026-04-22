# REST API Documentation

## Base URL
```
http://localhost:8080/api
```

## Домен "Локации" (Services)

| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/api/locations?location={name}` | Список локаций с фильтрацией |
| GET | `/api/locations/:id` | Одна локация |
| POST | `/api/locations` | Создание локации с картинкой и видео |

### Примеры для Postman:

**1. GET список локаций:**
```
GET http://localhost:8080/api/locations?location=пляж
```

**2. GET одна локация:**
```
GET http://localhost:8080/api/locations/1
```

**3. POST создание локации с файлами:**
```
POST http://localhost:8080/api/locations
Content-Type: multipart/form-data

Fields:
- name: "Пляж Золотой"
- description: "Красивый песчаный пляж"
- players: "2-6"
- image: <выбрать файл>
- video: <выбрать файл>
```

---

## Домен "Заявки" (Requests)

| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/api/requests/cart` | Иконка корзины (id черновика + счётчик локаций) |
| GET | `/api/requests?status=&start_date=&end_date=` | Список (кроме черновиков и удаленных) |
| GET | `/api/requests/:id` | Одна заявка со списком локаций |
| PUT | `/api/requests/:id` | Обновление только nickname |
| DELETE | `/api/requests/:id` | Удаление (статус → "удалён") |
| PUT | `/api/requests/:id/form` | Сформировать черновик (переход в "сформирован") |
| PUT | `/api/requests/:id/complete` | Завершить/отклонить (модератором) |

### Примеры для Postman:

**1. GET корзина текущего пользователя:**
```
GET http://localhost:8080/api/requests/cart
```
Ответ:
```json
{
  "draft_id": 5,
  "locations_cnt": 2
}
```

**2. GET список заявок с фильтром:**
```
GET http://localhost:8080/api/requests?status=сформирован&start_date=2025-01-01&end_date=2025-12-31
```

**3. GET одна заявка:**
```
GET http://localhost:8080/api/requests/2
```

**4. PUT обновление заявки:**
```
PUT http://localhost:8080/api/requests/2
Content-Type: application/json

{
  "nickname": "Иван Петров"
}
```

**5. PUT сформировать заявку:**
```
PUT http://localhost:8080/api/requests/2/form
Content-Type: application/json
```

**6. PUT завершить/отклонить заявку:**
```
PUT http://localhost:8080/api/requests/2/complete
Content-Type: application/json

{
  "approve": true
}
```

**7. DELETE удалить заявку:**
```
DELETE http://localhost:8080/api/requests/2
```

---

## Домен "М-М связь" (Локации в Заявке)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/api/requests/:requestId/locations/:locationId` | Добавить локацию в черновик |
| PUT | `/api/requests/:requestId/locations/:locationId` | Изменить приоритет локации |
| DELETE | `/api/requests/:requestId/locations/:locationId` | Удалить локацию из заявки |

### Примеры для Postman:

**1. POST добавить локацию в заявку:**
```
POST http://localhost:8080/api/requests/2/locations/5
Content-Type: application/json
```

**2. PUT изменить приоритет:**
```
PUT http://localhost:8080/api/requests/2/locations/5
Content-Type: application/json

{
  "priority": 10
}
```

**3. DELETE удалить локацию:**
```
DELETE http://localhost:8080/api/requests/2/locations/5
```

---

## Домен "Пользователь" (Profile)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/api/profile/register` | Регистрация пользователя |
| POST | `/api/profile/login` | Аутентификация (заглушка) |
| POST | `/api/profile/logout` | Деавторизация (заглушка) |

### Примеры для Postman:

**1. POST регистрация:**
```
POST http://localhost:8080/api/profile/register
Content-Type: application/json

{
  "email": "ivan@example.com",
  "name": "Иван Петров",
  "password": "password123"
}
```

**2. POST логин:**
```
POST http://localhost:8080/api/profile/login
Content-Type: application/json

{
  "email": "ivan@example.com",
  "password": "password123"
}
```

**3. POST логаут:**
```
POST http://localhost:8080/api/profile/logout
Content-Type: application/json
```

