# Bank Service

This service simulates a Bank API provider.

## API:

| Endpoint            | Type   | Parameters         |
| ------------------- | ------ | ------------------ |
| `/payment/{id}`     | GET    | **id**: payment id |
| `/payment/{id}`     | DELETE | **id**: payment id |
| `/payment/new`      | PUT    |                    |
| `/payment/pay/{id}` | POST   | **id**: payment id |

## Disclaimer

The payment ID needs to be updated in the client component in order to work.

## How to start it

Start with docker:

```sh
docker compose up -d
docker build . -t bankservice
docker run --network=host -dp 127.0.0.1:8000:8000 bankservice:latest
```

Start manually:

```sh
docker compose up -d
cd server
go mod tidy
go run main.go
curl -X PUT -i http://localhost:8000/payment/new --data '{"user":"mario","amount":100,"description":"volo bologna milano"}'
```

Copy the id and update the client, then do

```sh
cd client
pnpm i
pnpm dev
```

Now, you can use the client to pay the newly created payment.

Use [swaggo](https://github.com/swaggo/swag) to generate the Swagger API documentation:

```sh
cd server
swag init --parseDependency --instanceName bank 
```
