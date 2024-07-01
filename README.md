# ACMEsky

Progetto di Ingegneria del Software Orientata ai Servizi A.A. 2023/2024

## Deploy using Docker compose

Avvio di camunda per ACMESky

```bash
docker compose -f docker-compose-camunda.yaml up -d
```

Avvio di ACMESky platform

```bash
docker compose -f ./docker-compose-ACMESky.yaml up -d
```

Avvio di FlightCompany platform

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
