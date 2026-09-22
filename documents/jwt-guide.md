# JWT Guide

How tokens are issued, what they contain, and how to use them against this API.

## How to get a token

Register a user (once), then login:

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "User Test",
    "email": "user.test@mail.com",
    "password": "secret123"
  }'

curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user.test@mail.com",
    "password": "secret123"
  }'
```

The login response contains the token:

```json
{
  "status": "success",
  "code": "200",
  "user": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

## How to use the token

Send it on every protected endpoint as a `Bearer` token in the `Authorization` header:

```bash
curl http://localhost:8080/users \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

Protected routes: `GET /users`, `GET /users/{id}`, `PATCH /users/{id}`, `DELETE /users/{id}`.
`/auth/register` and `/auth/login` don't require a token.