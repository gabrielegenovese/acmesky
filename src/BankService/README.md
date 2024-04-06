# Bank Service

Servizio di gestione dei pagamenti per la compagnia ACMESky.

## Commands

```bash
docker compose up -d
cd server
go mod tidy
go run main.go
cd ../client
pnpm i
pnpm dev
curl -X POST -i http://localhost:3000/payment/new --data '{"user":"mario","amount":100,"description":"volo bologna milano"}'
```

In seguito, sostituire nel componente l'id.