# NCC Service

This service manages NCCs.

[API documentation](http://0.0.0.0:8080/swagger/index.html)

## How to start it

### Manually

Middleware:

```
cd middleware
go mod tidy
swag init
go run .
```

Backend:

```
cd backend/
jolie main.ol
```

### Docker Compose

```
docker compose up -d
```

## Test

```
cd client
./curl.sh
```
