# ACMEsky

Progetto di Ingegneria del Software Orientata ai Servizi A.A. 2023/2024

## Deploy using Docker compose

Deploy all

```bash
docker-compose \
    -f docker-compose-camunda.yaml \
    -f docker-compose-ACMESky.yaml \
    -f docker-compose-FlightCompany.yaml \
    up -d
```

Deploy camunda needed for ACMESky

```bash
docker compose -f docker-compose-camunda.yaml up -d
```

Deploy ACMESky platform

```bash
docker compose -f ./docker-compose-ACMESky.yaml up -d
```

Deploy (a) FlightCompany platform

```bash
docker compose -f ./docker-compose-FlightCompany.yaml up -d
```


## Docs

The documentation is generated with [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/).

Arch linux installation:
```shell
yay -S python-mkdocs-material pymdown-extensions
mkdocs serve
```
