# API Examples

Base URL: `http://localhost:8080`

All authenticated endpoints require the header:

```
Authorization: Bearer <token>
```

---

## Register

`POST /auth/register`

**Request**

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "User Test",
    "email": "user.test@mail.com",
    "password": "secret123"
  }'
```

**Response `201`**

```json
{
  "status": "success",
  "code": "201",
  "user": {
    "id": "3f9a1c2e-...."
  }
}
```

**Response `400`** (email already exists / invalid body / validation error)

```json
{
  "status": "fail",
  "code": "400",
  "message": "email already exists"
}
```

---

## Login

`POST /auth/login`

**Request**

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user.test@mail.com",
    "password": "secret123"
  }'
```

**Response `200`**

```json
{
  "status": "success",
  "code": "200",
  "user": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Response `401`** (wrong email or password)

```json
{
  "status": "fail",
  "code": "401",
  "message": "invalid email or password"
}
```

---

## List Users

`GET /users?limit=20&offset=0`

**Request**

```bash
curl http://localhost:8080/users?limit=20&offset=0 \
  -H "Authorization: Bearer <token>"
```

**Response `200`**

```json
{
  "status": "success",
  "code": "200",
  "users": [
    {
      "id": "3f9a1c2e-....",
      "name": "User Test",
      "email": "user.test@mail.com",
      "created_at": "2026-09-22T10:00:00+07:00"
    }
  ],
  "meta": {
    "limit": 20,
    "offset": 0
  }
}
```

Optional query params: `limit`, `offset`, `from`, `to` (`from`/`to` are RFC3339 timestamps, filter by `created_at`).

---

## Get User by ID

`GET /users/{id}`

**Request**

```bash
curl http://localhost:8080/users/3f9a1c2e-.... \
  -H "Authorization: Bearer <token>"
```

**Response `200`**

```json
{
  "status": "success",
  "code": "200",
  "user": {
    "id": "3f9a1c2e-....",
    "name": "User Test",
    "email": "user.test@mail.com",
    "created_at": "2026-09-22T10:00:00+07:00"
  }
}
```

**Response `404`**

```json
{
  "status": "fail",
  "code": "404",
  "message": "user not found"
}
```

---

## Update User

`PATCH /users/{id}`

**Request** (only send the field(s) you want to change)

```bash
curl -X PATCH http://localhost:8080/users/3f9a1c2e-.... \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Name"
  }'
```

**Response `200`**

```json
{
  "status": "success",
  "code": "200",
  "user": {
    "id": "3f9a1c2e-....",
    "name": "New Name",
    "email": "user.test@mail.com",
    "created_at": "2026-09-22T10:00:00+07:00"
  }
}
```

**Response `400`** (no fields sent, or new email already used by another user)

```json
{
  "status": "fail",
  "code": "400",
  "message": "no fields to update"
}
```

---

## Delete User

`DELETE /users/{id}`

**Request**

```bash
curl -X DELETE http://localhost:8080/users/3f9a1c2e-.... \
  -H "Authorization: Bearer <token>"
```

**Response `200`**

```json
{
  "status": "success",
  "code": "200",
  "message": "user deleted successfully"
}
```

**Response `404`**

```json
{
  "status": "fail",
  "code": "404",
  "message": "user not found"
}
```

---

## Common Error: Missing/Invalid Token

Any authenticated endpoint, without a valid `Authorization: Bearer <token>` header:

**Response `401`**

```json
{
  "status": "fail",
  "code": "401",
  "message": "missing token"
}
```
