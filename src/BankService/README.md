# Bank Service

Servizio di gestione dei pagamenti per la compagnia ACMESky.

## Commands

```bash
docker compose up -d
cd server
go mod tidy
go run main.go
cd ../client
npx tailwindcss -i ./src/input.css -o ./src/output.css --watch
pnpm i
pnpm dev
curl -X POST -i http://localhost:3000/payment/new --data '{"user":"mario","amount":100,"description":"volo bologna milano"}'
```
