# Geographical Distance Service

This service calculates the distance between two points (coordinate or address). It uses the free [distancematrix.ai](https://distancematrix.ai) APIs.

## API:

| Endpoint    | Type | Parameters                                          |
| ----------- | ---- | --------------------------------------------------- |
| `/distance` | GET  | **from**: starting point; **to**: destination point |

## Use example

```sh
curl -X GET -i 'http://localhost:8093/distance?from=Bologna&to=Milan'
```

## How to start it

Command line:

```sh
go mod tidy
go run main.go
```

Docker:

```sh
docker build . -t geodistance/service
docker run -dp 127.0.0.1:8093:8093 geodistance/service:latest
```

Use [swaggo](https://github.com/swaggo/swag) to generate the Swagger API documentation:

```sh
swag init --parseDependency --instanceName geo
```
