# Flight Company Service

## Docker build

```sh
docker build . -t flightcompany/service
```

## Docs

To generate swagger docs:

```
swag init --parseDependency --instanceName fc
```
